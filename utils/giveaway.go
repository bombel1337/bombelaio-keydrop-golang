package utils

import (
	"fmt"
	"encoding/json"
	"io"
	"time"
	"net/url"
	// "bytes"
	"strings"
	"log"
    "math/rand"
	
	"bombelaio-keydrop-golang/models"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func init() {
	Logger = logrus.New()
	Logger.Formatter = &CustomFormatter{}
    Logger.SetOutput(colorable.NewColorableStdout())
}


func GettingLoggedIn(cookiesData string, raffleType string , integerUser int) {
	userNumber := fmt.Sprintf("%03d", integerUser)
	var proxyURL string
	var options []tls_client.HttpClientOption = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_108),
	}
	
	if !proxyLess {
		var randomProxy = proxyList[integerUser - 1]
		proxyArr := strings.Split(randomProxy, ":")
		proxyURL = fmt.Sprintf("http://%s:%s@%s:%s", proxyArr[2], proxyArr[3], proxyArr[0], proxyArr[1])
		options = append(options, tls_client.WithProxyUrl(proxyURL))
	} 


	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
		return
	}

	req, err := http.NewRequest(http.MethodGet, "https://key-drop.com/pl/apiData/Init/index", nil)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
		return
	}

	req.Header = http.Header{
		"cookie" : {cookiesData},
		"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
	}

	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
		return
	}
	defer resp.Body.Close()
	

    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
			return
		}

		var loggedInStruct models.GettingLoggedInStruct
		err = json.Unmarshal(bodyBytes, &loggedInStruct)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
			return
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
		Sleep(5000)
		GettingLoggedIn(cookiesData, raffleType, integerUser)
		
	}
}

func openFreeChest(index int, user Users){
	userNumber := fmt.Sprintf("%03d", index)
	postData := url.Values{}
	postData.Add("level", "0")
	var options []tls_client.HttpClientOption = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_108),
	}
	
	if !proxyLess {
		options = append(options, tls_client.WithProxyUrl(user.ProxyURL))
	} 


	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
	}


	req, err := http.NewRequest(http.MethodPost, "https://key-drop.com/pl/apiData/DailyFree/open", strings.NewReader(postData.Encode()))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header = http.Header{
		"content-type": {"application/x-www-form-urlencoded"},
		"cookie": {user.Cookies},
		"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
	}

	resp, err := client.Do(req)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error opening free chest: %v.", err))
		openFreeChest(index, user)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Error opening free chest: %v.",userNumber, err))
			openFreeChest(index, user)
			return
		}
		var freeCaseStruct models.FreeCaseStruct
		err = json.Unmarshal(bodyBytes, &freeCaseStruct)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Error opening free chest: %v.",userNumber, err))
			openFreeChest(index, user)
			return
		}

		if freeCaseStruct.Status {
			if freeCaseStruct.WinnerData.PrizeValue.Title == "" {
				Log(Logger, logrus.InfoLevel,  fmt.Sprintf("[%v] Opened free chest: %v , %v",userNumber,freeCaseStruct.WinnerData.PrizeValue.Title, freeCaseStruct.WinnerData.PrizeValue.Subtitle))
			} else {
				Log(Logger, logrus.InfoLevel,  fmt.Sprintf("[%v] Opened free chest: %v , %v",userNumber,freeCaseStruct, freeCaseStruct))
			}
			Sleep(1000 * 60 * 60 * 24)
			openFreeChest(index, user)
			} else {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Can't open free case yet: %v",userNumber,freeCaseStruct))
			Sleep(1000 * 60 * 60 * 24)
			openFreeChest(index, user)
		}

	} else if resp.StatusCode >= 500 {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%v] Error opening free chest: %v.",userNumber,resp.StatusCode))
		openFreeChest(index, user)
		return
	}
}


