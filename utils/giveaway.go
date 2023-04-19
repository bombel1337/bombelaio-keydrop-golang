package utils

import (
	"bombelaio-keydrop-golang/models"
	"github.com/sirupsen/logrus"
	"github.com/mattn/go-colorable"
	"encoding/json"
	"fmt"
	//"bytes"
	"io"
	"time"
	"net/http"
	//"mime/multipart"
	"net/url"
    "strings"
)

func init() {
	Logger = logrus.New()
	Logger.Formatter = &CustomFormatter{}
    Logger.SetOutput(colorable.NewColorableStdout())
}


func GettingLoggedIn(cookiesData string, raffleType string , integerUser int) {
	userNumber := fmt.Sprintf("%03d", integerUser)

	var client *http.Client
	var proxyURL string 
	if !proxyLess {
		var randomProxy = proxyList[integerUser - 1]
		proxyArr := strings.Split(randomProxy, ":")
		proxyURL = fmt.Sprintf("http://%s:%s@%s:%s", proxyArr[2], proxyArr[3], proxyArr[0], proxyArr[1])
		urlProxy, err := url.Parse(proxyURL)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error: %s.", err))
		}

		client = &http.Client {
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest("GET", "https://key-drop.com/apiData/Init/index", nil)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
	}
    req.Header.Set("cookie", cookiesData)
    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 7.1; vivo 1716 Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
	}
	defer resp.Body.Close()
	

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
		}

		var loggedInStruct models.GettingLoggedInStruct
		err = json.Unmarshal(bodyBytes, &loggedInStruct)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
		}
		
		if len(loggedInStruct.Avatar) != 0 {
			if !proxyLess {
				AddUserToArray("usernames", Users{Name: loggedInStruct.UserName, SteamID: loggedInStruct.SteamID, Avatar: loggedInStruct.Avatar, Tries: 1, ProxyURL: proxyURL, Cookies: cookiesData})
			} else {
				AddUserToArray("usernames", Users{Name: loggedInStruct.UserName, SteamID: loggedInStruct.SteamID, Avatar: loggedInStruct.Avatar, Tries: 1, Cookies: cookiesData})
			}
			Log(Logger, logrus.InfoLevel,  fmt.Sprintf("[%s] Successfuly restored session for task.", userNumber))
		} else {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Request, logged error: %v", userNumber, loggedInStruct.Message ))
		}

	} else {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Request, logged error: %v", userNumber, resp.StatusCode))
	}
}

// func openFreeChest(index int, cookiesData string){
// 	userNumber := fmt.Sprintf("%03d", index)
// 	client := &http.Client{}
//     body := &bytes.Buffer{}
//     writer := multipart.NewWriter(body)
// 	writer.WriteField("level", "0")
// 	writer.Close()
// 	req, err := http.NewRequest("POST", "https://key-drop.com/pl/apiData/DailyFree/open", body)
// 	if err != nil {
// 		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error opening free chest: %v.", err))
// 		Sleep(2500)
// 		openFreeChest(index, cookiesData)
// 		return
// 	}
// 	req.Header.Set("authority", "key-drop.com")
// 	req.Header.Set("accept", "*/*")
// 	req.Header.Set("accept-language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")
// 	req.Header.Set("cookie", cookiesData)
// 	req.Header.Set("dnt", "1")
// 	req.Header.Set("origin", "https://key-drop.com")
// 	req.Header.Set("referer", "https://key-drop.com/pl/daily-case")
// 	req.Header.Set("sec-ch-ua", `"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`)
// 	req.Header.Set("sec-ch-ua-mobile", "?0")
// 	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
// 	req.Header.Set("sec-fetch-dest", "empty")
// 	req.Header.Set("sec-fetch-mode", "cors")
// 	req.Header.Set("sec-fetch-site", "same-origin")
// 	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 7.1; vivo 1716 Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error opening free chest: %v.", err))
// 		openFreeChest(index, cookiesData)
// 		return
// 	}
// 	defer resp.Body.Close()
// fmt.Println(resp.StatusCode)

// 	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
// 		bodyBytes, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Error opening free chest: %v.",userNumber, err))
// 			openFreeChest(index, cookiesData)
// 			return
// 		}
// 		var freeCaseStruct models.FreeCaseStruct
// 		err = json.Unmarshal(bodyBytes, &freeCaseStruct)
// 		if err != nil {
// 			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Error opening free chest: %v.",userNumber, err))
// 			openFreeChest(index, cookiesData)
// 			return
// 		}

