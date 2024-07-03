package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
)

// CreateEntity создание сущности
func (p *PgStorage) CreateEntity(ctx context.Context, entity entity.EntityModel) (int32, error) {

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	query := "INSERT INTO entities (user_id, etype, created_at) VALUES ($1, $2, $3) RETURNING id"
	_, err = tx.ExecContext(ctx, query, entity.UserID, entity.Etype, time.Now())
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	q := "SELECT LASTVAL() id"
	r := tx.QueryRowContext(ctx, q)
	var idEntity int32
	err = r.Scan(&idEntity)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// заносим свойства
	for _, prop := range entity.Props {
		query := "INSERT INTO properties (entity_id, field_id, value) VALUES ($1, $2, $3)"
		_, err = tx.ExecContext(ctx, query, idEntity, prop.FieldID, prop.Value)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// заносим метаинформацию
	for _, meta := range entity.Metainfo {
		query := "INSERT INTO metainfo (entity_id, title, value) VALUES ($1, $2, $3)"
		_, err = tx.ExecContext(ctx, query, idEntity, meta.Title, meta.Value)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	return idEntity, nil

}

// UpdateEntity Сохранение отредактированной сущности
func (p *PgStorage) UpdateEntity(ctx context.Context, entity entity.EntityModel) error {

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	// заносим свойства
	for _, prop := range entity.Props {
		query := "UPDATE properties SET value = $1 WHERE entity_id = $2 AND field_id = $3"
		_, err = tx.ExecContext(ctx, query, prop.Value, entity.ID, prop.FieldID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// заносим метаинформацию
	// Удаляем старую метаинформацию
	queryDel := "DELETE FROM metainfo WHERE entity_id = $1"
	_, err = tx.ExecContext(ctx, queryDel, entity.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Добавляем новые
	for _, meta := range entity.Metainfo {
		query := "INSERT INTO metainfo (entity_id, title, value) VALUES ($1, $2, $3)"
		_, err = tx.ExecContext(ctx, query, entity.ID, meta.Title, meta.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil

}

func (p *PgStorage) GetEntity(ctx context.Context, id int32) (entity.EntityModel, error) {

	empty := entity.EntityModel{}

	query := "SELECT user_id, etype FROM entities WHERE id = $1"
	var userID int32
	var etype string
	row := p.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&userID, &etype)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return empty, fmt.Errorf("no entity with id: %v", id)
		} else {
			return empty, err
		}
	}

	ent := entity.EntityModel{
		ID:     id,
		UserID: userID,
		Etype:  etype,
	}

	// получаем свойства
	var props []entity.Property
	var prop entity.Property

	query = `SELECT id, field_id, value FROM properties WHERE entity_id = $1`
	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		return empty, fmt.Errorf("GetEntity error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&prop.ID, &prop.FieldID, &prop.Value)
		if err != nil {
			return empty, fmt.Errorf("scan Property error: %w", err)
		}
		prop.EntityID = id
		props = append(props, prop)
	}

	// получаем метаинформацию
	var metainfo []entity.Metainfo
	var meta entity.Metainfo

	query = `SELECT id, title, value FROM metainfo WHERE entity_id = $1`
	rows, err = p.db.QueryContext(ctx, query, id)
	if err != nil {
		return empty, fmt.Errorf("select metainfo error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&meta.ID, &meta.Title, &meta.Value)
		if err != nil {
			return empty, fmt.Errorf("scan Property error: %w", err)
		}
		meta.EntityID = id
		metainfo = append(metainfo, meta)
	}

	ent.Props = props
	ent.Metainfo = metainfo

	return ent, nil
}

// GetBinaryFilenameByEntityID Получение данных по именам файлов хранения бинарных данных
func (p *PgStorage) GetBinaryFilenameByEntityID(ctx context.Context, entityID int32) (string, error) {
	query := "SELECT p.value FROM entities e, properties p WHERE e.id = $1 AND e.id = p.entity_id LIMIT 1"
	var filename string
	row := p.db.QueryRowContext(ctx, query, entityID)
	err := row.Scan(&filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no property with entityID: %v", entityID)
		} else {
			return "", err
		}
	}

	return filename, nil
}

type BinaryFileDataProperty struct {
	Servername string `json:"servername"`
	Clientname string `json:"clientname"`
	Chunkcount int32  `json:"chunkcount"`
}

// SetChunkCountForCryptoBinary Сохранение кол-ва фрагментов, на которые разбит бинарный файл
func (p *PgStorage) SetChunkCountForCryptoBinary(ctx context.Context, entityID int32, chunkCount int32) error {
	query := "SELECT p.id property_id, p.value FROM entities e, properties p WHERE e.id = $1 AND e.id = p.entity_id LIMIT 1"
	var filedata string
	var propertyID int32
	row := p.db.QueryRowContext(ctx, query, entityID)
	err := row.Scan(&propertyID, &filedata)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no property with entityID: %v", entityID)
		} else {
			return err
		}
	}

	fd := &BinaryFileDataProperty{}
	err = json.Unmarshal([]byte(filedata), fd)
	if err != nil {
		return err
	}

	fd.Chunkcount = chunkCount
	filedataStr, err := json.Marshal(fd)
	if err != nil {
		return err
	}

	query = "UPDATE properties SET value = $1 WHERE id = $2"
	_, err = p.db.ExecContext(ctx, query, filedataStr, propertyID)
	if err != nil {
		return err
	}

	return nil
}

// GetEntityListByType Получение списка сущностей указанного типа для конкретного пользователя
// Простая карта с кодом сущности и названием(составляется из метаданных)
func (p *PgStorage) GetEntityListByType(ctx context.Context, etype string, userID int32) (map[int32][]string, error) {

	query := `SELECT e.id, m.title, m.value FROM entities e LEFT JOIN metainfo m 
                        ON e.id = m.entity_id
                        WHERE e.etype = $1 AND e.user_id = $2`
	rows, err := p.db.QueryContext(ctx, query, etype, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var id int32
	var title, value sql.NullString
	var list = make(map[int32][]string)
	for rows.Next() {
		err := rows.Scan(&id, &title, &value)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		list[id] = append(list[id], title.String+":"+value.String)
	}

	return list, nil
}

// DeleteEntity Удаление данных сущности из базы
func (p *PgStorage) DeleteEntity(ctx context.Context, id int32, userID int32) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM entities WHERE id = $1 AND user_id = $2 "
	res, err := p.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected > 0 {
		queryProps := "DELETE FROM properties WHERE entity_id = $1"
		_, err := p.db.ExecContext(ctx, queryProps, id)
		if err != nil {
			tx.Rollback()
			return err
		}

		queryMetas := "DELETE FROM metainfo WHERE entity_id = $1"
		_, err = p.db.ExecContext(ctx, queryMetas, id)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	tx.Commit()

	return nil
}
