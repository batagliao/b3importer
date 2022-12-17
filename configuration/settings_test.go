package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {

	cfg := Get()
	assert.Equal(t, 2022, cfg.YearStart)
}
