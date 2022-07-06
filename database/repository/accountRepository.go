package repository

import (
	"database/sql"
	"social_network_project/entities"
	"strings"
)

type AccountRepository struct {
	Db *sql.DB
}

func (p *AccountRepository) InsertAccount(account *entities.Account) error {
	sqlStatement := `
		INSERT INTO account (id, username, name, description, email, password, created_at, updated_at, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := p.Db.Exec(sqlStatement, account.ID, account.Username, account.Name, account.Description,
		account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.Deleted)
	if err != nil {
		return err
	}

	return nil
}

func (p *AccountRepository) ExistsAccountByEmailAndPassword(email string, password string) (*bool, error) {
	sqlStatement := `
		SELECT email, password 
		FROM account
		WHERE email = $1
		AND password = $2`
	rows, err := p.Db.Query(sqlStatement, email, password)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *AccountRepository) FindAccountIDbyEmail(email string) (*string, error) {
	sqlStatement := `
		SELECT id
		FROM account
		WHERE email = $1`
	rows, err := p.Db.Query(sqlStatement, email)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var id *string
	_ = rows.Scan(&id)

	return id, nil
}

func (p *AccountRepository) FindAccountbyID(id string) (*entities.Account, error) {
	sqlStatement := `
		SELECT id, username, name, description, email, password, created_at, updated_at, deleted
		FROM account
		WHERE id = $1`
	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var account entities.Account
	err = rows.Scan(
		&account.ID,
		&account.Username,
		&account.Name,
		&account.Description,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Deleted,
	)
	if err != nil {
		return nil, err
	}

	account.CreatedAt = strings.Join(strings.Split(account.CreatedAt, "T00:00:00Z"), "")
	account.UpdatedAt = strings.Join(strings.Split(account.CreatedAt, "T00:00:00Z"), "")

	return &account, nil
}