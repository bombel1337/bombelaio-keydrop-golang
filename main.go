package main


import (
    "fmt"
	"bombelaio-keydrop-golang/utils"
)


func main() {
	err := utils.EnsureDataFile()
    if err != nil {
		fmt.Printf("Error: %s", err)
    }
	fmt.Println("Hello, world!")
}