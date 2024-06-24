package postgresql

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupDatabase настройка чистой тестовой БД с начальными данными
func setupDatabase() (*PgStorage, error) {
	ctx := context.Background()

	dbname := "users"
	user := "user"
	password := "password"

	// Start the postgres container
	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	// очистка БД
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	query := "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// конец очистки БД

	pgs, err := NewPostgresqlStorage(dsn)
	if err != nil {
		return nil, err
	}

	// заполняем базу через миграции
	// не забываем импортировать
	// _ "github.com/golang-migrate/migrate/v4/source/file"
	// _ "github.com/golang-migrate/migrate/v4/database/postgres"

	path, _ := os.Getwd()
	parts := strings.Split(path, "internal")
	sourceURL := "file://" + parts[0] + "cmd/server/migrations"
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return nil, err
	} else {
		err = m.Up()
		if err != nil && err.Error() != "no change" {
			return nil, err
		}

	}

	return pgs, nil
}
