package selenoid

import (
	"path/filepath"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestGetConfigDir(t *testing.T) {
	t.Run("Get config dir", func(t *testing.T) {
		assert.Equal(t, DefaultConfigDir, GetSelenoidConfigDir())
		assert.Equal(t, DefaultConfigDir, GetSelenoidUIConfigDir())
		assert.True(t, filepath.IsAbs(GetSelenoidConfigDir()))
	})
}
