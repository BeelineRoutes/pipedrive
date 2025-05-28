/** ****************************************************************************************************************** **
	Calls related to users (crew)

    
** ****************************************************************************************************************** **/

package pipedrive 

import (
    "github.com/pkg/errors"
    
    "fmt"
    "net/http"
    "net/url"
    "context"
    "time"
    "strconv"
    "strings"
)

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ----------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- STRUCTS ---------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

type Activity struct {
    Id, Company_id, User_id, Org_id, Deal_id int
    Type, Due_date, Due_time, Duration, Subject, Location, Org_name, Lead_title, Note string 
    Done bool 

    Start time.Time // these get filled in by us
    Dur time.Duration 
}

// we want to convert our date/time info into a timestamp we can use
func (this *Activity) setupStartTimes (start time.Time) (err error) {
    if len(this.Due_date) < 8 { return nil } // nothing more to do

    // figure out the time for this
    if len(this.Due_time ) < 4 { // expected format "03:00"
        // not sure what to do with this, unscheduled job for the date maybe?
        
        this.Start, err = time.Parse("2006-01-02 15:04:05", this.Due_date + " " + start.Format("15:04:05"))
        if err != nil { return err }
    } else {
        // we have time too
        this.Start, err = time.Parse("2006-01-02 15:04", this.Due_date + " " + this.Due_time)
        if err != nil { return err }
    }

    // we got a start time, now do our end time
    parts := strings.Split(this.Duration, ":")
    if len(parts) != 2 { 
        return nil  // not returning this as an error
    }

    hours, err := strconv.Atoi(parts[0])
    if err != nil { return err }
    minutes, err := strconv.Atoi(parts[1])
    if err != nil { return err }

    this.Dur = time.Hour * time.Duration(hours) + time.Minute * time.Duration(minutes)
    
    return nil 
}

type activityResponse struct {
    Data []*Activity
    Success bool 
    Additional_data struct {
        Pagination struct {
            More_items_in_collection bool 
        }
    }
}

// takes the jobs out of whatever this parent object is for
func (this activityResponse) toActivites (start, finish time.Time) (ret []*Activity, err error) {
    for _, m := range this.Data {
        if m.Done { continue } // it's done, so don't worry about it

        err = m.setupStartTimes (start) // figure this conversion out
        if err != nil { return }
        
        if m.Start.IsZero() { continue } // no date set, probably can't get pulled in, but whatever
        
        // see if it's in our target window
        if m.Start.Before(start) { continue } // this isn't in our window
        if m.Start.Add(m.Dur).After(finish) { continue } // this isn't in our window either        
        
        // we got one
        ret = append (ret, m)
    }
    return 
}

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- PRIVATE FUNCTIONS -----------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

  //-----------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS -------------------------------------------------------------------------------------------------------//
//-----------------------------------------------------------------------------------------------------------------------//

// returns all jobs that match our conditions
func (this *Pipedrive) ListActivities (ctx context.Context, bearer, domain string, start, finish time.Time) (ret []*Activity, err error) {
    
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + bearer 

    // we need some query params for this
    params := url.Values{}
    params.Set("user_id", "0") // this tells pipedrive to return all the activities for all users
    params.Set("done", "0") // don't return done jobs, although we check for that anyway
    params.Set("limit", "200")

    // figure out our dates
    params.Set("start_date", start.Format("2006-01-02")) // just use this as a base
    params.Set("end_date", finish.AddDate(0, 0, 1).Format("2006-01-02")) // we need to add a day to wrap around
    // we then "double check" the dates before including the activity to be returned

    for i := 0; i <= 5; i++ { // stay in a loop as long as we're pulling jobs
        params.Set("start", fmt.Sprintf("%d", i * 200)) // this isn't a page, it's the running count
        
        resp := &activityResponse{}
        err = this.send (ctx, http.MethodGet, fmt.Sprintf("%s/v1/activities?%s", domain, params.Encode()), header, nil, resp)
        if err != nil { return } // bail

        acts, lErr := resp.toActivites(start, finish)
        if lErr != nil { return nil, lErr }

        ret = append (ret, acts...) // add these to our return

        // see if we're done
        if resp.Additional_data.Pagination.More_items_in_collection == false {
            return // we're good
        }
    }

    return // we hit our limit
}

// updates what we need to about an activity
func (this *Pipedrive) UpdateActivity (ctx context.Context, bearer, domain string, id, userId int, start time.Time) error {
    
    header := make(map[string]string)
    header["Authorization"] = "Bearer " + bearer 

    var data struct {
        DueDate string `json:"due_date"`
        DueTime string `json:"due_time"`
        User int `json:"user_id"`
    }

    data.DueDate = start.Format("2006-01-02")
    data.DueTime = start.Format("15:04")
    data.User = userId

    var resp struct {
        Success bool 
    }

    err := this.send (ctx, http.MethodPut, fmt.Sprintf("%s/v1/activities/%d", domain, id), header, data, &resp)
    if err != nil { return err } // bail

    if resp.Success != true {
        return errors.Errorf ("update didn't return success: %s : %d : %d : %s", domain, id, userId, start)
    }

    return nil // we're good 
}
