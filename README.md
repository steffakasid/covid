# Covid [![Go](https://github.com/steffakasid/covid/actions/workflows/go.yml/badge.svg)](https://github.com/steffakasid/covid/actions/workflows/go.yml)

This tool just downloads the RKI covid raw data from link:https://media.githubusercontent.com/media/robert-koch-institut/SARS-CoV-2_Infektionen_in_Deutschland/master/Aktuell_Deutschland_SarsCov2_Infektionen.csv[] and shows them on the command line.

*Example Output: `covid --aggregate-month --region 8222`*
```
AgeGroup  A00-A04   A05-A14   A15-A34   A35-A59   A60-A79   A80+      unbekannt 
----------------------------------------------------------------------------------
2020-03   1         3         87        108       24        4         0         
2020-04   2         7         70        90        38        11        0         
2020-05   3         5         13        12        4         1         0         
2020-06   2         1         11        6         1         1         0         
2020-07   3         10        17        22        3         0         0         
2020-08   2         12        108       76        12        1         0         
2020-09   19        44        147       113       20        4         0         
2020-10   25        124       483       453       121       27        3         
2020-11   58        209       884       943       347       230       7         
2020-12   75        297       980       1135      464       314       2         
2021-01   52        101       521       566       225       179       0         
2021-02   33        63        247       286       89        26        0         
2021-03   117       224       639       684       192       29        0         
2021-04   92        306       750       920       234       40        1         
2021-05   56        169       401       463       119       23        0         
2021-06   5         29        81        68        15        4         0         
2021-07   2         23        123       34        8         5         0         
2021-08   18        95        416       282       39        9         0         
2021-09   77        299       639       539       101       23        0         
SUM       642       2021      6617      6800      2056      931       13  
```

# Flags
`covid --help`
```
Usage of covid: 

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

Flags:
  -a, --age-group string   Age group
      --aggregate-month    Aggregate cases by month
      --aggregate-year     Aggregate cases by year
  -h, --help               Show help
  -r, --region string      German region code e.g. 8222 for Mannheim
  -u, --update             Update data
  -y, --year string        Filter the result by year
```