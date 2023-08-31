package utils


import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"bombelaio-keydrop-golang/models"	
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
	"net/url"
)

func gettingCaptchaCapmonster(giveawayID string, user Users, index int) (string, error) {
	userNumber := fmt.Sprintf("%03d", index)

	var client *http.Client
	var proxyURL string = user.ProxyURL
	if !proxyLess {
		urlProxy, err := url.Parse(proxyURL)
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("[%s] Error parsing proxy joingiveaway: %v.",userNumber, err))
			return "err", err
		}
		client = &http.Client {
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
		}
	} else {
		client = &http.Client{}
	}
	data := strings.NewReader(fmt.Sprintf(`{"clientKey": "%s", "task": {"type": "RecaptchaV2EnterpriseTaskProxyless", "websiteURL": "https://key-drop.com/pl/giveaways/keydrop/%s", "websiteKey": "6Ld2uggaAAAAAG9YRZYZkIhCdS38FZYpY9RRYkwN"}}`, CaptchaKey, giveawayID))
	req, err := http.NewRequest("POST", "https://api.capmonster.cloud/createTask", data)
    if err != nil {
        return "err", err
    }
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 OPR/97.0.0.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
        return "err", err
	}
	defer resp.Body.Close()


	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "err", err
		}

		// Parse the JSON response
		var captchaTaskIDStruck models.CaptchaTaskID
		err = json.Unmarshal(bodyBytes, &captchaTaskIDStruck)
		if err != nil {
			return "err", err
		}

	if captchaTaskIDStruck.ErrorID == 0 && captchaTaskIDStruck.TaskID != 0 {
		
		jsonData, err := json.Marshal(map[string]interface{}{
			"clientKey": CaptchaKey,
			"taskId": captchaTaskIDStruck.TaskID,
		})
		if err != nil {
			return "err", err
		}
		
		req, err := http.NewRequest("POST", "https://api.capmonster.cloud/getTaskResult", bytes.NewBuffer(jsonData))
		if err != nil {
			return "err", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 OPR/97.0.0.0")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		for {
			resp, err := client.Do(req)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {

				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return "err", err
				}
		
				var CaptchaResultStruct models.CaptchaResult
				err = json.Unmarshal(bodyBytes, &CaptchaResultStruct)
				if err != nil {
					return "err", err
				}

				if CaptchaResultStruct.Status == "ready" {
					return CaptchaResultStruct.Solution.GRecaptchaResponse, nil
				}

			} else {
				return strconv.Itoa(resp.StatusCode), err
			}
	

			Sleep(2000)
		}

	} else if captchaTaskIDStruck.ErrorID == 1 {
		return captchaTaskIDStruck.ErrorCode, err
	}

	} else {
		return strconv.Itoa(resp.StatusCode), err
	}

	return "err", err
}
