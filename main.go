/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	rkiRawDataURL = "https://media.githubusercontent.com/media/robert-koch-institut/SARS-CoV-2_Infektionen_in_Deutschland/master/Aktuell_Deutschland_SarsCov2_Infektionen.csv"
	tmpFile       = "/tmp/Aktuell_Deutschland_SarsCov2_Infektionen.csv"
)

var (
	regionCode                                      string
	ageGroup                                        string
	help, aggregateYear, aggregateMonth, updateData bool
)

const (
	regionCodeID = 0
	ageGroupID   = 1
	dateID       = 3
	countID      = 9
	countDeadID  = 10
)

func init() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.StringVar(&regionCode, "region", "", "German region code e.g. 8222 for Mannheim")
	flag.StringVar(&ageGroup, "age-group", "", "Age group")
	flag.BoolVar(&aggregateMonth, "aggregate-month", false, "Aggregate cases by month")
	flag.BoolVar(&aggregateYear, "aggregate-year", false, "Aggregate cases by year")
	flag.BoolVar(&updateData, "update", false, "Update data")
}

func main() {
	flag.Parse()

	if help {
		flag.CommandLine.Usage()
	} else {
		if updateData {
			err := os.Remove(tmpFile)
			logIfFatal(err)
		}
		if _, err := os.Stat(tmpFile); err != nil {
			downloadRKIrawData()
		}
		parseData()
	}
}

func parseData() {
	file, err := os.Open(tmpFile)
	logIfFatal(err)

	csvReader := csv.NewReader(file)

	if aggregateYear {
		casesPerYear(csvReader)
	} else if aggregateMonth {
		casesPerMonth(csvReader)
	} else if ageGroup != "" {
		allCases(csvReader)
	} else {
		casesPerAge(csvReader)
	}
}

func casesPerYear(csvReader *csv.Reader) {
	countCases := map[string]map[string]int{}
	header := []string{}
	i := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		if i > 0 && (regionCode == "" || record[regionCodeID] == regionCode) {
			key1 := strings.Split(record[dateID], "-")[0]
			key2 := record[ageGroupID]

			if !contains(header, key2) {
				header = append(header, key2)
			}

			count, err := strconv.Atoi(record[countID])
			logIfFatal(err)
			if countCases[key1] == nil {
				countCases[key1] = map[string]int{}
			}
			countCases[key1][key2] += count
		}
		i++
	}

	printSimpleTable(countCases, header)
}

func casesPerMonth(csvReader *csv.Reader) {
	countCases := map[string]map[string]int{}
	header := []string{}
	i := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		if i > 0 && (regionCode == "" || record[regionCodeID] == regionCode) {
			date := strings.Split(record[dateID], "-")
			year := date[0]
			month := date[1]
			key1 := fmt.Sprintf("%s-%s", year, month)
			key2 := record[ageGroupID]

			if !contains(header, key2) {
				header = append(header, key2)
			}

			count, err := strconv.Atoi(record[countID])
			logIfFatal(err)
			if countCases[key1] == nil {
				countCases[key1] = map[string]int{}
			}
			countCases[key1][key2] += count
		}
		i++
	}
	printSimpleTable(countCases, header)
}

func allCases(csvReader *csv.Reader) {
	countCases := 0
	i := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		if i > 0 && (regionCode == "" || record[regionCodeID] == regionCode && record[ageGroupID] == ageGroup) {
			count, err := strconv.Atoi(record[countID])
			logIfFatal(err)
			countCases += count
		}
		i++
	}
	fmt.Printf("Cases in %s: %d\n", ageGroup, countCases)
}

func casesPerAge(csvReader *csv.Reader) {
	countCase := map[string]int{}
	i := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		if i > 0 && (regionCode == "" || record[regionCodeID] == regionCode) {
			count, err := strconv.Atoi(record[countID])
			logIfFatal(err)
			countCase[record[ageGroupID]] += count
		}
		i++
	}

	keys := make([]string, 0, len(countCase))
	maxKeyLen := 0
	for k := range countCase {
		keys = append(keys, k)
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}
	sort.Strings(keys)
	fmt.Println("AgeGroup", "\t", "Count")
	fmt.Println("-----------------------")
	for _, key := range keys {
		div := maxKeyLen / len(key)
		tabs := strings.Repeat("\t", div)
		fmt.Println(key, tabs, countCase[key])
	}
}

func printSimpleTable(countCases map[string]map[string]int, header []string) {
	tableHeader := "AgeGroup\t"

	sort.Strings(header)

	for _, hd := range header {
		tableHeader = tableHeader + fmt.Sprintf("%s\t", hd)
	}
	fmt.Println(tableHeader)
	tabs := strings.Count(tableHeader, "\t")
	for i := 0; i < (len(tableHeader) + tabs + 3); i++ {
		fmt.Print("-")
	}
	fmt.Println()

	keys := make([]string, 0, len(countCases))
	for k := range countCases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {

		fmt.Printf("%s\t\t", key)

		ageGroups := countCases[key]
		for _, ag := range header {
			fmt.Printf("%d\t", ageGroups[ag])
		}
		fmt.Println()
	}
}

func downloadRKIrawData() error {
	// Get the data
	resp, err := http.Get(rkiRawDataURL)
	logIfFatal(err)
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(tmpFile)
	logIfFatal(err)

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func logIfFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func contains(arr []string, val string) bool {
	for _, s := range arr {
		if s == val {
			return true
		}
	}
	return false
}
