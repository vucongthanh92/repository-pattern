package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"demo-echo/article/repository"
	models "demo-echo/model"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

// NewMysqlArticleRepository func
func NewMysqlArticleRepository(conn *sql.DB) models.ArticleRepository {
	return &mysqlArticleRepository{conn}
}

func (m *mysqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []models.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]models.Article, 0)
	for rows.Next() {
		t := models.Article{}
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = models.Author{
			ID: authorID,
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *mysqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []models.Article, nextCursor string, err error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
		FROM article ORDER BY created_at LIMIT ? `
	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}
	fmt.Println("decodedCursor:", decodedCursor)
	res, err = m.fetch(ctx, query, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}
	return
}

func (m *mysqlArticleRepository) GetByID(ctx context.Context, id int64) (res models.Article, err error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
		  FROM article WHERE ID = ?`
	listAr, err := m.fetch(ctx, query, id)
	if err != nil {
		return models.Article{}, err
	}

	if len(listAr) > 0 {
		res = listAr[0]
	} else {
		return res, models.ErrNotFound
	}
	return res, nil
}

func (m *mysqlArticleRepository) GetByTitle(ctx context.Context, title string) (res models.Article, err error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
		  FROM article WHERE title = ?`
	listAr, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(listAr) > 0 {
		res = listAr[0]
	} else {
		return res, models.ErrNotFound
	}
	return res, nil
}

func (m *mysqlArticleRepository) Store(ctx context.Context, ar *models.Article) (err error) {
	query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	smtp, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := smtp.ExecContext(ctx, ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.CreatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	ar.ID = lastID
	return
}

func (m *mysqlArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM article WHERE id = ?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rows != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rows)
		return
	}
	return
}

func (m *mysqlArticleRepository) Update(ctx context.Context, ar *models.Article) (err error) {
	query := `UPDATE article set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.CreatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}
	return
}
