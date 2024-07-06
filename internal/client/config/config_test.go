package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	path, _ := os.Getwd()
	parts := strings.Split(path, "internal")

	t.Run("test1", func(t *testing.T) {
		file := parts[0] + "cmd/client/config.yaml"
		os.Args = append(os.Args, "-c="+file)
		cfg, err := NewClientConfig()
		require.NoError(t, err)
		assert.Equal(t, "local", cfg.Env)
		assert.Equal(t, "localhost:9090", cfg.ServerAddress)
		assert.Equal(t, "Secret", cfg.SecretKey)
	})

}
