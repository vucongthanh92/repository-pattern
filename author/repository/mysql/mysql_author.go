package mysql

import (
	"context"
	"database/sql"

	models "demo-echo/model"
)

type mysqlAuthorRepo struct {
	DB *sql.DB
}

// NewMysqlAuthorRepository func
func NewMysqlAuthorRepository(db *sql.DB) models.AuthorRepository {
	return &mysqlAuthorRepo{
		DB: db,
	}
}

func (m *mysqlAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (res models.Author, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return models.Author{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return res, nil
}

func (m *mysqlAuthorRepo) GetByID(ctx context.Context, id int64) (models.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
