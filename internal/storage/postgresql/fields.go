package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dnsoftware/gophkeeper/internal/server/domain/field"
)

// GetEntityFields получение набора полей сущности
func (p *PgStorage) GetEntityFields(ctx context.Context, etype string) ([]field.EntityFields, error) {

	query := `SELECT id, name, ftype, validate_rules, validate_messages FROM fields WHERE etype = $1`
	rows, err := p.db.QueryContext(ctx, query, etype)
	if err != nil {
		return nil, fmt.Errorf("GetEntityCodes error: %w", err)
	}
	defer rows.Close()

	var ef field.EntityFields
	var fields []field.EntityFields
	for rows.Next() {
		err := rows.Scan(&ef.ID, &ef.Name, &ef.Ftype, &ef.ValidateRules, &ef.ValidateMessages)
		if err != nil {
			return nil, fmt.Errorf("GetEntityFields: %w", err)
		}
		ef.Etype = etype
		fields = append(fields, ef)
	}

	return fields, nil

}

// GetFieldByEtypeAndName получить код записи и тип по типу сущности и имени
func (p *PgStorage) GetFieldByEtypeAndName(ctx context.Context, etype string, name string) (int32, string, error) {
	query := `SELECT id, ftype  FROM fields WHERE etype = $1 AND name = $2`

	var id int32
	var ftype string
	row := p.db.QueryRowContext(ctx, query, etype, name)
	err := row.Scan(&id, &ftype)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", nil
		} else {
			return 0, "", err
		}
	}

	return id, ftype, nil
}

// IsFieldType имеет ли поле с указанным идентификатором определенный тип?
func (p *PgStorage) IsFieldType(ctx context.Context, id int32, ftype string) (bool, error) {
	query := `SELECT name FROM fields WHERE id = $1 AND ftype = $2`

	var name string
	row := p.db.QueryRowContext(ctx, query, id, ftype)
	err := row.Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
