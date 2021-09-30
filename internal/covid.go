package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	regionCodeID = 0
	ageGroupID   = 1
	dateID       = 3
	countID      = 9
	countDeadID  = 10
)

type Aggregation int

const (
	NoAggregation Aggregation = iota
	AggregateYear
	AggregateMonth
)

type Covid struct {
	rkiRawDataURL string
	tmpFile       string
}

func GetInstance(rawDataUrl, tmpFile string) Covid {
	return Covid{rawDataUrl, tmpFile}
}

func (c Covid) ParseData(regionCode, ageGroup, yearFilter string, aggregate Aggregation) (map[string]map[string]int, []string) {

	c.initData()

	file, err := os.Open(c.tmpFile)
	logIfFatal(err)

	csvReader := csv.NewReader(file)

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

			key1 := ""
			key2 := strings.TrimSpace(record[ageGroupID])
			date := strings.Split(record[dateID], "-")
			year := date[0]
			month := date[1]

			if aggregate == AggregateYear {
				key1 = year
			}
			if aggregate == AggregateMonth {
				key1 = fmt.Sprintf("%s-%s", year, month)
			}

			if yearFilter == "" || yearFilter == year {
				if !contains(header, key2) {
					header = append(header, key2)
				}

				count, err := strconv.Atoi(record[countID])
				logIfFatal(err)
				if key1 != "" {
					if countCases[key1] == nil {
						countCases[key1] = map[string]int{}
					}
					countCases[key1][key2] += count
				}
				countCases["SUM"][key2] += count
			}
		}
		i++
	}
	return countCases, header
}

func (c Covid) initData() {
	if fi, err := os.Stat(c.tmpFile); err != nil {
		c.UpdateData()
	} else {
		fmt.Printf("Downloaded %s at %s\n\n", c.tmpFile, fi.ModTime())
	}
}

func (c Covid) UpdateData() {
	if _, err := os.Stat(c.tmpFile); err == nil {
		logIfFatal(os.Remove(c.tmpFile))
	}
	logIfFatal(c.downloadRKIrawData())
}

func (c Covid) downloadRKIrawData() error {
	// Get the data
	resp, err := http.Get(c.rkiRawDataURL)
	logIfFatal(err)
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(c.tmpFile)
	logIfFatal(err)

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
