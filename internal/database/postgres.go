package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgres(conn string) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return nil, err
	}

	p := &PostgresDatabase{db}

	if _, err = p.getCurrentVersion(); err != nil {
		p.setVersion(0)
	}

	return p, nil
}

func (p *PostgresDatabase) DB() *sql.DB {
	return p.db
}

func (p *PostgresDatabase) getCurrentVersion() (int, error) {
	currentVersionRow := p.db.QueryRow("SELECT current_setting('my.version') AS verion")
	if currentVersionRow == nil {
		return 0, errors.New("missing user_version row")
	}

	currentVersionBuf := make([]byte, 1)
	if err := currentVersionRow.Scan(&currentVersionBuf); err != nil {
		return 0, err
	}

	currentVersion, err := strconv.Atoi(string(currentVersionBuf))
	if err != nil {
		return 0, err
	}

	return currentVersion, nil
}

func (p *PostgresDatabase) setVersion(version int) error {
	pragmaQuery := fmt.Sprintf("SET my.version TO %d", version)
	_, err := p.db.Exec(pragmaQuery)
	return err
}
