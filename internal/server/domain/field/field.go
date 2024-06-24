package field

import "context"

type FieldStorage interface {
	GetEntityFields(ctx context.Context, etype string) ([]EntityFields, error)
	IsFieldType(ctx context.Context, id int32, ftype string) (bool, error)
}

type EntityFields struct {
	ID               int32
	Name             string
	Ftype            string
	ValidateRules    string
	ValidateMessages string
}

type Field struct {
	storage FieldStorage
}

func NewField(storage FieldStorage) (*Field, error) {
	f := &Field{
		storage: storage,
	}

	return f, nil
}

func (f *Field) Fields(ctx context.Context, etype string) ([]EntityFields, error) {

	fields, err := f.storage.GetEntityFields(ctx, etype)
	if err != nil {
		return nil, err
	}

	return fields, nil
}

func (f *Field) isFieldType(ctx context.Context, id int32, ftype string) bool {

	isType, _ := f.storage.IsFieldType(ctx, id, ftype)

	return isType
}
