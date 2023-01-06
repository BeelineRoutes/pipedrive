/** ****************************************************************************************************************** **
    Pipedrive API wrapper
    written for GoLang
    Created 2022-06-20 by Nathan Thomas 
    Courtesy of BeelineRoutes.com

    current docs in v1
    https://developer.pipedrive.com/

** ****************************************************************************************************************** **/

package pipedrive 

import (
    "github.com/pkg/errors"

    "fmt"
    "context"
    "net/url"
    "encoding/base64"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

const apiURL = "https://beelineroutes.pipedrive.com/v1"

var (
    ErrUnexpected       = errors.New("idk...")
	ErrAuthExpired      = errors.New("Auth Expired")
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type apiError struct {
    Success bool 
    Error, Error_info string 
    ErrorCode int 
}

type Oauth struct {
    AccessToken string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ApiDomain string `json:"api_domain"`
    Expires int `json:"expires_in"`
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CLASS -----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type Pipedrive struct {
	OAuthClientId, OAuthClientSecret, RedirectURI string
}

// inits our pipedrive api object with the required default info for everyone
// returns the object and a bool to indicate we're goods
func NewPipedrive (clientId, clientSecret, redirectUri string) (*Pipedrive, bool) {
    ret := &Pipedrive {
        OAuthClientId: clientId,
        OAuthClientSecret: clientSecret,
        RedirectURI: redirectUri,
    }

    return ret, ret.valid()
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// validation check to be used on startup
func (this *Pipedrive) valid () bool {
    if len(this.OAuthClientId) == 0 { return false }
    if len(this.OAuthClientSecret) == 0 { return false }
    if len(this.RedirectURI) == 0 { return false }

    return true 
}

// returns a base64 encoded hash of the id and secret
// this is needed for the header Authorization to get a bearer token and refresh it
func (this *Pipedrive) hashAuth () string {
    return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", this.OAuthClientId, this.OAuthClientSecret)))
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PUBLIC FUNCTIONS ------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// inital oauth call using the short lived token from the callback
func (this *Pipedrive) OAuth (ctx context.Context, code string) (*Oauth, error) {
    resp := &Oauth{}
    
    // start our header
    header := make(map[string]string)
    header["Authorization"] = "Basic " + this.hashAuth()

    // start our form params
    params := url.Values{}
    params.Set("grant_type", "authorization_code")
    params.Set("code", code)
    params.Set("redirect_uri", this.RedirectURI)

    err := this.sendForm (ctx, "https://oauth.pipedrive.com/oauth/token", header, params, &resp)
    return resp, err
}

// refreshes our api token using our old refresh token
func (this *Pipedrive) RefreshToken (ctx context.Context, oldRefresh string) (*Oauth, error) {
    resp := &Oauth{}

    // start our header
    header := make(map[string]string)
    header["Authorization"] = "Basic " + this.hashAuth()

    // start our form params
    params := url.Values{}
    params.Set("grant_type", "refresh_token")
    params.Set("refresh_token", oldRefresh)

    err := this.sendForm (ctx, "https://oauth.pipedrive.com/oauth/token", header, params, &resp)
    return resp, err
}
