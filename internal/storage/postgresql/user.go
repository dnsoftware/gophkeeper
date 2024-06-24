package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dnsoftware/gophkeeper/internal/constants"
	"github.com/dnsoftware/gophkeeper/internal/utils"
	"github.com/dnsoftware/gophkeeper/logger"
)

// GetUser получение данных пользователя (возвращает ID и дату добавления, если ID = 0 - такого пользователя нет)
func (p *PgStorage) GetUser(ctx context.Context, login string) (int, time.Time, error) {

	query := `SELECT id, created_at FROM users WHERE login = $1`
	row := p.db.QueryRowContext(ctx, query, login)

	var (
		id        int
		createdAt time.Time
	)

	err := row.Scan(&id, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, time.Time{}, nil
		} else {
			return 0, time.Time{}, fmt.Errorf("GetUser: %w", err)
		}
	}

	return id, createdAt, nil
}

// UserCreate регистрация нового пользователя.
func (p *PgStorage) UserCreate(ctx context.Context, login string, password string, salt string) (int, error) {

	query := "INSERT INTO users (login, password, salt, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	_, err := p.db.ExecContext(ctx, query, login, password, salt, time.Now())
	if err != nil {
		return 0, err
	}
	q := "SELECT LASTVAL() id"
	r := p.db.QueryRowContext(ctx, q)
	var idNew int
	err = r.Scan(&idNew)
	if err != nil {
		return 0, err
	}

	return idNew, nil
}

// LoginUser проверка наличия пары логин-пароль, пароль подается в исходном виде
// возвращает ID пользователя и пустую строку в случае успеха или 0 с текстом описания, если пользователя нет в базе
func (p *PgStorage) LoginUser(ctx context.Context, login string, password string) (int, string) {

	var (
		salt string
		id   int
	)

	query := `SELECT salt FROM users WHERE login = $1`
	row := p.db.QueryRowContext(ctx, query, login)
	err := row.Scan(&salt)
	if err != nil {
		logger.Log().Error("LoginCheckUser, get salt error: " + err.Error())
		return 0, constants.ErrNoSuchUser
	}

	passHash := utils.PassGenerate(password, salt)

	query = `SELECT id FROM users WHERE login = $1 AND password = $2`
	row = p.db.QueryRowContext(ctx, query, login, passHash)

	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, constants.ErrBadPassword
		} else {
			return 0, err.Error()
		}
	}

	return id, ""
}
