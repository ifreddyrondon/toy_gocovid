package domain

import (
	"sort"
	"strings"
)

type Stats struct {
	Updated                int64   `json:"updated"`
	Cases                  int64   `json:"cases"`
	TodayCases             int64   `json:"todayCases"`
	Deaths                 int64   `json:"deaths"`
	TodayDeaths            int64   `json:"todayDeaths"`
	Recovered              int64   `json:"recovered"`
	Active                 int64   `json:"active"`
	Critical               int64   `json:"critical"`
	CasesPerOneMillion     float64 `json:"casesPerOneMillion"`
	DeathsPerOneMillion    float64 `json:"deathsPerOneMillion"`
	Tests                  int64   `json:"tests"`
	TestsPerOneMillion     float64 `json:"testsPerOneMillion"`
	Population             int64   `json:"population"`
	ActivePerOneMillion    float64 `json:"activePerOneMillion"`
	RecoveredPerOneMillion float64 `json:"recoveredPerOneMillion"`
	CriticalPerOneMillion  float64 `json:"criticalPerOneMillion"`
}

// All represents up to date Global totals
type All struct {
	Stats
	AffectedCountries int `json:"affectedCountries"`
}

type Country struct {
	Stats
	Continent   string `json:"continent"`
	Country     string `json:"country"`
	CountryInfo struct {
		ID   int     `json:"_id"`
		Iso2 string  `json:"iso2"`
		Iso3 string  `json:"iso3"`
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
		Flag string  `json:"flag"`
	} `json:"countryInfo"`
}

type FilterFunc func([]Country) []Country

func (f FilterFunc) Filter(countries []Country) []Country {
	return f(countries)
}

// SortByCases sorts the data by number of total cases
var FilterByName = func(prefix string) FilterFunc {
	return func(countries []Country) []Country {
		return filter(countries, func(country Country) bool {
			lowerString, lowerPrefix := strings.ToLower(country.Country), strings.ToLower(prefix)
			return strings.HasPrefix(lowerString, lowerPrefix)
		})
	}
}

func filter(arr []Country, cond func(Country) bool) []Country {
	var result []Country
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}

type SortFunc func([]Country)

func (f SortFunc) Sort(countries []Country) {
	f(countries)
}

// SortByCases sorts the data by number of total cases
var SortByCases = SortFunc(func(countries []Country) {
	sort.SliceStable(countries, func(i, j int) bool {
		return countries[i].Cases > countries[j].Cases
	})
})

// SortByCasesToday sorts the data by number of cases today
var SortByCasesToday = SortFunc(func(countries []Country) {
	sort.SliceStable(countries, func(i, j int) bool {
		return countries[i].TodayCases > countries[j].TodayCases
	})
})

// SortByDeaths sorts the data by number of total deaths
var SortByDeaths = SortFunc(func(countries []Country) {
	sort.SliceStable(countries, func(i, j int) bool {
		return countries[i].Deaths > countries[j].Deaths
	})
})

// SortByDeathsToday sorts the data by number of deaths today
var SortByDeathsToday = SortFunc(func(countries []Country) {
	sort.SliceStable(countries, func(i, j int) bool {
		return countries[i].TodayDeaths > countries[j].TodayDeaths
	})
})

// SortByActive sorts the data by number of active cases
var SortByActive = SortFunc(func(countries []Country) {
	sort.SliceStable(countries, func(i, j int) bool {
		return countries[i].Active > countries[j].Active
	})
})

// SortByRecoveries sorts the data by number of total recoveries
var SortByRecoveries = SortFunc(func(countries []Country) {
	sort.SliceStable(countries, func(i, j int) bool {
		return countries[i].Recovered > countries[j].Recovered
	})
})

// SortByMortality sorts the data by mortality rate
var SortByMortality = SortFunc(func(countries []Country) {
	sort.SliceStable(countries, func(i, j int) bool {
		return float64(countries[i].Deaths)/float64(countries[i].Cases) > float64(countries[j].Deaths)/float64(countries[j].Cases)
	})
})
