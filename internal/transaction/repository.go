package transaction

import "eth_fetcher/internal/model"

type StoreFinder interface {
	Storer
	Finder
}

type Storer interface {
	Store(transactions *model.Transaction) error
	StoreHashesPerUser(user,hash string) error
}

type Finder interface {
	FindAll() ([]*model.Transaction, error)
	FindByHash(hash string) (*model.Transaction, error)
	FindTransactionHashesPerUser(user string) ([]string, error)
}
