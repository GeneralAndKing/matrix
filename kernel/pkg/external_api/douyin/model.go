package douyin

type WordCover struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

type Sentence struct {
	ChallengeID       string      `json:"challenge_id"`
	DiscussVideoCount int         `json:"discuss_video_count"`
	DisplayStyle      int         `json:"display_style"`
	DriftInfo         interface{} `json:"drift_info"`
	EventTime         int64       `json:"event_time"`
	GroupID           string      `json:"group_id"`
	HotValue          int         `json:"hot_value"`
	HotlistParam      string      `json:"hotlist_param"`
	Label             int         `json:"label"`
	RelatedWords      interface{} `json:"related_words"`
	SentenceID        string      `json:"sentence_id"`
	VideoCount        int         `json:"video_count"`
	Word              string      `json:"word"`
	WordCover         WordCover   `json:"word_cover"`
	WordSubBoard      interface{} `json:"word_sub_board"`
	WordType          int         `json:"word_type"`
	PostAwemeInfo     string      `json:"post_aweme_info,omitempty"`
	ViewCount         int         `json:"view_count,omitempty"`
}

type Extra struct {
	FatalItemIDs []interface{} `json:"fatal_item_ids"`
	LogID        string        `json:"logid"`
	Now          int64         `json:"now"`
}

type LogPB struct {
	ImprID string `json:"impr_id"`
}

type RecommendHotspotResult struct {
	AllSentences []Sentence    `json:"all_sentences"`
	Extra        Extra         `json:"extra"`
	LogPB        LogPB         `json:"log_pb"`
	RecSentences []interface{} `json:"rec_sentences"`
	StatusCode   int           `json:"status_code"`
}

type SearchHotspotResult struct {
	Sentences    []Sentence    `json:"sentences"`
	Extra        Extra         `json:"extra"`
	LogPB        LogPB         `json:"log_pb"`
	RecSentences []interface{} `json:"rec_sentences"`
	StatusCode   int           `json:"status_code"`
}

type ChallengeSugResult struct {
	SugList []struct {
		ChaName   string `json:"cha_name"`
		ViewCount int    `json:"view_count"`
		Cid       string `json:"cid"`
		GroupId   string `json:"group_id"`
		Tag       int    `json:"tag"`
	} `json:"sug_list"`
	StatusCode       int    `json:"status_code"`
	StatusMsg        string `json:"status_msg"`
	Rid              string `json:"rid"`
	WordsQueryRecord struct {
		Info        string `json:"info"`
		WordsSource string `json:"words_source"`
		QueryId     string `json:"query_id"`
	} `json:"words_query_record"`
	Extra struct {
		Now             int64         `json:"now"`
		Logid           string        `json:"logid"`
		FatalItemIds    []interface{} `json:"fatal_item_ids"`
		SearchRequestId string        `json:"search_request_id"`
	} `json:"extra"`
	LogPb struct {
		ImprId string `json:"impr_id"`
	} `json:"log_pb"`
}

type ActivityResult struct {
	ActivityList []struct {
		ActivityId     string   `json:"activity_id"`
		ActivityLevel  int      `json:"activity_level"`
		ActivityName   string   `json:"activity_name"`
		ActivityStatus int      `json:"activity_status"`
		ActivityType   int      `json:"activity_type"`
		Challenge      []string `json:"challenge"`
		ChallengeIds   []int64  `json:"challenge_ids"`
		CollectId      int      `json:"collect_id"`
		CollectStatus  bool     `json:"collect_status"`
		CoverImage     string   `json:"cover_image"`
		GameId         string   `json:"game_id"`
		HotScore       int      `json:"hot_score"`
		IfWellChosen   bool     `json:"if_well_chosen"`
		JumpLink       string   `json:"jump_link"`
		JumpType       int      `json:"jump_type"`
		QueryTag       int      `json:"query_tag"`
		RewardType     int      `json:"reward_type"`
		ShowEndTime    string   `json:"show_end_time"`
		ShowStartTime  string   `json:"show_start_time"`
	} `json:"activity_list"`
	Extra struct {
		Logid string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
	StatusCode int `json:"status_code"`
}

type FlashmobShootResult struct {
	Display int `json:"display"`
	Extra   struct {
		FatalItemIds []interface{} `json:"fatal_item_ids"`
		Logid        string        `json:"logid"`
		Now          int64         `json:"now"`
	} `json:"extra"`
	HotList int `json:"hot_list"`
	LogPb   struct {
		ImprId string `json:"impr_id"`
	} `json:"log_pb"`
	RecommendList []struct {
		Text string `json:"text"`
	} `json:"recommend_list"`
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FlashmobInfoResult struct {
	Extra struct {
		Logid string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
	FlashMobInfoMap []struct {
		FlashMobInfos struct {
			Count int `json:"count"`
			Cover struct {
				Uri     string   `json:"uri"`
				UrlList []string `json:"url_list"`
			} `json:"cover"`
			Id             string      `json:"id"`
			UserAvatarList interface{} `json:"user_avatar_list"`
		} `json:"flash_mob_infos"`
		FlashMobText string `json:"flash_mob_text"`
	} `json:"flash_mob_info_map"`
}

type FlashmobRankResult struct {
	StatusCode   int    `json:"status_code"`
	StatusMsg    string `json:"status_msg"`
	RankCellList []struct {
		FlashMobId      string `json:"flash_mob_id"`
		AwemeId         string `json:"aweme_id"`
		Count           int    `json:"count"`
		Text            string `json:"text"`
		Order           int    `json:"order"`
		CreatorNickname string `json:"creator_nickname"`
		InitiatorUid    string `json:"initiator_uid"`
		Cover           struct {
			Uri     string   `json:"uri"`
			UrlList []string `json:"url_list"`
			Width   int      `json:"width"`
			Height  int      `json:"height"`
		} `json:"cover"`
	} `json:"rank_cell_list"`
}