// 		if freeCaseStruct.Status {
// 			Log(Logger, logrus.InfoLevel,  fmt.Sprintf("[%v] Opened free chest: %v , %v",userNumber,freeCaseStruct.WinnerData.PrizeValue.Title, freeCaseStruct.WinnerData.PrizeValue.Subtitle))
// 			Sleep(1000 * 60 * 60)
// 			openFreeChest(index, cookiesData)
// 		} else {
// 			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Can't open free case yet: %v",userNumber,freeCaseStruct.Error))
// 			Sleep(1000 * 60 * 60)
// 			openFreeChest(index, cookiesData)
// 		}

// 	} else if resp.StatusCode >= 500 {
// 		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Error opening free chest: %v.",userNumber,resp.StatusCode))
// 		openFreeChest(index, cookiesData)
// 		return
// 	}
// }


func monitoringGiveaway(raffleType string) {
		var retriesInteger int = 0
		prevGiveawayID := ""
		for {
			client := &http.Client{}
			req, err := http.NewRequest("GET", "https://ws-2061.key-drop.com/v1/giveaway//list?type=active&page=0&perPage=5&status=active&sort=latest", nil)
			if err != nil {
				Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error: %v.", err))
				monitoringGiveaway(raffleType)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 7.1; vivo 1716 Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")
			resp, err := client.Do(req)
			if err != nil {
				Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error: %v.", err))
				monitoringGiveaway(raffleType)
				return
			}
			

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error: %v.", err))
					monitoringGiveaway(raffleType)
					return
				}

				var giveawayStruct models.MonitoringGiveawayStruct
				err = json.Unmarshal(bodyBytes, &giveawayStruct)
				if err != nil {
					Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error: %v.", err))
					monitoringGiveaway(raffleType)
					return
				}

				for i := 0; i < len(giveawayStruct.Data); i++ {
					if giveawayStruct.Data[i].Frequency == raffleType && prevGiveawayID != giveawayStruct.Data[i].ID {
						Log(Logger, logrus.WarnLevel,  fmt.Sprintf("Found new giveaway: %s, sending tasks!", giveawayStruct.Data[i].ID))
						for index, user := range users["usernames"] {							
							go gettingBearer(raffleType, giveawayStruct.Data[i].ID, user, index, retriesInteger)
							if err != nil {
								Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error: %v.", err))
								monitoringGiveaway(raffleType)
								return
							}
							
						}

						go readWinners(prevGiveawayID, raffleType)
						prevGiveawayID = giveawayStruct.Data[i].ID
					} else if (i == len(giveawayStruct.Data)){
						Log(Logger, logrus.WarnLevel,  fmt.Sprintf("Couldn't find any matching giveaways for: : %s!", raffleType))

					}
				}

			} else {
				Log(Logger, logrus.ErrorLevel,  fmt.Sprintf( "Error monitoring giveaway: %v", resp.StatusCode))
			}
		Sleep(randomIntFromInterval(5000, 12000))	
		}
}

func gettingBearer(raffleType string, giveawayID string, user Users, index int, retriesInteger int)  {
	userNumber := fmt.Sprintf("%03d", index)
	var cookiesData string = user.Cookies
	var client *http.Client
	var proxyURL string = user.ProxyURL
	if !proxyLess {
		urlProxy, err := url.Parse(proxyURL)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error parsing proxy bearer: %v.", userNumber ,err))
			return
		}
		client = &http.Client {
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	} else {
		client = &http.Client{}
	}

	url := "https://key-drop.com/token?t=" + fmt.Sprint(time.Now().Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error setting up request bearer: %v.", userNumber ,err))
		return
	}
    req.Header.Set("cookie", cookiesData)
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 7.1; vivo 1716 Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error requesting bearer: %v.",userNumber, err))
		return
	}
	defer resp.Body.Close()
	

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        body, err := io.ReadAll(resp.Body)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error reading bearer body: %v.",userNumber, err))
			return
		}
		retriesInteger = 1
		joinGiveaway(cookiesData, raffleType, giveawayID, string(body), user, false, index, retriesInteger)
	} else if resp.StatusCode >= 500 {
		if retriesInteger<=3 {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error getting bearer, retrying: %v. Retry number: %v", userNumber , resp.StatusCode, retriesInteger))
			Sleep(500)
			gettingBearer(raffleType, giveawayID, user, index, retriesInteger)
			retriesInteger++
		} else {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error getting bearer, max tries: %v.", userNumber ,resp.StatusCode))
		}

	}
}

