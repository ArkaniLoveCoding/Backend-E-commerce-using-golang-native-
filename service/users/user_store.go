package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
)

type Store struct {
	store *sqlx.DB
}

func NewStore(store *sqlx.DB) *Store {
	return &Store{store: store}
}

func (s *Store) UpdateToken(
	ctx context.Context, id uuid.UUID, token string, token_refresh string, user *types.User) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
	ctx_tx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx_tx, tx_options)
	if err != nil {
		return errors.New("Failed to doing transactions!")
	}

	defer tx.Rollback()

	query := `
		UPDATE users 
		SET token = $2,
			refresh_token = $3
		WHERE id = $1;
	`

	if err := tx.QueryRowContext(
		ctx,
		query,
		id,
		token,
		token_refresh,
	).Scan(
		&user.Id,
		&user.Token,
		&user.Rerfresh_token,
	); err != nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transactions!")
	}

	return nil

}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	query := `SELECT 
	id, firstname, lastname, password, email, country, address, role, token, refresh_token, created_at, updated_at
	FROM users WHERE email = $1;`
	err := s.store.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (s *Store) GetUserById(id uuid.UUID) (*types.User, error) {

	var user types.User
	err := s.store.Get(&user, "SELECT id, firstname, lastname, password, email, country, address, role, token, refresh_token, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil

}

func (s *Store) CreateUser(ctx context.Context, user *types.User) error  {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing transactions!")
	}

	defer tx.Rollback()
	
	query := `
		INSERT INTO users (id, firstname, lastname, 
		password, email, country, address, role, token, refresh_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING *;
	`

	if err := tx.QueryRowContext(
		ctx,
		query,
		user.Id,
		user.Firstname,
		user.Lastname,
		user.Password,
		user.Email,
		user.Country,
		user.Address,
		user.Role,
		user.Token,
		user.Rerfresh_token,
		user.Created_at,
		user.Updated_at,
	).Scan(
		&user.Id,
		&user.Firstname,
		&user.Lastname,
		&user.Password,
		&user.Email,
		&user.Country,
		&user.Address,
		&user.Role,
		&user.Token,
		&user.Rerfresh_token,
		&user.Created_at,
		&user.Updated_at,
		); err != nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil
}

func (s *Store) GetAllUser() ([]types.User, error) {

	query := `
	SELECT firstname, lastname, password, email, country, address, token, refresh_token, created_at, updated_at
	FROM users;
	`
	var users []types.User

	rows, err := s.store.Queryx(query, users)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to get all data of the user!" + err.Error())
		}
		return nil, errors.New(err.Error())
	}

	for rows.Next() {
		var u types.User
		if err := rows.StructScan(&u); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("Failed to scan all data of the user!" + err.Error())
			}
			return nil, errors.New(err.Error())
		}
		users = append(users, u)
	}
	defer rows.Close()

	return users, nil

}

func (s *Store) GetUsersRole(role string) (*types.User, error) {

	var users types.User
	query := `
		SELECT id, firstname, lastname, password, email, country, address, role, token, refresh_token
		created_at FROM users WHERE role = $1;
	`
	if err := s.store.Get(&users, query, role); err != nil {
		return nil, errors.New("Cannot find the user role!")
	}

	return nil, nil

}

func (s *Store) UpdateDataUser(
	id uuid.UUID, 
	ctx context.Context, 
	firstname string,
	lastname string,
	password string,
	email string,
	country string,
	address string,
	users *types.User,
	) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing transactions!")
	}

	defer tx.Rollback()

	query := `
		UPDATE users 
		SET firstname = $2,
			lastname = $3,
			password = $4,
			email = $5,
			country = $6,
			address = $7
		WHERE id = $1;
	`
	var u = users

	if err := tx.QueryRowContext(
		ctx,
		query,
		id,
		firstname,
		lastname,
		password,
		email,
		country,
		address,
	).Scan(
		&u.Id,
		&u.Firstname,
		&u.Lastname,
		&u.Password,
		&u.Email,
		&u.Country,
		&u.Address,
	); err != nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transactions!")
	}

	return nil

}

func (s *Store) DeleteUsersOnlyAdmin(id uuid.UUID, ctx context.Context) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing the transactions!")
	}

	defer tx.Rollback()

	query := `
		DELETE FROM users WHERE id = $1; 
	`
	var users types.User

	result, err := tx.ExecContext(ctx, query, users.Id)
	if err != nil {
		return errors.New("Failed to execute the context from db" + err.Error())
	}
	
	rows_affected, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to get the rows from db, no one be execute from your db!")
	}

	if rows_affected == 0 {
		return errors.New("Failed to checking the rows affected from your db!")
	}

	return nil

}

