package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type DataFile struct {
	Config  EauConfig  `json:"config"`
	Entries []EauEntry `json:"entries"`
}

type EauConfig struct {
	DailyIntake int `json:"dailyIntake"`
}

type EauEntry struct {
	Date   int64 `json:"date"`
	Amount int   `json:"amount"`
}

const dataFileName string = ".eau"

var dataFilePath string

func initDataFile() {
	// TODO: Support for Windows
	home := os.Getenv("HOME")
	dataFilePath = path.Join(home, dataFileName)

	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		// TODO: This doesn't look so pretty
		data := []byte("{\"config\":{\"dailyIntake\":3700},\"entries\":[]}")
		err := ioutil.WriteFile(dataFilePath, data, 0644)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Data file initialized")
	}
}

func getSum(entries []EauEntry, date string) (sum int) {
	for _, entry := range entries {
		if date == strings.Split(time.Unix(entry.Date, 0).String(), " ")[0] {
			sum = sum + entry.Amount
		}
	}

	return
}

func addEntry(quantity int) {
	// TODO: Handle errors

	jsonFile, _ := os.Open(dataFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var dataFile DataFile
	json.Unmarshal(byteValue, &dataFile)

	dataFile.Entries = append(dataFile.Entries, EauEntry{
		time.Now().Unix(),
		quantity,
	})

	newByteValue, _ := json.Marshal(dataFile)
	ioutil.WriteFile(dataFilePath, newByteValue, 0644)

	today := strings.Split(time.Now().String(), " ")[0]
	fmt.Printf("[Eau] New entry added. %d/%d\n", getSum(dataFile.Entries, today), dataFile.Config.DailyIntake)
}

func printStatus() {
	// TODO: Handle errors

	// TODO: Move readfile to helper func
	jsonFile, _ := os.Open(dataFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var dataFile DataFile
	json.Unmarshal(byteValue, &dataFile)

	today := strings.Split(time.Now().String(), " ")[0]
	fmt.Printf("[Eau] Today: %d/%d\n", getSum(dataFile.Entries, today), dataFile.Config.DailyIntake)

}

func main() {
	initDataFile()

	quantity := flag.Int("a", 0, "Add quantity in milliliters")
	flag.Parse()

	if *quantity > 0 {
		addEntry(*quantity)
		return
	}

	printStatus()
}
