package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

	query := `SELECT id, name, stock, image, price, expired, category, created_at, updated_at FROM product_clients`

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
	SELECT id, name, stock, image, price, expired, category FROM product_clients
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

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New(err.Error())
	}

	defer tx.Rollback()

	query := `
	INSERT INTO product_clients (id, name, stock, image, price, expired, category, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING *
	`

	if err := tx.QueryRowContext(
		ctx, 
		query,
		products.Id,
		products.Name,
		products.Stock,
		products.Image,
		products.Price,
		products.Expired,
		products.Category,
		products.Created_at,
		products.Updated_at,
	).Scan(
		&products.Id,
		&products.Name,
		&products.Stock,
		&products.Image,
		&products.Price,
		&products.Expired,
		&products.Category,
		&products.Created_at,
		&products.Updated_at,
	); err != nil {
		return errors.New("Failed to scan query and store!" + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to do commit on this seasons!")
	}

	return nil

} 

func (s *Store) DeleteProductsOnlyAdmin(id uuid.UUID, ctx context.Context) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelLinearizable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing some transactions!")
	}

	defer tx.Rollback()

	query := `
		DELETE FROM product_clients WHERE id = $1;
	`

	result, err := tx.ExecContext(ctx, query, id)
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

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to comit some transactions!")
	}

	return nil

}

func (s *Store) UpdateProductsOnlyAdmin(
	id uuid.UUID, 
	name string,
	stock int,
	image string,
	category string,
	expired string,
	price string,
	ctx_update context.Context,
) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelLinearizable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing some transactions!")
	}

	query := `
		UPDATE product_clients 
		SET name = $2,
			price = $3,
			stock = $4,
			image = $5,
			category = $6,
			expired = $7
		WHERE id = $1
		RETURNING *;
	`

	result, err := tx.ExecContext(
		ctx_update, 
		query,
		id,
		name,
		price,
		stock,
		image,
		category,
		expired,
	)
	if err != nil {
		return errors.New("Failed to execute the update query!")
	}

	result_affected, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to scan the rows affected in your db!")
	}

	if result_affected == 0 {
		return errors.New("no one data is changing in your db!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to comit some transactions!")
	}

	return nil

}

func (s *Store) SearchManyProducts(ctx context.Context, keyword string, offset int, limit int) ([]types.Products, error) {

	query := `
		SELECT id, name, price, stock, category, expired, 
		FROM product_clients WHERE name LIKE ?
		LIMIT = ? OFFSET = ?
	`

	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	search_query := "%" + keyword + "%"

	rows, err := s.store.QueryContext(ctx, query, search_query, offset, limit)
	if err != nil {
		return nil, errors.New("Failed to get data from db!")
	}
	defer rows.Close()

	var products []types.Products

	for rows.Next() {
		var p types.Products
		if err := rows.Scan(&p); err != nil {
			return nil, errors.New("Failed to get data from db!")
		}
		products = append(products, p)
	}

	return nil, nil

}
