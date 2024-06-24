package douyin

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

var douyinClient = resty.New().SetBaseURL("https://creator.douyin.com")

func FetchRecommendHotspot() ([]HotspotResponse, error) {
	var (
		result    RecommendHotspotResult
		responses = make([]HotspotResponse, 0)
	)
	resp, err := douyinClient.R().
		SetResult(&result).
		Get("/aweme/v1/hotspot/recommend")

	if err != nil {
		return nil, fmt.Errorf("failed to request douyin recommend hotspot")
	}
	if resp.IsSuccess() {
		for _, sentence := range result.AllSentences {
			responses = append(responses, HotspotResponse{
				EventTime:  time.Unix(sentence.EventTime, 0),
				HotValue:   sentence.HotValue,
				VideoCount: sentence.VideoCount,
				Word:       sentence.Word,
				Cover:      sentence.WordCover.URLList[0],
			})
		}
		return responses, nil
	}
	return nil, fmt.Errorf("failed to request douyin recommend hotspot")

}

func FetchSearchHotspot(keyword string) ([]HotspotResponse, error) {
	var (
		result    SearchHotspotResult
		responses = make([]HotspotResponse, 0)
	)
	resp, err := douyinClient.R().SetResult(&result).
		SetQueryParams(map[string]string{
			"query": keyword,
			"count": "50",
		}).
		Get("/aweme/v1/hotspot/search/")
	if err != nil {
		return nil, fmt.Errorf("failed to request douyin search hotspot")
	}
	if resp.IsSuccess() {
		for _, sentence := range result.Sentences {
			responses = append(responses, HotspotResponse{
				EventTime:  time.Unix(sentence.EventTime, 0),
				HotValue:   sentence.HotValue,
				VideoCount: sentence.VideoCount,
				Word:       sentence.Word,
				Cover:      sentence.WordCover.URLList[0],
			})
		}
		return responses, nil
	}
	return nil, fmt.Errorf("failed to request douyin search hotspot")
}

func FetchChallengeSug(keyword string) ([]ChallengeSugResponse, error) {
	var (
		result    ChallengeSugResult
		responses = make([]ChallengeSugResponse, 0)
	)
	resp, err := douyinClient.R().SetResult(&result).
		SetQueryParams(map[string]string{
			"source":  "challenge_create",
			"aid":     "2906",
			"keyword": keyword,
		}).
		Get("/aweme/v1/search/challengesug/")
	if err != nil {
		return nil, fmt.Errorf("failed to request challenge sug")
	}
	if resp.IsSuccess() {
		for _, sug := range result.SugList {
			responses = append(responses, ChallengeSugResponse{
				Name:      sug.ChaName,
				ViewCount: sug.ViewCount,
			})
		}
		return responses, nil
	}
	return nil, fmt.Errorf("failed to request challenge sug")
}

func FetchActivity(cookies []*http.Cookie) ([]ActivityResponse, error) {
	var (
		result    ActivityResult
		responses = make([]ActivityResponse, 0)
	)
	resp, err := douyinClient.R().SetResult(&result).
		SetCookies(cookies).
		SetQueryParams(map[string]string{
			"query_tag":      "0",
			"page":           "0",
			"size":           "9999",
			"need_challenge": "1",
		}).
		Get("/web/api/media/activity/get/")
	if err != nil {
		return nil, fmt.Errorf("failed to request activity")
	}
	if resp.IsSuccess() {
		for _, activity := range result.ActivityList {
			responses = append(responses, ActivityResponse{
				Cover:      activity.CoverImage,
				HotScore:   activity.HotScore,
				Name:       activity.ActivityName,
				Challenges: activity.Challenge,
				StartTime:  activity.ShowStartTime,
				EndTime:    activity.ShowEndTime,
			})
		}
		return responses, nil
	}
	return nil, fmt.Errorf("failed to request activity")

}

func FetchFlashmob(cookies []*http.Cookie) ([]FlashmobResponse, error) {
	return nil, nil
}
