package content

import (
	"database/sql"
	"fmt"
	"log"
	"txp/restapistarter/internal/app/module/content/entity"
	sqlUtil "txp/restapistarter/pkg/data/sql"
)

type ContentRepository struct {
	db *sql.DB
}

func NewContentRepository(db *sql.DB) *ContentRepository {
	r := new(ContentRepository)
	r.db = db
	return r
}

func (r *ContentRepository) Create(e *entity.Content) error {
	_, err := r.db.Exec(
		"INSERT INTO contents (name)"+
			"VALUES ($1)",
		e.Name,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *ContentRepository) ReadMany() (*sql.Rows, error) {
	rows, err := r.db.Query(
		"SELECT * FROM contents", // WHERE id IS NOT NULL
	)
	if err != nil {
		return nil, fmt.Errorf("ReadMany %v", err)
	}
	return rows, nil
}

func (r *ContentRepository) ReadOne(id string) *sql.Row {
	row := r.db.QueryRow(
		"SELECT * FROM contents WHERE id = $1 LIMIT 1",
		id,
	)
	return row
}

func (r *ContentRepository) Update(id string, e *entity.Content) (int64, error) {
	q := "UPDATE contents SET name = $2 WHERE id = $1"
	res, err := r.db.Exec(
		q,
		id,
		e.Name,
	)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return sqlUtil.GetRowsAffected(res), nil
}

func (r *ContentRepository) Delete(id string) (int64, error) {
	q := "DELETE FROM contents WHERE id = $1"
	res, err := r.db.Exec(
		q,
		id,
	)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return sqlUtil.GetRowsAffected(res), nil
}
