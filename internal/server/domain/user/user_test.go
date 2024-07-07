package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dnsoftware/gophkeeper/internal/server/mocks"
)

func TestUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockUserStorage(ctrl)
	user, err := NewUser(mockStorage)
	ctx := context.Background()

	token, err := user.Registration(ctx, "login", "pass", "repeat")
	assert.Error(t, err)
	assert.Equal(t, "", token)

	mockStorage.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(0, time.Now(), errors.New("testerr"))
	token, err = user.Registration(ctx, "login", "pass", "pass")
	assert.Error(t, err)
	assert.Equal(t, "", token)

	mockStorage.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(1, time.Now(), nil)
	token, err = user.Registration(ctx, "login", "pass", "pass")
	assert.Error(t, err)
	assert.Equal(t, "", token)

}
