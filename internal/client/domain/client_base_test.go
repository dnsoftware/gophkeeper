package domain

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
	sender.EXPECT().EntityList(gomock.Any()).Return(entList, nil)

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

}
