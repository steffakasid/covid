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
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/steffakasid/covid/internal"
)

var (
	regionCode, ageGroup, yearFilter                string
	help, aggregateYear, aggregateMonth, updateData bool
)

func init() {
	flag.BoolVarP(&help, "help", "h", false, "Show help")
	flag.StringVarP(&regionCode, "region", "r", "", "German region code e.g. 8222 for Mannheim")
	flag.StringVarP(&ageGroup, "age-group", "a", "", "Age group")
	flag.StringVarP(&yearFilter, "year", "y", "", "Filter the result by year")
	flag.BoolVar(&aggregateMonth, "aggregate-month", false, "Aggregate cases by month")
	flag.BoolVar(&aggregateYear, "aggregate-year", false, "Aggregate cases by year")
	flag.BoolVarP(&updateData, "update", "u", false, "Update data")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: \n", os.Args[0])
		fmt.Fprintln(os.Stderr, `
This tool can be used to download the raw covid data from the german RKI GitHub repository.
Per default it will just print the fully aggregated overall sums per age group. Multiple flags
can be used to change this behavior:

Examples:
covid --age-group A00-A04   - Just print the aggregated sum for age group A00-A04
covid --region 8222         - Just print the results for region 8222 (Mannheim)
covid --aggregate-month     - Instead calculating the full sum aggregate sums per month
covid --aggregate-year      - Instead calculating the full sum aggregate sums per year
covid --update              - Download the latest data from GitHub
covid --year 2021           - Only show results for 2021

Full Example:
covid -region 8222 -aggregate-month -update

Flags:`)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
	} else {
		covid := internal.GetInstance("https://media.githubusercontent.com/media/robert-koch-institut/SARS-CoV-2_Infektionen_in_Deutschland/master/Aktuell_Deutschland_SarsCov2_Infektionen.csv", "/tmp/Aktuell_Deutschland_SarsCov2_Infektionen.csv")

		if updateData {
			covid.UpdateData()
		}

		var aggregate internal.Aggregation
		if aggregateMonth && aggregateYear {
			panic("You should only define either --aggregate-year or --aggregate-month not both")
		} else if aggregateMonth {
			aggregate = internal.AggregateMonth
		} else if aggregateYear {
			aggregate = internal.AggregateYear
		} else {
			aggregate = internal.NoAggregation
		}

		data, header := covid.ParseData(regionCode, ageGroup, yearFilter, aggregate)
		printSimpleTable(data, header)
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
