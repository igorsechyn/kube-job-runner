package pq

import (
	"database/sql"
	"fmt"

	"kube-job-runner/pkg/app/config"
	"kube-job-runner/pkg/app/data"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(config config.Config) (*Store, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", config.PgUser, config.PgPassword, config.PgHost, config.PgPort, config.PgDatabase)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (store *Store) PutDocument(document data.Document) error {
	statement, err := store.db.Prepare("INSERT INTO executions(image, tag, status, jobId, timestamp) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(document.Image, document.Tag, document.Status, document.JobID, document.Timestamp)
	return err
}

func (store *Store) GetDocuments(id string) ([]data.Document, error) {
	rows, err := store.db.Query("SELECT image, tag, status, jobId, timestamp FROM executions WHERE jobId = $1 ORDER BY timestamp DESC", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]data.Document, 0)
	for rows.Next() {
		var document data.Document
		err := rows.Scan(&document.Image, &document.Tag, &document.Status, &document.JobID, &document.Timestamp)
		if err == nil {
			documents = append(documents, document)
		}
	}
	return documents, nil
}
