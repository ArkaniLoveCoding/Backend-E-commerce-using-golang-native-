package product

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
)

type Store struct {
	store *sqlx.DB
}

func NewStoreProduct (store *sqlx.DB) *Store {
	return &Store{store: store}
}

func (s *Store) GetAllProduct() ([]types.Products, error) {

	query := `SELECT id, name, stock, price, expired, category, created_at, updated_at FROM product_clients`

	rows, err := s.store.Queryx(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Faield to get all product!")
		}
		return nil, nil
	}
	defer rows.Close()

	var products []types.Products

	for rows.Next() {
		var p types.Products
		if err := rows.StructScan(&p); err != nil {
			return nil, errors.New("Failed to scan data!")
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *Store) GetProductByID(id uuid.UUID) (*types.Products, error) {

	var products types.Products
	query := 
	`
	SELECT id, name, stock, price, expired, category FROM product_clients
	WHERE id = $1
	`

	if err := s.store.Get(&products, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to find the id!")
		}
		return nil, nil
	}

	return &products, nil

}

func (s *Store) CreateNewProduct(ctx context.Context, products *types.Products) error {

	query := `
	INSERT INTO product_clients (id, name, stock, price, expired, category, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING *
	`

	if err := s.store.QueryRowContext(
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

func (s *Store) DeleteProductsOnlyAdmin(id uuid.UUID, ctx context.Context) error {

	query := `
		DELETE FROM users WHERE id = $1;
	`
	var users types.User

	result, err := s.store.ExecContext(ctx, query, users.Id)
	if err != nil {
		return errors.New("Failed to get id from users db!" + err.Error())
	}

	rows_affected, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to check the rows from db!" + err.Error())
	}

	if rows_affected == 0 {
		return errors.New("No one data is executed in your db, thats why the rows is confirmed zero value!")
	}

	return nil

}