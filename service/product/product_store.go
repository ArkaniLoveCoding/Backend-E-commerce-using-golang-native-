package product

import (
	"context"
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

func (s *Store) GetAllProduct([]types.Products) (*types.Products, error) {

	query := `SELECT name, stock, price, expired, category, created_at, updated_at FROM product_clients`
	var products []types.Products

	rows, err := s.db.Queryx(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Faield to get all product!")
		}
		return nil, nil
	}

	for rows.Next() {
		if err := rows.StructScan(&products); err != nil {
			return nil, errors.New("Failed to scan data!")
		}
	}

	return nil, nil
}

func (s *Store) GetProductByID(id string) (*types.Products, error) {

	var products types.Products
	query := 
	`
	SELECT id, name, stock, price, expired, category FROM product_clients
	WHERE id = $1
	`

	if err := s.db.Get(&products, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to find the id!")
		}
		return nil, nil
	}

	return nil, nil

}

func (s *Store) CreateNewProduct(ctx context.Context, products *types.Products) error {

	query := `
	INSERT INTO product_clients (id, name, stock, price, expired, category, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING *
	`

	if err := s.db.QueryRowContext(
		ctx, 
		query,
		products.Id,
		products.Name,
		products.Stock,
		products.Price,
		products.Expired,
		products.Category,
		products.Created_at,
		products.Updated_at,
	).Scan(
		&products.Id,
		&products.Name,
		&products.Stock,
		&products.Price,
		&products.Expired,
		&products.Category,
		&products.Created_at,
		&products.Updated_at,
	); err != nil {
		return errors.New("Failed to scan query and store!" + err.Error())
	}

	return nil

} 