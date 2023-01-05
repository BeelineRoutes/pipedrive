
package pipedrive 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestUsers (t *testing.T) {
	pd, cfg := parseConfig(t)
	
	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of users
	users, err := pd.ListUsers (ctx, cfg.AccessToken, cfg.ApiDomain)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(users) > 0, "expecting at least 1 user")
	assert.NotEqual (t, 0, users[0].Id, "not filled in")
	assert.NotEqual (t, "", users[0].Name, "not filled in")
	
	for _, j := range users {
		t.Logf ("%+v\n", j)
	}
	
}
