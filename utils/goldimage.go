package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
	// "bytes"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/sirupsen/logrus"
	"strings"
)

var lastDiscordContent string
var firstPing bool = true

func DiscordMonitorGold(users map[string][]Users) {

	var options []tls_client.HttpClientOption = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_108),
	}
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
		DiscordMonitorGold(users)
		return
	}
	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/v9/channels/868574854536888401/messages?limit=50", nil)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
		DiscordMonitorGold(users)
		return
	}

	req.Header = http.Header{
		"cookie":             {"__dcfduid=401d2320df8111edbc64f153b4456d25; __sdcfduid=401d2321df8111edbc64f153b4456d25ae69b7e45a6cdad1af635a2f1d400d5fea5d94c0414ba8b920a6f1594e802770; __cfruid=8f2fb00d9db926a38c342922f243d3124ce7060b-1682092298; __cf_bm=0FJRjS_g2Cx9e0Ib0CX3ieILMvbvW3WeItGq3fLCwO0-1682092301-0-AWS0oEw3w1XugpgmdMjYnNbi3qvje+ptSSiGzrLxbO40wmJAVrkTM5hfhB1mGIDbrsasRyxpoIEUL9v3zLyaKqRQvn3ZfOkA35mhdsIWT+E2"},
		"accept":             {"*/*"},
		"accept-language":    {"pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7"},
		"authorization":      {"MjQzNDg2NDM2NzYwNzQ4MDMz.GnBDlr.XJRVWC2rffg2VW6OqxunFYMJE3GISgbRJRtmhQ"},
		"referer":            {"https://discord.com/channels/862273703709900822/868574854536888401"},
		"referrer-policy":    {"strict-origin-when-cross-origin"},
		"sec-ch-ua":          {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"x-debug-options":    {"bugReporterEnabled"},
		"x-discord-locale":   {"pl"},
		"x-super-properties": {"eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiQ2hyb21lIiwiZGV2aWNlIjoiIiwic3lzdGVtX2xvY2FsZSI6InBsLVBMIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzExMi4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiYnJvd3Nlcl92ZXJzaW9uIjoiMTEyLjAuMC4wIiwib3NfdmVyc2lvbiI6IjEwIiwicmVmZXJyZXIiOiIiLCJyZWZlcnJpbmdfZG9tYWluIjoiIiwicmVmZXJyZXJfY3VycmVudCI6IiIsInJlZmVycmluZ19kb21haW5fY3VycmVudCI6IiIsInJlbGVhc2VfY2hhbm5lbCI6InN0YWJsZSIsImNsaWVudF9idWlsZF9udW1iZXIiOjE5MTAyNiwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbCwiZGVzaWduX2lkIjowfQ=="},
	}

	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
		DiscordMonitorGold(users)
		return
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
		DiscordMonitorGold(users)
		return
	}

	var monitorDiscordGoldCodesStruct GoldenCodesDiscord
	err = json.Unmarshal(bodyText, &monitorDiscordGoldCodesStruct)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error while monitoring discord gold codes: %v.", err))
		DiscordMonitorGold(users)
		return
	}
	if !firstPing {
		if monitorDiscordGoldCodesStruct[0].Content != lastDiscordContent {
			Log(Logger, logrus.WarnLevel, fmt.Sprintf("Found new free gold voucher: %v. Sending tasks.", monitorDiscordGoldCodesStruct[0].Content))
			for _, user := range users["usernames"] {
				Sleep(randomIntFromInterval(150, 250))
				go EnterGold(monitorDiscordGoldCodesStruct[0].Content, user)
				if err != nil {
					Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
					DiscordMonitorGold(users)
					return
				}
			}

		}
	}
	lastDiscordContent = monitorDiscordGoldCodesStruct[0].Content
	Sleep(randomIntFromInterval(5000, 9000))
	firstPing = false
	DiscordMonitorGold(users)
}

