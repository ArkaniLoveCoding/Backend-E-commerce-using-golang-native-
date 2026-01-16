package product

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
)

type Store struct {
	db *sqlx.DB
}

func NewStore (db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAllProduct ([]types.Products) (*[]types.Products, error) {

	query := `SELECT name, stock, price, expired, category FROM products`
	var products []types.Products

	if err := s.db.Select(products, query); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No one data is matching with the result!")
		}
		return nil, errors.New("Failed to load products")
	}

	return nil, nil
}