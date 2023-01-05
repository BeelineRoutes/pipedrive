/** ****************************************************************************************************************** **
	The actual sending and receiving stuff
	Reused for most of the calls to Pipedrive
	
** ****************************************************************************************************************** **/

package pipedrive 

import (
    "github.com/pkg/errors"

    "net/http"
	"net/url"
    "context"
    "encoding/json"
    "io/ioutil"
    "bytes"
	"strings"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// handles making the request and reading the results from it 
// if there's an error the Error object will be set, otherwise it will be nil
func (this *Pipedrive) finish (req *http.Request, out interface{}) error {
	resp, err := http.DefaultClient.Do (req)
	
	if err != nil { return errors.WithStack (err) }
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll (resp.Body)

    if resp.StatusCode > 399 { 
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			// special error 
			return errors.Wrapf (ErrAuthExpired, "Unauthorized : %d : %s", resp.StatusCode, string(body))
		}
		// just a default
		err = errors.Wrapf (ErrUnexpected, "Pipedrive Error : %d : %s", resp.StatusCode, string(body))

		// see if we can figure out the error
		errResp := &apiError{}
		jErr := errors.WithStack (json.Unmarshal (body, errResp))
		if jErr == nil {
			// we don't know what to do with this error
			err = errors.Wrapf (err, "%s : %s", errResp.Error, errResp.Error_info)
		} else {
			// different error object than expected
			err = errors.Wrapf (err, "unmarshal : %s", jErr.Error())
		}
        
        return err
    }
	
	if out != nil { 
		err = errors.WithStack (json.Unmarshal (body, out))
		if err != nil {
			err = errors.Wrap (err, string(body)) // if it didn't unmarshal, include the body so we know what it did look like
		}
	}
	
	return err // we're good
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

func (this *Pipedrive) send (ctx context.Context, requestType, link string, header map[string]string, in, out interface{}) error {
	var jstr []byte 
	var err error 

	if in != nil {
		jstr, err = json.Marshal (in)
		if err != nil { return errors.WithStack (err) }

		header["Content-Type"] = "application/json; charset=utf-8"
	}
	
	req, err := http.NewRequestWithContext (ctx, requestType, link, bytes.NewBuffer(jstr))
	if err != nil { return errors.Wrap (err, link) }

	for key, val := range header { req.Header.Set (key, val) }
	err = this.finish (req, out)
	
	return errors.Wrapf (err, " %s : %s", link, string(jstr))
}

// sends a form-urlencoded request
func (this *Pipedrive) sendForm (ctx context.Context, link string, header map[string]string, data url.Values, out interface{}) error {
	
	header["Content-Type"] = "application/x-www-form-urlencoded"
	
	req, err := http.NewRequestWithContext (ctx, http.MethodPost, link, strings.NewReader(data.Encode()))
	if err != nil { return errors.Wrap (err, link) }

	for key, val := range header { req.Header.Set (key, val) }
	err = this.finish (req, out)
	
	return errors.Wrapf (err, " %s : %s", link, data.Encode())
}

