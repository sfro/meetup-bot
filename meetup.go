package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// https://secure.meetup.com/meetup_api/console/?path=/2/categories
// Tech is id 34 ;)

type meetupResponse struct {
	Results []Result `json:"results"`
}

type Venue struct {
	Country              string  `json:"country"`
	LocalizedCountryName string  `json:"localized_country_name"`
	City                 string  `json:"city"`
	Address1             string  `json:"address_1"`
	Name                 string  `json:"name"`
	Longitude            float64 `json:"lon"`
	ID                   int     `json:"id"`
	Latitude             float64 `json:"lat"`
	Repinned             bool    `json:"repinned"`
}

type Group struct {
	JoinMode       string  `json:"join_mode"`
	Created        int     `json:"created"`
	Name           string  `json:"name"`
	GroupLongitude float64 `json:"group_lon"`
	ID             int     `json:"id"`
	URLName        string  `json:"urlname"`
	GroupLatitude  float64 `json:"group_lat"`
	Who            string  `json:"who"`
}

type Result struct {
	UTOffset       int     `json:"utc_offset"`
	Venue          Venue   `json:"venue"`
	RSVPLimit      int     `json:"rsvp_limit"`
	HeadCount      int     `json:"headcount"`
	Distance       float64 `json:"distance"`
	Visibility     string  `json:"visibility"`
	WaitlistCount  int     `json:"waitlist_count"`
	Created        int     `json:"created"`
	MaybeRSVPCount int     `json:"maybe_rsvp_count"`
	Description    string  `json:"description"`
	EventURL       string  `json:"event_url"`
	YesRSVPCount   int     `json:"yes_rsvp_count"`
	Duration       int     `json:"duration"`
	Name           string  `json:"name"`
	ID             string  `json:"id"`
	Time           int64   `json:"time"`
	Updated        int     `json:"updated"`
	Group          Group   `json:"group"`
	Status         string  `json:"status"`
}

func getTechMeetups() ([]Result, error) {
	url := fmt.Sprintf("https://api.meetup.com/2/open_events?key=%s&sign=true&photo-host=public&country=se&city=stockholm&category=34&page=20", apiKeys["meetup"])

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	meetupResponse := meetupResponse{}

	err = json.NewDecoder(resp.Body).Decode(&meetupResponse)
	if err != nil {
		return nil, err
	}

	return meetupResponse.Results, nil
}

func meetupSuggestion() (*Result, error) {
	rand.Seed(time.Now().Unix())

	results, err := getTechMeetups()
	if err != nil {
		return nil, err
	}

	runnerUps := []Result{}

	for _, result := range results {
		if result.RSVPLimit-result.YesRSVPCount > 0 && result.Status == "upcoming" {
			runnerUps = append(runnerUps, result)
		}
	}

	return &runnerUps[rand.Intn(len(runnerUps))], nil
}
