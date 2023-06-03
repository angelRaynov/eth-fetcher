package transaction

import "eth_fetcher/internal/model"

type StoreFetcher interface {
	Storer
}

type Storer interface {
	Store(transactions []model.Transaction)
}
