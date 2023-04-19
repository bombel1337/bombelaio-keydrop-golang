package models




type GettingLoggedInStruct struct {
	UserName string `json:"userName"`
	SteamID  string `json:"steamId"`
	Avatar   string `json:"avatar"`
}



type MonitoringGiveawayStruct struct {
	Success    bool `json:"success"`
	Pagination struct {
		ItemsCount   int `json:"itemsCount"`
		ItemsPerPage int `json:"itemsPerPage"`
		CurrentPage  int `json:"currentPage"`
	} `json:"pagination"`
	Data []struct {
		ID                string      `json:"id"`
		Status            string      `json:"status"`
		MaxUsers          int         `json:"maxUsers"`
		MinUsers          int         `json:"minUsers"`
		HaveIJoined       bool        `json:"haveIJoined"`
		MySlot            interface{} `json:"mySlot"`
		PublicHash        string      `json:"publicHash"`
		DeadlineTimestamp int64       `json:"deadlineTimestamp"`
		Frequency         string      `json:"frequency"`
		Prizes            []struct {
			ID         int     `json:"id"`
			Color      string  `json:"color"`
			ItemImg    string  `json:"itemImg"`
			Title      string  `json:"title"`
			Subtitle   string  `json:"subtitle"`
			Price      float64 `json:"price"`
			Condition  string  `json:"condition"`
			WeaponType string  `json:"weaponType"`
			Currency   string  `json:"currency"`
		} `json:"prizes"`
		ParticipantCount int           `json:"participantCount"`
		Winners          []interface{} `json:"winners"`
	} `json:"data"`
}

type JoinGiveawayStruct struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		IDGiveaway  int    `json:"idGiveaway"`
		IDSteam     string `json:"idSteam"`
		Username    string `json:"username"`
		SteamAvatar string `json:"steamAvatar"`
		ClientSeed  string `json:"clientSeed"`
		Ticket      int    `json:"ticket"`
		Slot        int    `json:"slot"`
	} `json:"data"`
}

type CaptchaTaskID struct {
	ErrorID int `json:"errorId"`
	TaskID  int `json:"taskId"`
	ErrorCode        string `json:"errorCode"`
}



type CaptchaResult struct {
	Solution struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
		Cookies            struct {
			Nocookies string `json:"nocookies"`
		} `json:"cookies"`
	} `json:"solution"`
	Status           string      `json:"status"`
	ErrorID          int         `json:"errorId"`
	ErrorCode        interface{} `json:"errorCode"`
	ErrorDescription interface{} `json:"errorDescription"`
}



type WinnersChecker struct {
	Success bool `json:"success"`
	Data    struct {
		ID                    string `json:"id"`
		MySteamID             string `json:"mySteamId"`
		MaxUsers              int    `json:"maxUsers"`
		MinUsers              int    `json:"minUsers"`
		DepositAmountRequired int    `json:"depositAmountRequired"`
		DepositAmountMissing  int    `json:"depositAmountMissing"`
		PublicHash            string `json:"publicHash"`
		DeadlineTimestamp     int64  `json:"deadlineTimestamp"`
		Status                string `json:"status"`
		Prizes                []struct {
			ID         int     `json:"id"`
			Color      string  `json:"color"`
			ItemImg    string  `json:"itemImg"`
			Title      string  `json:"title"`
			Subtitle   string  `json:"subtitle"`
			Price      float64 `json:"price"`
			Condition  string  `json:"condition"`
			WeaponType string  `json:"weaponType"`
			Currency   string  `json:"currency"`
		} `json:"prizes"`
		CanIJoin         bool        `json:"canIJoin"`
		BlockedUntil     interface{} `json:"blockedUntil"`
		HaveIJoined      bool        `json:"haveIJoined"`
		MySlot           int         `json:"mySlot"`
		Participants     []string    `json:"participants"`
		ParticipantCount int         `json:"participantCount"`
		Frequency        string      `json:"frequency"`
		Winners          []struct {
			PrizeID  int `json:"prizeId"`
			Userdata struct {
				IDSteam     string `json:"idSteam"`
				Username    string `json:"username"`
				SteamAvatar string `json:"steamAvatar"`
				Ticket      int    `json:"ticket"`
				Slot        int    `json:"slot"`
				ClientSeed  string `json:"clientSeed"`
			} `json:"userdata"`
		} `json:"winners"`
	} `json:"data"`
}

type PrizesWinner struct {
		ID         int     `json:"id"`
		Color      string  `json:"color"`
		ItemImg    string  `json:"itemImg"`
		Title      string  `json:"title"`
		Subtitle   string  `json:"subtitle"`
		Price      float64 `json:"price"`
		Condition  string  `json:"condition"`
		WeaponType string  `json:"weaponType"`
		Currency   string  `json:"currency"`
}