
package pipedrive 

import (
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"

	"testing"
	"encoding/json"
	"os"
	"context"
	"time"
)

type testingConfig struct {
	RefreshToken, ApiDomain, AccessToken string
}

// for local testing. fill out the example_config.json and rename it
// but don't put your stuff in the repo!
func parseConfig (t *testing.T) (*Pipedrive, *testingConfig) {
	config, err := os.Open("./config.json")
	if err != nil { t.Fatal (err) }
	
	jsonParser := json.NewDecoder (config)

    ret := &Pipedrive{}
	err = jsonParser.Decode (ret)
	if err != nil { t.Fatal (err) }

	if ret.Valid() == false {
		t.Fatal (errors.Errorf("config.json has missing params"))
	}

	config.Close()
	// now do it again cause the config also has our previous refresh token
	config, err = os.Open("./config.json")
	if err != nil { t.Fatal (err) }

	defer config.Close()
	
	jsonParser = json.NewDecoder (config)

	cfg := &testingConfig{}
	err = jsonParser.Decode (cfg)
	if err != nil { t.Fatal (err) }

	return ret, cfg
}

func getBearer (t *testing.T) *Oauth {
	pd, cfg := parseConfig(t)
	
	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	oauth, err := pd.RefreshToken (ctx, cfg.RefreshToken)
	if err != nil { t.Fatal(err) }

	return oauth
}

func TestRefreshToken (t *testing.T) {
	oauth := getBearer (t) 
	
	assert.NotEqual (t, 0, len(oauth.AccessToken))
	assert.NotEqual (t, 0, len(oauth.RefreshToken))
	assert.NotEqual (t, 0, len(oauth.ApiDomain))
	assert.Equal (t, 3599, oauth.Expires)

	jstr, err := json.Marshal(oauth)
	if err != nil { t.Fatal(err) }

	t.Logf("%s\n", string(jstr))
}
