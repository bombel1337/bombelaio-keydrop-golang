package utils


import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"bombelaio-keydrop-golang/models"	
	"io"
	"strconv"
)

func gettingCaptchaCapmonster(giveawayID string) (string, error) {
    jsonStr := []byte(fmt.Sprintf(`{
        "clientKey": "%v",
        "task": {
            "type": "NoCaptchaTaskProxyless",
            "websiteURL": "https://key-drop.com/giveaways/keydrop/%s",
            "websiteKey": "6Ld2uggaAAAAAG9YRZYZkIhCdS38FZYpY9RRYkwN"
        }
    }`,CaptchaKey, giveawayID))

	req, err := http.NewRequest("POST", "https://api.capmonster.cloud/createTask", bytes.NewBuffer(jsonStr))
    if err != nil {
        return "err", err
    }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 OPR/97.0.0.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{}
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
