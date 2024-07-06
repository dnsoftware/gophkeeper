package domain

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/gophkeeper/internal/constants"
)

func TestCreateMetainfo(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	var metas []*Metainfo
	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input("Название поля метаданных:", "required", `{"required": "Укажите название поля метаданных"}`).Return("testmetaname", nil).AnyTimes()
	mockReadline.EXPECT().input("Значение поля метаданных:", "required", `{"required": "Укажите значение поля метаданных"}`).Return("testmetavalue", nil).AnyTimes()

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	metas = client.createMetainfo(metas)

	require.Equal(t, metas[0].Title, "testmetaname")

}

// Добавление сущности
func TestBaseAddEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil).AnyTimes()
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("1", nil)

	fieldGroup := []*Field{&Field{
		Id:               3,
		Name:             "Номер банковской карты",
		Etype:            "card",
		Ftype:            "string",
		ValidateRules:    "credit_card",
		ValidateMessages: `{"credit_card": "Неправильный формат номера карты"}`,
	}, &Field{
		Id:               4,
		Name:             "Месяц/Год (mm/yy) до которого действует карта",
		Etype:            "card",
		Ftype:            "string",
		ValidateRules:    "len=5",
		ValidateMessages: `{"len": "Месяц/год должны быть в формате mm/dd"}`,
	}, &Field{
		Id:               5,
		Name:             "Код проверки подлинности",
		Etype:            "card",
		Ftype:            "string",
		ValidateRules:    "len=3,number",
		ValidateMessages: `{"len": "Код должен состоять из трех цифр", "number": "Только число"}`,
	}}

	mockReadline.EXPECT().input("Номер банковской карты:", "credit_card", gomock.Any()).Return("12345678901234", nil)
	mockReadline.EXPECT().input("Месяц/Год (mm/yy) до которого действует карта:", "len=5", gomock.Any()).Return("12/24", nil)
	mockReadline.EXPECT().input("Код проверки подлинности:", "len=3,number", gomock.Any()).Return("12/24", nil)
	mockReadline.EXPECT().GetFieldsGroup("card").Return(fieldGroup)

	mockReadline.EXPECT().input("Метаданные или сохранение>>", "required,number", gomock.Any()).Return("1", nil).AnyTimes()

	mockReadline.EXPECT().input("Название поля метаданных:", "required", gomock.Any()).Return("testmetaname", nil).AnyTimes()
	mockReadline.EXPECT().input("Значение поля метаданных:", "required", gomock.Any()).Return("testmetavalue", nil).AnyTimes()

	mockReadline.EXPECT().input("Еще метаданные или сохранение>>", "required,number", gomock.Any()).Return("2", nil).AnyTimes()

	mockReadline.EXPECT().input("Сохранить или заново>>", "required,number", gomock.Any()).Return("1", nil).AnyTimes()

	sender.EXPECT().AddEntity(gomock.Any()).Return(int32(1), nil)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().GetEtypeName("card").Return("").AnyTimes()
	mockReadline.EXPECT().GetField(int32(3)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().GetField(int32(4)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().GetField(int32(5)).Return(&Field{}).AnyTimes()

	entCodes := []*EntityCode{&EntityCode{
		Etype: "card",
		Name:  "Банковская карта",
	}}

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}

// Просмотр сущности
func TestBaseViewEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil).AnyTimes()
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("2", nil)

	entList := make(map[int32]string)
	entList = map[int32]string{1: "111", 2: "222"}
	sender.EXPECT().EntityList(gomock.Any()).Return(entList, nil).AnyTimes()

	mockReadline.EXPECT().input("Просмотр объекта>>", "required,number", gomock.Any()).Return("1", nil)
	sender.EXPECT().Entity(gomock.Any()).Return(&Entity{
		Id:       1,
		UserID:   1,
		Etype:    "card",
		Props:    nil,
		Metainfo: nil,
	}, nil)

	mockReadline.EXPECT().input("Действия с объектом>>", "required,number", gomock.Any()).Return("0", nil)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().GetEtypeName("card").Return("").AnyTimes()
	mockReadline.EXPECT().GetField(int32(3)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().GetField(int32(4)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().GetField(int32(5)).Return(&Field{}).AnyTimes()

	entCodes := []*EntityCode{&EntityCode{
		Etype: "card",
		Name:  "Банковская карта",
	}}

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	entCodes = []*EntityCode{&EntityCode{
		Etype: "binary",
		Name:  "Бинарные данные",
	}}
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("2", nil)
	mockReadline.EXPECT().GetEtypeName("binary").Return("").AnyTimes()
	mockReadline.EXPECT().input("Просмотр объекта>>", "required,number", gomock.Any()).Return("1", nil)

	props := []*Property{&Property{
		EntityId: 1,
		FieldId:  1,
		Value:    `{"servername":"/servername","clientname":"/clientname","chunkcount":59}`,
	}}
	metas := []*Metainfo{&Metainfo{
		EntityId: 1,
		Title:    "meta",
		Value:    "meta",
	}}

	sender.EXPECT().Entity(gomock.Any()).Return(&Entity{
		Id:       1,
		UserID:   1,
		Etype:    "binary",
		Props:    props,
		Metainfo: metas,
	}, nil)
	sender.EXPECT().DownloadCryptoBinary(gomock.Any(), gomock.Any()).Return("path", nil)
	mockReadline.EXPECT().input("Действия с объектом>>", gomock.Any(), gomock.Any()).Return("1", nil)
	field := Field{
		Id:               1,
		Name:             "Произвольные бинарные данные (путь к файлу)",
		Etype:            constants.BinaryEntity,
		Ftype:            "path",
		ValidateRules:    "required,file",
		ValidateMessages: `{"required": "Путь к файлу не может быть пустым", "file": "Файла не существует"}`,
	}
	mockReadline.EXPECT().GetField(int32(1)).Return(&field)
	mockReadline.EXPECT().edit(field.Name+":", "", field.ValidateRules, field.ValidateMessages).Return("newval", nil)
	mockReadline.EXPECT().edit("Название метаданных:", gomock.Any(), gomock.Any(), gomock.Any()).Return("metanew", nil)
	mockReadline.EXPECT().edit("Значение метаданных:", gomock.Any(), gomock.Any(), gomock.Any()).Return("metanewval", nil)
	sender.EXPECT().SaveEntity(gomock.Any()).Return(int32(2), nil)
	sender.EXPECT().UploadCryptoBinary(int32(2), gomock.Any()).Return(int32(12345), nil)

	res, err = client.Base(entCodes)
	require.Equal(t, "again", res)

}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil).AnyTimes()
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("2", nil)

	entList := make(map[int32]string)
	entList = map[int32]string{1: "111", 2: "222"}
	sender.EXPECT().EntityList(gomock.Any()).Return(entList, nil).AnyTimes()

	mockReadline.EXPECT().input("Просмотр объекта>>", "required,number", gomock.Any()).Return("1", nil)
	sender.EXPECT().Entity(gomock.Any()).Return(&Entity{
		Id:       1,
		UserID:   1,
		Etype:    "card",
		Props:    nil,
		Metainfo: nil,
	}, nil)

	mockReadline.EXPECT().input("Действия с объектом>>", gomock.Any(), gomock.Any()).Return("2", nil)
	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().GetEtypeName("card").Return("").AnyTimes()
	mockReadline.EXPECT().GetField(int32(3)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().GetField(int32(4)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().GetField(int32(5)).Return(&Field{}).AnyTimes()

	entCodes := []*EntityCode{&EntityCode{
		Etype: "card",
		Name:  "Банковская карта",
	}}

	client, err = NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().input("Уверены (Y or N)>>", gomock.Any(), gomock.Any()).Return("y", nil)
	sender.EXPECT().DeleteEntity(gomock.Any()).Return(nil)

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}

func TestBaseNegative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", errors.New("testerr"))
	mockReadline.EXPECT().interrupt("1", errors.New("testerr")).Return(loopNone)
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("1", errors.New("testerr"))
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("1", nil)

	fieldGroup := []*Field{&Field{
		Id:               3,
		Name:             "Номер банковской карты",
		Etype:            "card",
		Ftype:            "string",
		ValidateRules:    "credit_card",
		ValidateMessages: `{"credit_card": "Неправильный формат номера карты"}`,
	}}
	mockReadline.EXPECT().GetFieldsGroup("card").Return(fieldGroup).AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("testerr"))
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("testerr"))
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("5", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("5", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("0", nil)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	entCodes := []*EntityCode{&EntityCode{
		Etype: "card",
		Name:  "Банковская карта",
	}}

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}

func TestBaseNegative2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("2", nil)

	entCodes := []*EntityCode{&EntityCode{
		Etype: constants.BinaryEntity,
		Name:  "Бинарные данные",
	}}

	fieldGroup := []*Field{&Field{
		Id:               7,
		Name:             "Произвольные бинарные данные (путь к файлу)",
		Etype:            "binary",
		Ftype:            "path",
		ValidateRules:    "required",
		ValidateMessages: `{"required": "Не может быть пустым"}`,
	}}
	mockReadline.EXPECT().GetFieldsGroup(constants.BinaryEntity).Return(fieldGroup).AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("2", nil)

	mockReadline.EXPECT().GetEtypeName("binary").Return("").AnyTimes()
	mockReadline.EXPECT().GetField(int32(7)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("2", errors.New("testerr"))
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("5", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)
	sender.EXPECT().AddEntity(gomock.Any()).Return(int32(1), nil)
	sender.EXPECT().UploadCryptoBinary(gomock.Any(), gomock.Any()).Return(int32(1000), errors.New("testerr"))

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}

func TestBaseActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("2", nil)

	entCodes := []*EntityCode{&EntityCode{
		Etype: constants.BinaryEntity,
		Name:  "Бинарные данные",
	}}

	fieldGroup := []*Field{&Field{
		Id:               7,
		Name:             "Произвольные бинарные данные (путь к файлу)",
		Etype:            "binary",
		Ftype:            "path",
		ValidateRules:    "required",
		ValidateMessages: `{"required": "Не может быть пустым"}`,
	}}
	mockReadline.EXPECT().GetFieldsGroup(constants.BinaryEntity).Return(fieldGroup).AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("2", nil)

	mockReadline.EXPECT().GetEtypeName("binary").Return("").AnyTimes()
	mockReadline.EXPECT().GetField(int32(7)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)
	sender.EXPECT().AddEntity(gomock.Any()).Return(int32(1), errors.New("testerr"))

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}

