package fetching

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/gocovid/private/domain"
)

type All struct {
	Doer
}

func (f *All) Fetch() (*domain.All, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://disease.sh/v3/covid-19/all", nil)
	res, err := f.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var all domain.All
	if err := json.NewDecoder(res.Body).Decode(&all); err != nil {
		return nil, err
	}
	return &all, nil
}
