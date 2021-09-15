package domain

type HistoricalCountry struct {
	Country  string   `json:"country"`
	Province []string `json:"province"`
	Timeline Timeline `json:"timeline"`
}

type Timeline struct {
	Cases     map[string]int64 `json:"cases"`
	Deaths    map[string]int64 `json:"deaths"`
	Recovered map[string]int64 `json:"recovered"`
}
