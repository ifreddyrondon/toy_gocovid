package fetching

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/gocovid/private/domain"
)

type Countries struct {
	Doer
}

func (f *Countries) Fetch() ([]domain.Country, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://corona.lmao.ninja/v2/countries", nil)
	res, err := f.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	countries := make([]domain.Country, 0)
	if err := json.NewDecoder(res.Body).Decode(&countries); err != nil {
		return nil, err
	}
	return countries, err
}
