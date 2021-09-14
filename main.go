/*
Copyright © 2021 Steffen Rumpf <github@steffen-rumpf.de>

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
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily

		fmt.Fprintf(w, "Usage of %s: \n", os.Args[0])
		fmt.Fprintln(w, `
This tool can be used to download the raw covid data from the german RKI GitHub repository.
Per default it will just print the fully aggregated overall sums per age group. Multiple flags
can be used to change this behavior:

Examples:
covid -age-group A00-A04   - Just print the aggregated sum for age group A00-A04
covid -region 8222         - Just print the results for region 8222 (Mannheim)
covid -aggregate-month     - Instead calculating the full sum aggregate sums per month
covid -aggregate-year      - Instead calculating the full sum aggregate sums per year
covid -update              - Download the latest data from GitHub

Full Example:
covid -region 8222 -aggregate-month -update

Flags:`)

		flag.PrintDefaults()
	}
}

func main() {
	var err error
	flag.Parse()

	if help {
		flag.CommandLine.Usage()
	} else {
		var fi os.FileInfo
		if updateData {
			if _, err = os.Stat(tmpFile); err == nil {
				err := os.Remove(tmpFile)
				logIfFatal(err)
			}
		}
		if fi, err = os.Stat(tmpFile); err != nil {
			downloadRKIrawData()
		} else {
			fmt.Printf("Downloaded %s at %s\n\n", tmpFile, fi.ModTime())
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
	} else {
		allCases(csvReader)
	}
}

func casesPerYear(csvReader *csv.Reader) {
	countCases := map[string]map[string]int{}
	header := []string{}
	countCases["SUM"] = map[string]int{}

	i := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		if i > 0 &&
			(regionCode == "" || record[regionCodeID] == regionCode) &&
			(ageGroup == "" || record[ageGroupID] == ageGroup) {
			key1 := strings.Split(record[dateID], "-")[0]
			key2 := strings.TrimSpace(record[ageGroupID])

			if !contains(header, key2) {
				header = append(header, key2)
			}

			count, err := strconv.Atoi(record[countID])
			logIfFatal(err)
			if countCases[key1] == nil {
				countCases[key1] = map[string]int{}
			}
			countCases[key1][key2] += count
			countCases["SUM"][key2] += count
		}
		i++
	}

	printSimpleTable(countCases, header)
}

func casesPerMonth(csvReader *csv.Reader) {
	countCases := map[string]map[string]int{}
	header := []string{}
	countCases["SUM"] = map[string]int{}

	i := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		if i > 0 &&
			(regionCode == "" || record[regionCodeID] == regionCode) &&
			(ageGroup == "" || record[ageGroupID] == ageGroup) {
			date := strings.Split(record[dateID], "-")
			year := date[0]
			month := date[1]
			key1 := fmt.Sprintf("%s-%s", year, month)
			key2 := strings.TrimSpace(record[ageGroupID])

			if !contains(header, key2) {
				header = append(header, key2)
			}

			count, err := strconv.Atoi(record[countID])
			logIfFatal(err)
			if countCases[key1] == nil {
				countCases[key1] = map[string]int{}
			}
			countCases[key1][key2] += count
			countCases["SUM"][key2] += count
		}
		i++
	}
	printSimpleTable(countCases, header)
}

func allCases(csvReader *csv.Reader) {
	countCases := map[string]int{}
	i := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		if i > 0 &&
			(regionCode == "" || record[regionCodeID] == regionCode) &&
			(ageGroup == "" || record[ageGroupID] == ageGroup) {
			key1 := record[ageGroupID]

			count, err := strconv.Atoi(record[countID])
			logIfFatal(err)
			countCases[key1] += count
			countCases["SUM"] += count
		}
		i++
	}

	keys := make([]string, 0, len(countCases))
	maxKeyLen := 0
	for k := range countCases {
		keys = append(keys, k)
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}
	sort.Strings(keys)
	headerSpace := " "
	if maxKeyLen-8 > 0 {
		headerSpace = strings.Repeat(" ", maxKeyLen-6)
	}
	fmt.Printf("AgeGroup%sCount\n", headerSpace)
	fmt.Println(strings.Repeat("-", maxKeyLen+10))
	for _, key := range keys {
		space := strings.Repeat(" ", maxKeyLen-len(key)+2)
		fmt.Printf("%s%s%d\n", key, space, countCases[key])
	}
}

func printSimpleTable(countCases map[string]map[string]int, header []string) {
	tableHeader := "AgeGroup" + strings.Repeat(" ", 2)
	columWidth := 10
	sort.Strings(header)

	for _, hd := range header {
		tableHeader = tableHeader + hd + strings.Repeat(" ", columWidth-len(hd))
	}
	fmt.Println(tableHeader)
	fmt.Println(strings.Repeat("-", len(tableHeader)+2))

	keys := make([]string, 0, len(countCases))
	for k := range countCases {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s%s", key, strings.Repeat(" ", columWidth-len(key)))

		ageGroups := countCases[key]
		for _, ag := range header {
			fmt.Printf("%d%s", ageGroups[ag], strings.Repeat(" ", columWidth-len(strconv.Itoa(ageGroups[ag]))))
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
