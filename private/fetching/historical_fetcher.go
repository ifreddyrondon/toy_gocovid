package fetching

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/gocovid/private/domain"
)

type HistoricalCountryVal struct {
	HistoricalCountry domain.HistoricalCountry
	Err               error
}

const uriTemplate = "https://disease.sh/v3/covid-19/historical/%v?lastdays=%v"
const defaultLastDays = 30

type Historical struct {
	Doer
}

func (f *Historical) FetchCountryChan(iso string) <-chan HistoricalCountryVal {
	ch := make(chan HistoricalCountryVal, 1)
	go func() {
		v, err := f.fetchCountry(iso)
		if err != nil {
			ch <- HistoricalCountryVal{Err: err}
			return
		}
		ch <- HistoricalCountryVal{HistoricalCountry: *v}
	}()
	return ch
}

func (f *Historical) fetchCountry(iso string) (*domain.HistoricalCountry, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(uriTemplate, iso, defaultLastDays), nil)
	res, err := f.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var historical domain.HistoricalCountry
	if err := json.NewDecoder(res.Body).Decode(&historical); err != nil {
		return nil, err
	}
	return &historical, nil
}