func monitoringGiveaway(raffleType string) {
	for index, user := range users["usernames"] {
		go openFreeChest(index, user)
	}
	go DiscordMonitorGold(users)
		var retriesInteger int = 0
		prevGiveawayID := ""

		req, err := http.NewRequest(http.MethodGet, "https://wss-2061.key-drop.com/v1/giveaway//list?type=active&page=0&perPage=5&status=active&sort=latest", nil)
		if err != nil {
			log.Println(err)
			return
		}
	
		req.Header = http.Header{
			"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		}

		for {
			var proxy string

			// jar := tls_client.NewCookieJar()
			options := []tls_client.HttpClientOption{
				tls_client.WithTimeoutSeconds(30),
				tls_client.WithClientProfile(tls_client.Chrome_112),
				// tls_client.WithCookieJar(jar), 
			}

			if !proxyLess {
				randomIndex := rand.Intn(len(proxyList))
				randomProxy := proxyList[randomIndex]
				proxyArr := strings.Split(randomProxy, ":")
				proxy = fmt.Sprintf("http://%s:%s@%s:%s", proxyArr[2], proxyArr[3], proxyArr[0], proxyArr[1])
				options = append(options, tls_client.WithProxyUrl(proxy))
			}

			client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
			if err != nil {
				log.Println(err)
				return
			}
			
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
					Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error unmarshal giveaway: %v.", err))
					monitoringGiveaway(raffleType)
					return
				}
				var totalPrice float64 = 0.0

				for i := 0; i < len(giveawayStruct.Data); i++ {
					if giveawayStruct.Data[i].Frequency == raffleType && prevGiveawayID != giveawayStruct.Data[i].ID && giveawayStruct.Data[i].ParticipantCount != 1000 {
						for _, prize := range giveawayStruct.Data[i].Prizes {		
							totalPrice += prize.Price
						}

						if (totalPrice > 2) {
							Log(Logger, logrus.WarnLevel,  fmt.Sprintf("Found new giveaway: %s, sending tasks! Value: %v", giveawayStruct.Data[i].ID, totalPrice))
							for index, user := range users["usernames"] {							
								go gettingBearer(raffleType, giveawayStruct.Data[i].ID, user, index, retriesInteger)
								if err != nil {
									Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error: %v.", err))
									monitoringGiveaway(raffleType)
									return
								}
								
							}
						} else {
							Log(Logger, logrus.WarnLevel,  fmt.Sprintf("Found new giveaway: %s, but prize is too low, not sending tasks! Value: %v", giveawayStruct.Data[i].ID, totalPrice))
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

	var options []tls_client.HttpClientOption = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_108),
	}
	
	if !proxyLess {
		options = append(options, tls_client.WithProxyUrl(user.ProxyURL))
	} 


	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
	}

	req, err := http.NewRequest(http.MethodGet, "https://key-drop.com/token?t=" + fmt.Sprint(time.Now().Unix()), nil)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header = http.Header{
		"cookie" : {user.Cookies},
		"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
	}

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
		joinGiveaway(user.Cookies, raffleType, giveawayID, string(body), user, false, index, retriesInteger)
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
	var options []tls_client.HttpClientOption = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_108),
	}
	
	if !proxyLess {
		options = append(options, tls_client.WithProxyUrl(user.ProxyURL))
	} 


	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error: %v.", userNumber ,err))
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


	url := "https://wss-3002.key-drop.com/v1/giveaway//joinGiveaway/" + fmt.Sprint(giveawayID)
	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error setting request joingiveaway: %v.",userNumber, err))
		return
	}

	req.Header = http.Header{
		"cookie" : {user.Cookies},
		"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
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
	url := "https://wss-2061.key-drop.com/v1/giveaway//data/" + fmt.Sprint(giveawayID)

	var proxy string

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_112),
	}

	if !proxyLess {
		randomIndex := rand.Intn(len(proxyList))
		randomProxy := proxyList[randomIndex]
		proxyArr := strings.Split(randomProxy, ":")
		proxy = fmt.Sprintf("http://%s:%s@%s:%s", proxyArr[2], proxyArr[3], proxyArr[0], proxyArr[1])
		options = append(options, tls_client.WithProxyUrl(proxy))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header = http.Header{
		"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
	}

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