func TestBaseActions2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("2", nil)

	entCodes := []*EntityCode{&EntityCode{
		Etype: constants.BinaryEntity,
		Name:  "Бинарные данные",
	}}

	fieldGroup := []*Field{&Field{
		Id:               7,
		Name:             "Произвольные бинарные данные (путь к файлу)",
		Etype:            "binary",
		Ftype:            "path",
		ValidateRules:    "required",
		ValidateMessages: `{"required": "Не может быть пустым"}`,
	}}
	mockReadline.EXPECT().GetFieldsGroup(constants.BinaryEntity).Return(fieldGroup).AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("2", nil)

	mockReadline.EXPECT().GetEtypeName("binary").Return("").AnyTimes()
	mockReadline.EXPECT().GetField(int32(7)).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)
	sender.EXPECT().AddEntity(gomock.Any()).Return(int32(1), nil)
	sender.EXPECT().UploadCryptoBinary(gomock.Any(), gomock.Any()).Return(int32(1000), nil)

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}

func TestBaseActions3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("2", nil)

	entCodes := []*EntityCode{&EntityCode{
		Etype: constants.BinaryEntity,
		Name:  "Бинарные данные",
	}}

	sender.EXPECT().EntityList(gomock.Any()).Return(nil, errors.New("testerr"))

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}

