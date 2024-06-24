package douyin

import "time"

type HotspotResponse struct {
	EventTime  time.Time `json:"eventTime"`
	HotValue   int       `json:"hotValue"`
	VideoCount int       `json:"videoCount"`
	Word       string    `json:"word"`
	Cover      string    `json:"cover"`
}

type ChallengeSugResponse struct {
	Name      string `json:"name"`
	ViewCount int    `json:"viewCount"`
}

type ActivityResponse struct {
	Cover      string   `json:"cover"`
	HotScore   int      `json:"hotScore"`
	Name       string   `json:"name"`
	Challenges []string `json:"challenges"`
	StartTime  string   `json:"startTime"`
	EndTime    string   `json:"endTime"`
}

type FlashmobResponse struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Cover string `json:"cover"`
}
