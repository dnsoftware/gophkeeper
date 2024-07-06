package domain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("r", nil).AnyTimes()
	mockReadline.EXPECT().Registration().Return("login", "password", nil).AnyTimes()
	sender.EXPECT().Registration("login", "password", "password").Return("token", nil)

	var entCodes []*EntityCode

	sender.EXPECT().EntityCodes().Return(entCodes, nil)
	mockReadline.EXPECT().Close().Return(nil).AnyTimes()

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	err = client.Start()
	require.NoError(t, err)

	mockReadline.EXPECT().SetEtypeName("card", "Банковская карта").Return()
	field := Field{}
	fields := []*Field{&field}
	entCodes = []*EntityCode{&EntityCode{
		Etype: "card",
		Name:  "Банковская карта",
	}}
	sender.EXPECT().Fields("card").Return(fields, nil)
	mockReadline.EXPECT().MakeFieldsDescription(fields).Return()
	sender.EXPECT().EntityCodes().Return(entCodes, nil)

	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("", nil).AnyTimes()
	sender.EXPECT().Registration("login", "password", "password").Return("", nil)
	mockReadline.EXPECT().Login().Return("login", "password", nil).AnyTimes()
	sender.EXPECT().Login("login", "password").Return("", nil).AnyTimes()
	mockReadline.EXPECT().input(`Выберите номер объекта:`, gomock.Any(), gomock.Any()).Return("", fmt.Errorf("testerror"))
	mockReadline.EXPECT().interrupt("", gomock.Any()).Return(loopBreak)

	err = client.Start()
	require.NoError(t, err)

}

func TestStartNegative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sender := NewMockSender(ctrl)

	controller := gomock.NewController(t)
	mockReadline := NewMockReadline(controller)
	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("q", errors.New("testerr"))

	client, err := NewGophKeepClient(mockReadline, sender)
	require.NoError(t, err)

	err = client.Start()
	require.Error(t, err)

	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("q", nil)
	mockReadline.EXPECT().Login().Return("login", "password", errors.New("testerr"))

	err = client.Start()
	require.Error(t, err)

	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("r", nil)
	mockReadline.EXPECT().Registration().Return("login", "password", errors.New("testerr"))

	err = client.Start()
	require.Error(t, err)

	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("r", nil)
	mockReadline.EXPECT().Registration().Return("login", "password", nil)
	sender.EXPECT().Registration("login", "password", "password").Return("", errors.New("testerr"))
	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("q", errors.New("testerr"))

	err = client.Start()
	require.Error(t, err)

	mockReadline.EXPECT().input(`Нажмите [Enter] для входа или "r" для регистрации>>`, "", gomock.Any()).Return("q", nil)
	mockReadline.EXPECT().Login().Return("login", "password", nil)
	sender.EXPECT().Login("login", "password").Return("", errors.New("testerr"))
	mockReadline.EXPECT().Login().Return("login", "password", nil)
	sender.EXPECT().Login("login", "password").Return("token", nil)
	sender.EXPECT().EntityCodes().Return(nil, errors.New("testerr"))
	sender.EXPECT().Fields("card").Return(nil, errors.New("testerr")).AnyTimes()

	err = client.Start()
	require.NoError(t, err)

}
