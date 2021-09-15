package fetching

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/ifreddyrondon/gocovid/private/domain"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type IndexVal struct {
	All       domain.All
	Countries []domain.Country
	Err       error
}

type Index struct {
	all       All
	countries Countries
}

func NewIndex(client Doer) *Index {
	return &Index{
		all:       All{Doer: client},
		countries: Countries{Doer: client},
	}
}

func (f *Index) FetchChan(ctx context.Context) <-chan IndexVal {
	ch := make(chan IndexVal, 1)
	go func() {
		group, _ := errgroup.WithContext(ctx)
		var all domain.All
		var countries []domain.Country
		group.Go(func() error {
			v, err := f.all.Fetch()
			all = *v
			return err
		})
		group.Go(func() error {
			v, err := f.countries.Fetch()
			countries = v
			return err
		})
		if err := group.Wait(); err != nil {
			ch <- IndexVal{Err: err}
			return
		}
		ch <- IndexVal{
			All: all, Countries: countries,
		}
	}()
	return ch
}
