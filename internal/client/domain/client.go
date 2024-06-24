package domain

import "context"

// Sender Интерфейс отправки/приема данных с сервера
type Sender interface {
	Registration(ctx context.Context, login string, password string, password2 string) (string, error)
	Login(ctx context.Context, login string, password string) (string, error)
	EntityCodes(ctx context.Context) ([]*EntityCode, error)
	Fields(ctx context.Context, etype string) ([]*Field, error)
	AddEntity(ctx context.Context, ae Entity) (int32, error)
	UploadBinary(ctx context.Context, entityId int32, file string) (int32, error)
	DownloadBinary(ctx context.Context, entityId int32, fileName string) (string, error)
	UploadCryptoBinary(ctx context.Context, entityId int32, file string) (int32, error)
	DownloadCryptoBinary(ctx context.Context, entityId int32, fileName string) (string, error)
}

type Entity struct {
	Id       int32       // ID сущности
	UserID   int32       // ID пользователя
	Etype    string      // тип сущности: card, text, logopas, binary и т.д.
	Props    []*Property // массив значений свойств
	Metainfo []*Metainfo // массив значений метаинформации
}

type Property struct {
	EntityId int32  // код сущности
	FieldId  int32  // код описания поля свйоства
	Value    string // значение свойства
}

type Metainfo struct {
	EntityId int32  // код сущности
	Title    string // наименование метаинформации
	Value    string // значение метаинформации
}

type EntityCode struct {
	Etype string
	Name  string
}

type Field struct {
	Id               int32
	Name             string
	Ftype            string
	ValidateRules    string
	ValidateMessages string
}

type GophKeepClient struct {
	Sender *Sender
}

func NewGophKeepClient(sender *Sender) (*GophKeepClient, error) {

	client := &GophKeepClient{
		Sender: sender,
	}

	return client, nil

}