func EnterGold(promoCode string, user Users) {

	var data = strings.NewReader(fmt.Sprintf(`{"promoCode":"%s","recaptcha":"null"}`, promoCode))
	var options []tls_client.HttpClientOption = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_108),
	}

	if !proxyLess {
		options = append(options, tls_client.WithProxyUrl(user.ProxyURL))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://key-drop.com/pl/Api/activation_code", data)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
		return
	}

	req.Header = http.Header{
		"content-type":     {"application/json"},
		"cookie":           {user.Cookies},
		"user-agent":       {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"x-requested-with": {"XMLHttpRequest"},
	}
	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
			return
		}

		var adddingCodesStruct AddingCodes
		err = json.Unmarshal(bodyText, &adddingCodesStruct)
		if err != nil {
			Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error: %v.", err))
			return
		}

		if adddingCodesStruct.ErrorCode == "spamError" {
			Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error while entering discord gold codes: %v. Data: %v", user.Name, adddingCodesStruct.Info))
			Sleep(5000)
			EnterGold(promoCode, user)
		} else if adddingCodesStruct.ErrorCode == "expiredCode" {
			Log(Logger, logrus.ErrorLevel, fmt.Sprintf("Error, discount code expired: %v.", user.Name))

		} else if adddingCodesStruct.ErrorCode == `usedCode` {
			Log(Logger, logrus.InfoLevel, fmt.Sprintf("Code already used for acc: %v. Data: %v", user.Name, adddingCodesStruct))

		} else if adddingCodesStruct.Title == "Z\u0142oty kod zosta\u0142 aktywowany" || adddingCodesStruct.Title == "The golden code has been activated" {
			Log(Logger, logrus.InfoLevel, fmt.Sprintf("Sucessfully added gold code for: %v. Data: %v", user.Name, adddingCodesStruct))
		} else {
			Log(Logger, logrus.WarnLevel, fmt.Sprintf("Other response: %v. Data: %v", user.Name, adddingCodesStruct))
		}
	} else {
		Log(Logger, logrus.WarnLevel, fmt.Sprintf("Bad request for user: %v, Request status: %v", user.Name, resp.StatusCode))
	}

}

type GoldenCodesDiscord []struct {
	ID        string `json:"id"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
	ChannelID string `json:"channel_id"`
	Author    struct {
		ID               string      `json:"id"`
		Username         string      `json:"username"`
		GlobalName       interface{} `json:"global_name"`
		DisplayName      interface{} `json:"display_name"`
		Avatar           string      `json:"avatar"`
		Discriminator    string      `json:"discriminator"`
		PublicFlags      int         `json:"public_flags"`
		Bot              bool        `json:"bot"`
		AvatarDecoration interface{} `json:"avatar_decoration"`
	} `json:"author"`
	Attachments     []interface{} `json:"attachments"`
	Embeds          []interface{} `json:"embeds"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	Pinned          bool          `json:"pinned"`
	MentionEveryone bool          `json:"mention_everyone"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Flags           int           `json:"flags"`
	Components      []interface{} `json:"components"`
	Reactions       []struct {
		Emoji struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"emoji"`
		Count        int `json:"count"`
		CountDetails struct {
			Burst  int `json:"burst"`
			Normal int `json:"normal"`
		} `json:"count_details"`
		BurstColors []interface{} `json:"burst_colors"`
		MeBurst     bool          `json:"me_burst"`
		Me          bool          `json:"me"`
	} `json:"reactions,omitempty"`
}

type AutoGenerated struct {
	PromoCode    string      `json:"promoCode"`
	Status       bool        `json:"status"`
	Bonus        interface{} `json:"bonus"`
	GoldBonus    string      `json:"goldBonus"`
	DepositBonus interface{} `json:"depositBonus"`
	History      []struct {
		Title        string `json:"title"`
		Type         string `json:"type"`
		PromoCode    string `json:"promoCode"`
		GoldBonus    string `json:"goldBonus,omitempty"`
		Date         string `json:"date"`
		DepositBonus string `json:"depositBonus,omitempty"`
		Bonus        struct {
			Value    float64 `json:"value"`
			Currency string  `json:"currency"`
		} `json:"bonus,omitempty"`
	} `json:"history"`
	Title string `json:"title"`
}

type AddingCodes struct {
	Status    bool   `json:"status"`
	Info      string `json:"info"`
	ErrorCode string `json:"errorCode"`
	Title     string `json:"title"`
}
