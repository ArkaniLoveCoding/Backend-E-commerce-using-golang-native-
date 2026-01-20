package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (s *Store) UpdateToken(id uuid.UUID, token string, token_refresh string) error {

	query := `
		UPDATE users 
		SET token = $2,
			refresh_token = $3
		WHERE id = $1
	`

	_, err := s.store.DB.Exec(query, id, token, token_refresh)
	if err != nil {
		return errors.New("Failed to update token!" + err.Error())
	}

	return nil

}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	query := `SELECT 
	id, firstname, lastname, password, email, country, address, role, token, refresh_token, created_at
	FROM users WHERE email = $1`
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
	err := s.store.Get(&user, "SELECT id, firstname, lastname, password, email, created_at FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil

}

func (s *Store) CreateUser(ctx context.Context, user *types.User) error  {
	
	query := `
		INSERT INTO users (id, firstname, lastname, 
		password, email, country, address, role, token, refresh_token, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING *
	`

	if err := s.store.QueryRowContext(
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
		user.CreatedAt,
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
		&user.CreatedAt,
		); err != nil {
		return nil
	}

	return nil
}

func (s *Store) GetAllUser([]types.User) (*[]types.User, error) {

	query := `SELECT * FROM users`
	var users []types.User

	rows, err := s.store.Queryx(query, users)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to get all data of the user!" + err.Error())
		}
		return nil, errors.New(err.Error())
	}

	for rows.Next() {
		if err := rows.StructScan(&users); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("Failed to scan all data of the user!" + err.Error())
			}
			return nil, errors.New(err.Error())
		}
	}

	return &users, nil

}
