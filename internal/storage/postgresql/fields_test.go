package postgresql

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

func TestFields(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	ctx := context.Background()
	fields, err := db.GetEntityFields(ctx, "card")
	require.NoError(t, err)
	require.Equal(t, len(fields), 3)
}

func TestGetFieldByEtypeAndName(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	ctx := context.Background()
	id, ftype, err := db.GetFieldByEtypeAndName(ctx, "card", "Номер банковской карты")
	require.NoError(t, err)
	require.Greater(t, id, int32(0))
	require.Equal(t, ftype, "string")
}

func TestIsFieldType(t *testing.T) {
	db, err := setupDatabase()
	require.NoError(t, err)

	ctx := context.Background()
	id, ftype, err := db.GetFieldByEtypeAndName(ctx, "card", "Номер банковской карты")

	isType, err := db.IsFieldType(ctx, id, ftype)
	require.NoError(t, err)
	require.True(t, isType)
}
