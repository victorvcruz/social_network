package repository

import (
	"database/sql"
	"social_network_project/database/postgresql"
	"social_network_project/entities"
	"social_network_project/entities/response"
	"strings"
	"time"
)

type PostRepository interface {
	InsertPost(post *entities.Post) error
	FindPostsByAccountID(id *string) ([]interface{}, error)
	UpdatePostDataByID(postID, accountID, content *string) error
	FindPostByID(id *string) (*response.PostResponse, error)
	ExistsPostByID(id *string) (*bool, error)
	RemovePostByID(postID, accountID *string) error
	ExistsPostByPostIDAndAccountID(postID, accountID *string) (*bool, error)
	CountInteractionsForPost(postID *string, typeValue int) (*int, error)
}

type PostRepositoryStruct struct {
	Db *sql.DB
}

func NewPostRepository() PostRepository {
	return &PostRepositoryStruct{postgresql.Db}
}

func (p *PostRepositoryStruct) InsertPost(post *entities.Post) error {
	sqlStatement := `
		INSERT INTO post (id, account_id, content, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6)`

	row := p.Db.QueryRow(sqlStatement, post.ID, post.AccountID, post.Content, post.CreatedAt,
		post.UpdatedAt, post.Removed)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *PostRepositoryStruct) FindPostsByAccountID(id *string) ([]interface{}, error) {
	sqlStatement := `
		SELECT id, account_id, content, created_at, updated_at
		FROM post
		WHERE account_id = $1
		AND removed = false`

	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var post response.PostResponse
	for rows.Next() {
		err = rows.Scan(
			&post.ID,
			&post.AccountID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = strings.Join(strings.Split(post.CreatedAt, "T00:00:00Z"), "")
		post.UpdatedAt = strings.Join(strings.Split(post.CreatedAt, "T00:00:00Z"), "")

		like, err := p.CountInteractionsForPost(&post.ID, response.INTERACTION_TYPE_LIKED.EnumIndex())
		if err != nil {
			return nil, err
		}

		dislike, err := p.CountInteractionsForPost(&post.ID, response.INTERACTION_TYPE_DISLIKED.EnumIndex())
		if err != nil {
			return nil, err
		}

		post.Like = *like
		post.Dislike = *dislike
		list = append(list, post)
	}

	return list, nil
}

func (p *PostRepositoryStruct) UpdatePostDataByID(postID, accountID, content *string) error {
	sqlStatement := `
		UPDATE post 
		SET content = $1, updated_at = $2
		WHERE id = $3
		AND account_id = $4  
		AND removed = false`

	updateTime := time.Now().UTC().Format("2006-01-02")

	row := p.Db.QueryRow(sqlStatement, content, updateTime, postID, accountID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *PostRepositoryStruct) FindPostByID(id *string) (*response.PostResponse, error) {
	sqlStatement := `
		SELECT id, account_id, content, created_at, updated_at
		FROM post
		WHERE id = $1
		AND removed = false`

	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var post response.PostResponse
	err = rows.Scan(
		&post.ID,
		&post.AccountID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepositoryStruct) ExistsPostByID(id *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM post
		WHERE id = $1
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *PostRepositoryStruct) RemovePostByID(postID, accountID *string) error {
	sqlStatement := `
		UPDATE post 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	row := p.Db.QueryRow(sqlStatement, postID, accountID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *PostRepositoryStruct) ExistsPostByPostIDAndAccountID(postID, accountID *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM post
		WHERE id = $1
		AND account_id = $2
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, postID, accountID)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *PostRepositoryStruct) CountInteractionsForPost(postID *string, typeValue int) (*int, error) {
	sqlStatement := `
		SELECT count(type) 
		FROM interaction
		WHERE post_id = $1
		AND type = $2`

	rows, err := p.Db.Query(sqlStatement, postID, typeValue)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var count *int
	err = rows.Scan(&count)
	if err != nil {
		return nil, err
	}

	return count, nil
}
