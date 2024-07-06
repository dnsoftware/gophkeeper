package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	path, _ := os.Getwd()
	parts := strings.Split(path, "internal")
	file := parts[0] + "cmd/server/config.yaml"
	os.Args = append(os.Args, "-c="+file)
	cfg, err := NewServerConfig()
	require.NoError(t, err)
	assert.Equal(t, "local", cfg.Env)
	assert.Equal(t, "localhost:9090", cfg.ServerAddress)
}
