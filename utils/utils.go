package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"io"
    "math/rand"
	"strings"
	"sync"
	"github.com/sirupsen/logrus"
)

var proxyList []string // global variable to store the array
var proxyLess bool
var CaptchaKey string 
var IsWebhookEnabled bool
var DiscordWebhook string 




type CustomFormatter struct{}


type Users struct {
    Name string
	SteamID string
    Tries  int
	Avatar string
	Wins int
	ProxyURL string
	Cookies string
}

// global object containing arrays of users
var users = map[string][]Users{}
var Logger *logrus.Logger


func Log(logger *logrus.Logger, level logrus.Level, message string) {
    logger.WithFields(logrus.Fields{}).Log(level, message)
}

func AddUserToArray(arrayName string, newUser Users) {
    users[arrayName] = append(users[arrayName], newUser)
}

func UpdateUserTries(userName string) {
    for arrayName, userArray := range users {
        for i, user := range userArray {
            if user.Name == userName {
                // found the user, update their tries
                users[arrayName][i].Tries++
                return
            }
        }
    }
}

func UpdateUserWins(userName string) {
    for arrayName, userArray := range users {
        for i, user := range userArray {
            if user.Name == userName {
                // found the user, update their tries
                users[arrayName][i].Wins++
                return
            }
        }
    }
}


func Sleep(ms int) {
    time.Sleep(time.Duration(ms) * time.Millisecond)
}

func randomIntFromInterval(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func EnsureDataFile() error {
    // Check if the file exists
    _, err := os.Stat("data.csv")
    if os.IsNotExist(err) {
        // Create the file with headers
        file, err := os.Create("data.csv")
        if err != nil {
            return err
        }
        defer file.Close()

        writer := csv.NewWriter(file)
        defer writer.Flush()
        err = writer.Write([]string{"Proxies", "Cookies"})
        if err != nil {
            return err
        }
    } else if err != nil {
        return err
    }

    return nil
}

func ReadDataCsv(raffleType string) {
	file, err := os.Open("data.csv")
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error opening data file %v.", err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if (len(rows)==2){
		Log(Logger, logrus.WarnLevel,  fmt.Sprintf("Detected only %v account, restoring session.", len(rows) -1))

	} else if (len(rows)>2){
		Log(Logger, logrus.WarnLevel,  fmt.Sprintf("Detected wow a lot tbh: %v accounts, restoring sessions.", len(rows) -1))
	}

	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error reading data file %v.", err))
	}
	var wg sync.WaitGroup

	for i := 1; i < len(rows); i++ {
		wg.Add(1)

		go func(row []string, i int) {
			defer wg.Done()
			GettingLoggedIn(row[1], raffleType, i)
			Sleep(250)
		}(rows[i], i)
	}
	wg.Wait()

	monitoringGiveaway(raffleType)
}

func EnsureProxyFile() {
	file, err := os.Open("data.csv")
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error opening data file: %v.", err))
		return
	}
	defer file.Close()
	
	r := csv.NewReader(file)
	headers, err := r.Read()
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error reading data file row: %v.", err))
		return
	}
	
	proxiesIndex := -1
	for i, header := range headers {
		if header == "Proxies" {
			proxiesIndex = i
			break
		}
	}
	if proxiesIndex == -1 {
		Log(Logger, logrus.ErrorLevel,  "Error: no \"Proxies\" column found in data file")
		return
	}
	
	// Read each row and extract the proxies
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error reading data file row: %v.", err))
			return
		}
		if len(record) > proxiesIndex {
			proxyList = append(proxyList, record[proxiesIndex])
		}
	}
	
	if len(proxyList) == 0 {
		proxyLess = true
	}
}

func EnsureCaptchaKey()  ([]string, error) {
	file, err := os.Open("captcha_key.txt")
	if err != nil {
		Log(Logger, logrus.ErrorLevel,  fmt.Sprintf("Error reading txt file: %v.", err))
		return nil, err
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the contents of the file into a string variable
	var parts []string
	for scanner.Scan() {
		// Split the line by comma and append the resulting slice to the parts slice
		lineParts := strings.Split(strings.TrimSpace(scanner.Text()), ",")
		parts = append(parts, lineParts...)
	}

	// Split the first line by comma and return the resulting slice
	return parts, nil
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    var color int
    switch entry.Level {
    case logrus.InfoLevel:
        color = 32 // Green
    case logrus.WarnLevel:
        color = 34 // Yellow
    case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
        color = 31 // Red
    default:
        color = 0 // Default color
    }

    message := fmt.Sprintf("[%s] %s\n", entry.Time.Format("15:04:05"), entry.Message)
    return []byte("\x1b[" + fmt.Sprintf("%d", color) + "m" + message + "\x1b[0m"), nil
}