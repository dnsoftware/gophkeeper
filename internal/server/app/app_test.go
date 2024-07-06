package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNegative(t *testing.T) {
	err := ServerRun()
	require.Error(t, err)
}
