package postgresql

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/dnsoftware/gophkeeper/logger"
)

// PgStorage работает с Postgresql базой данных.
type PgStorage struct {
	db *sql.DB
}

func NewPostgresqlStorage(dsn string) (*PgStorage, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log().Error(err.Error())
		return nil, err
	}

	ps := &PgStorage{
		db: db,
	}

	return ps, nil
}

//func (p *PgStorage) GetEntityProperties(code string) []EntityProps {
//
//}