func joinGiveaway(cookiesData string, raffleType string, giveawayID string, bearerToken string ,user Users, iscaptcha bool, index int, retriesInteger int) {
	userNumber := fmt.Sprintf("%03d", index)

	var client *http.Client
	var proxyURL string = user.ProxyURL
	if !proxyLess {
		urlProxy, err := url.Parse(proxyURL)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error parsing proxy joingiveaway: %v.",userNumber, err))
			return
		}
		client = &http.Client {
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	} else {
		client = &http.Client{}
	}
	var payload io.Reader
	if iscaptcha {
		solution, err := gettingCaptchaCapmonster(giveawayID, user, index)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error getting solution: %v.",userNumber, err))
			return
		}
		Sleep(1000)
		payload = strings.NewReader(fmt.Sprintf(`{"captcha":"%v"}`, solution))
	} else {
		payload = nil
	}
	url := "https://ws-3002.key-drop.com/v1/giveaway//joinGiveaway/" + fmt.Sprint(giveawayID)
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error setting request joingiveaway: %v.",userNumber, err))
		return
	}

	if iscaptcha {
		req.Header.Set("content-type", "application/json")
		req.Header.Set("x-requested-with", "XMLHttpRequest")
	}

	req.Header.Set("authorization", "Bearer "+bearerToken)
	req.Header.Set("cookie", cookiesData)
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 7.1; vivo 1716 Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error sending request join giveaway: %v.",userNumber, err))
		return
	}
	defer resp.Body.Close()
	

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error reading body join giveaway: %v.",userNumber, err))
			return
		}

		var joinGiveawayStruct models.JoinGiveawayStruct
		err = json.Unmarshal(bodyBytes, &joinGiveawayStruct)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error marshal body join giveaway: %v.",userNumber, err))
			return
		}

		if (joinGiveawayStruct.Success){
			Log(Logger, logrus.InfoLevel,  fmt.Sprintf("[%v] User: %s, Successfuly joined giveaway: %v! Total entries: %d/%o",userNumber, user.Name, giveawayID, user.Wins, user.Tries))

			UpdateUserTries(user.Name)
		} else if !joinGiveawayStruct.Success && joinGiveawayStruct.Message == "captcha" {
			Log(Logger, logrus.InfoLevel,  fmt.Sprintf("[%v] User: %s, has a captcha, getting token!",userNumber, user.Name))

			Sleep(2500)
			joinGiveaway(cookiesData, raffleType, giveawayID, bearerToken, user, true, index, retriesInteger)
		} else {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] User: %s, unfortunately has got an error while joining! Error: %v",userNumber, user.Name, joinGiveawayStruct.Message))
		}

	} else if resp.StatusCode >= 500 {
		if retriesInteger <= 3 {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error joining giveaway, retrying: %v. Retry number: %v", userNumber , resp.StatusCode, retriesInteger))
			Sleep(500)
			joinGiveaway(cookiesData, raffleType, giveawayID, bearerToken ,user, iscaptcha, index, retriesInteger)
			retriesInteger++
		} else {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error joining giveaway, reached max retries: %v.", userNumber ,resp.StatusCode))
		}

	}
}


func readWinners(giveawayID string, raffleType string) {
	client := &http.Client{}
	url := "https://ws-2061.key-drop.com/v1/giveaway//data/" + fmt.Sprint(giveawayID)


	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error setting up request readwinners: %v.", err))
		return
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 7.1; vivo 1716 Build/N2G47H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error sending request readwinners: %v.", err))
		return
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error reading body readwinners: %v.", err))
		return
	}

	var winnerCheckersStruct models.WinnersChecker
	err = json.Unmarshal(bodyText, &winnerCheckersStruct)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error marshal body read winners: %v.", err))
		return
	}
	

	if winnerCheckersStruct.Data.Status == "ended" {
		for index, x := range winnerCheckersStruct.Data.Winners {
			for _, user := range users["usernames"] {
				if x.Userdata.IDSteam == user.SteamID {
					if IsWebhookEnabled {
						sendDiscordWebhook(user, winnerCheckersStruct.Data.Prizes[index], giveawayID)
					}

					UpdateUserWins(user.Name)
				}
			}
		}
	
	} else {
		Sleep(15000)
		readWinners(giveawayID, raffleType)
	}

}


