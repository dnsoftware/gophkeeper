// Package entity_code работа со справочником типов сущностей
package entity_code

import "context"

// EntityCodeStorage интерфейс работы с хранилищем
type EntityCodeStorage interface {
	// GetEntityCodes получение кодов доступных сущностей
	GetEntityCodes(ctx context.Context) (map[string]string, error)
}

// EntityCode код-название сущности
type EntityCode struct {
	storage EntityCodeStorage
}

func NewEntityCode(storage EntityCodeStorage) (*EntityCode, error) {
	user := &EntityCode{
		storage: storage,
	}

	return user, nil
}

// EntityCodes получить список типов сущностей
func (k *EntityCode) EntityCodes(ctx context.Context) (map[string]string, error) {
	codes, err := k.storage.GetEntityCodes(ctx)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

//// EntityProps свойство сущности
//type EntityProps struct {
//	name             string // название
//	ftype            string //
//	validateRules    string
//	validateMessages string
//}
