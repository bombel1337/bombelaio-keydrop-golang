package main

import (
	"bombelaio-keydrop-golang/utils"
	"fmt"
	"os"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)
const asciiLogo string =`___  ____  __  ______  ______     
/ _ )/ __ \/  |/  / _ )/ __/ /     
/ _  / /_/ / /|_/ / _  / _// /__    
/____/\____/_/  /_/____/___/____/`

func init() {
	utils.Logger = logrus.New()
	utils.Logger.Formatter = &utils.CustomFormatter{}
    utils.Logger.SetOutput(colorable.NewColorableStdout())
	txtfile, err := utils.EnsureCaptchaKey()
	if err != nil {
		utils.Log(utils.Logger, logrus.ErrorLevel,  fmt.Sprintf("Error reading txt file: %v.", err))
		
	}
	if len(txtfile) == 0 {
		utils.Log(utils.Logger, logrus.ErrorLevel, "Error reading captcha from file.", )
		utils.Sleep(1500)	
		os.Exit(0)
	} else if len(txtfile) == 1 {
		utils.CaptchaKey = txtfile[0]
	} else if  len(txtfile) == 2 {
		utils.IsWebhookEnabled = true
		utils.CaptchaKey = txtfile[0]
		utils.DiscordWebhook = txtfile[1]
	} else {
		utils.Log(utils.Logger, logrus.ErrorLevel, "Wrong data in txt file.")
	}
}

func main() {
	var Option string
	var Raffletype string


	err := utils.EnsureDataFile()
	if err != nil {
		utils.Log(utils.Logger, logrus.WarnLevel,fmt.Sprintf("Error: %v", err))
    }
	utils.EnsureProxyFile();
	if err != nil {
		utils.Log(utils.Logger, logrus.WarnLevel,fmt.Sprintf("Error: %v", err))
    }

	utils.Log(utils.Logger, logrus.WarnLevel,fmt.Sprintf("\n\n %s \n\n", asciiLogo))
	utils.Log(utils.Logger, logrus.WarnLevel, "Select option. \n1. KeyDrop raffle enter.")

	fmt.Scan(&Option)
	switch Option {
	case "1":
		utils.Log(utils.Logger, logrus.WarnLevel, " \n   What type of raffles? \n   1. Champion\n   2. Challenger\n   3. Legend\n   4. Contender\n   5. Amateur")


		fmt.Scan(&Raffletype)

		switch Raffletype {
		case "1":
			utils.Log(utils.Logger, logrus.WarnLevel, "Good choice... getting data.")


			utils.Sleep(500)
			utils.ReadDataCsv("champion") 

		case "2":
			utils.Log(utils.Logger, logrus.WarnLevel, "Good choice... getting data.")

			utils.Sleep(500)
			utils.ReadDataCsv("challenger") 

		case "3":
			utils.Log(utils.Logger, logrus.WarnLevel, "Good choice... getting data.")

			utils.Sleep(500)
			utils.ReadDataCsv("legend") 

		case "4":
			utils.Log(utils.Logger, logrus.WarnLevel, "Good choice... getting data.")

			utils.Sleep(500)
			utils.ReadDataCsv("contender") 

		case "5":
			utils.Log(utils.Logger, logrus.WarnLevel, "Good choice... getting data.")

			utils.Sleep(500)
			utils.ReadDataCsv("amateur") 
			
		default:
			utils.Log(utils.Logger, logrus.ErrorLevel, "Option doesnt't exist.")

			utils.Sleep(2000)
			fmt.Print("\033[H\033[2J")
			main()
		}

	default:
		utils.Log(utils.Logger, logrus.ErrorLevel, "Option doesnt't exist.")

		utils.Sleep(2000)
		fmt.Print("\033[H\033[2J")
		main()
	}

}