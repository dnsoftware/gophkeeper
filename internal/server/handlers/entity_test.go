package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/gophkeeper/internal/server/domain/entity"
	mock_domain "github.com/dnsoftware/gophkeeper/internal/server/mocks"
)

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoFields := mock_domain.NewMockFieldRepo(ctrl)
	repoEntity := mock_domain.NewMockEntityRepo(ctrl)
	entityService, _ := entity.NewEntity(repoEntity, repoFields)

	ctx := context.Background()
	ent := entity.EntityModel{
		ID:       1,
		UserID:   1,
		Etype:    "card",
		Props:    nil,
		Metainfo: nil,
	}

	t.Run("del2", func(t *testing.T) {
		repoEntity.EXPECT().GetEntity(ctx, int32(1)).Return(ent, errors.New("testerr"))

		err := entityService.DeleteEntity(ctx, 1, 1)
		require.Error(t, err)
	})

	t.Run("del2", func(t *testing.T) {
		repoEntity.EXPECT().GetEntity(ctx, int32(1)).Return(ent, nil)
		repoEntity.EXPECT().DeleteEntity(ctx, int32(1), int32(1)).Return(errors.New("testerr"))

		err := entityService.DeleteEntity(ctx, 1, 1)
		assert.Error(t, err)
	})

}