func TestBaseActionsNegative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input("Выберите номер объекта:", "required,number", gomock.Any()).Return("1", nil).AnyTimes()
	mockReadline.EXPECT().interrupt("1", nil).Return(loopNone)
	mockReadline.EXPECT().input("Действия для объекта>>", "required,number", gomock.Any()).Return("2", nil)

	entList := make(map[int32]string)
	entList = map[int32]string{1: "111", 2: "222"}
	sender.EXPECT().EntityList(gomock.Any()).Return(entList, nil).AnyTimes()

	props := []*Property{&Property{
		EntityId: 1,
		FieldId:  1,
		Value:    "",
	}}

	sender.EXPECT().Entity(gomock.Any()).Return(&Entity{
		Id:       1,
		UserID:   1,
		Etype:    constants.BinaryEntity,
		Props:    props,
		Metainfo: nil,
	}, nil)

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	entCodes := []*EntityCode{&EntityCode{
		Etype: constants.BinaryEntity,
		Name:  "Бинарные данные",
	}}

	mockReadline.EXPECT().GetEtypeName(constants.BinaryEntity).Return("").AnyTimes()
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", errors.New("testerr"))
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("bad", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)

	props[0].Value = `{"servername":"/servername","clientname":"/clientname","chunkcount":59}`
	sender.EXPECT().Entity(gomock.Any()).Return(&Entity{
		Id:       1,
		UserID:   1,
		Etype:    constants.BinaryEntity,
		Props:    props,
		Metainfo: nil,
	}, nil)

	sender.EXPECT().DownloadCryptoBinary(gomock.Any(), gomock.Any()).Return("", errors.New("testerr"))
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)
	sender.EXPECT().DownloadCryptoBinary(gomock.Any(), gomock.Any()).Return("", nil)
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", errors.New("testerr"))
	mockReadline.EXPECT().input(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)
	mockReadline.EXPECT().GetField(gomock.Any()).Return(&Field{}).AnyTimes()
	mockReadline.EXPECT().edit(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("test", nil)
	sender.EXPECT().SaveEntity(gomock.Any()).Return(int32(2), errors.New("testerr"))

	res, err := client.Base(entCodes)
	require.Equal(t, "again", res)

}
