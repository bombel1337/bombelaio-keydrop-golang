package utils

import (
    "encoding/csv"
    "os"
)

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
        err = writer.Write([]string{"Cookies"})
        if err != nil {
            return err
        }
    } else if err != nil {
        return err
    }

    return nil
}