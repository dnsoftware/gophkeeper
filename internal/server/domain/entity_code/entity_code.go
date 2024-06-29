package entity_code

import "context"

type EntityCodeStorage interface {
	// GetEntityCodes получение кодо доступных сущностей
	GetEntityCodes(ctx context.Context) (map[string]string, error)
}

type EntityCode struct {
	storage EntityCodeStorage
}

func NewEntityCode(storage EntityCodeStorage) (*EntityCode, error) {
	user := &EntityCode{
		storage: storage,
	}

	return user, nil
}

func (k *EntityCode) EntityCodes(ctx context.Context) (map[string]string, error) {
	codes, err := k.storage.GetEntityCodes(ctx)
	if err != nil {

	}

	return codes, nil
}

type EntityProps struct {
	name             string
	ftype            string
	validateRules    string
	validateMessages string
}
