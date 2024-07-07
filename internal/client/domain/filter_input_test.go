package domain

import (
	"testing"

	"github.com/chzyer/readline"
	"github.com/stretchr/testify/assert"

	"github.com/dnsoftware/gophkeeper/internal/constants"
)

func TestFilter(t *testing.T) {
	stopChan := make(chan bool, 1)
	f := NewFilter("/tmp/test", stopChan)
	r, ok := f.FilterInput(readline.CharCtrlZ)
	assert.Equal(t, false, ok)
	assert.Equal(t, int32(readline.CharCtrlZ), r)

	r, ok = f.FilterInput(constants.CharCtrlC)
	assert.Equal(t, true, ok)
	assert.Equal(t, int32(constants.CharCtrlC), r)

}
