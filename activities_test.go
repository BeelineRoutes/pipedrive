
package pipedrive 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestActivities (t *testing.T) {
	pd, cfg := parseConfig(t)
	
	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	start, err := time.Parse("2006-01-02 15:04", "2023-01-17 13:00")
	if err != nil { t.Fatal (err) }
	end := start.AddDate(0, 0, 1)

	// get our list of activities
	activities, err := pd.ListActivities (ctx, cfg.AccessToken, cfg.ApiDomain, start, end)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(activities) > 0, "expecting at least 1 activity")
	assert.NotEqual (t, 0, activities[0].Id)
	assert.NotEqual (t, "", activities[0].Location)
	assert.Equal (t, "2023-01-18 03:00", activities[0].Start.Format("2006-01-02 15:04"))
	assert.Equal (t, 15, int(activities[0].Dur.Minutes()))
	
	for _, j := range activities {
		t.Logf ("%+v\n", j)
	}
	
}

func TestActivitiesUpdate (t *testing.T) {
	pd, cfg := parseConfig(t)
	
	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	start, err := time.Parse("2006-01-02 15:04", "2023-01-17 23:00")
	if err != nil { t.Fatal (err) }
	
	// make the update
	err = pd.UpdateActivity (ctx, cfg.AccessToken, cfg.ApiDomain, 2, 17106790, start)
	if err != nil { t.Fatal (err) }	
}
