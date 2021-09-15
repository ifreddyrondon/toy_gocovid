package fetching_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/gocovid/private/fetching"
)

var (
	req200 = func() ClientFunc {
		return func(r *http.Request) (*http.Response, error) {
			json := `{"updated":1590006540224,"cases":5054252,"todayCases":71315,"deaths":327938,"todayDeaths":3384,"recovered":2005840,"active":2720474,"critical":45240,"casesPerOneMillion":648,"deathsPerOneMillion":42.1,"tests":64979271,"testsPerOneMillion":8385.55,"population":7748960018,"activePerOneMillion":351.08,"recoveredPerOneMillion":258.85,"criticalPerOneMillion":5.84,"affectedCountries":215}}`
			body := ioutil.NopCloser(strings.NewReader(json))
			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		}
	}
	req500 = func() ClientFunc {
		return func(r *http.Request) (*http.Response, error) {
			json := `}`
			body := ioutil.NopCloser(strings.NewReader(json))
			return &http.Response{
				StatusCode: 500,
				Body:       body,
			}, nil
		}
	}
	reqWithErr = func() ClientFunc {
		return func(r *http.Request) (*http.Response, error) {
			return nil, errors.New("test")
		}
	}
)

type ClientFunc func(r *http.Request) (*http.Response, error)

func (c ClientFunc) Do(req *http.Request) (*http.Response, error) {
	return c(req)
}

func TestFetcherAll_FetchOK(t *testing.T) {
	t.Parallel()
	fetcher := fetching.All{Doer: req200()}
	s, err := fetcher.Fetch()
	assert.Nil(t, err)
	assert.Equal(t, int64(5054252), s.Cases)
	assert.Equal(t, int64(71315), s.TodayCases)
}

func TestFetcherAll_FetchWithErrorForRequest(t *testing.T) {
	t.Parallel()
	fetcher := fetching.All{Doer: reqWithErr()}
	s, err := fetcher.Fetch()
	assert.EqualError(t, err, "test")
	assert.Nil(t, s)
}

func TestFetcherAll_FetchWithErrorForResponse(t *testing.T) {
	t.Parallel()
	fetcher := fetching.All{Doer: req500()}
	s, err := fetcher.Fetch()
	assert.EqualError(t, err, "invalid character '}' looking for beginning of value")
	assert.Nil(t, s)
}
