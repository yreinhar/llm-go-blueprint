package run

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultConfig(t *testing.T) {
	config := newDefaultConfig()

	assert.Equal(t, "8080", config.Port, "default port should be 8080")
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	mockEnv := func(key string) string { return "" }

	config, err := loadConfig("", mockEnv)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Port) // Should return default config.
}

func TestLoadConfig_InvalidPath(t *testing.T) {
	mockEnv := func(key string) string { return "" }

	config, err := loadConfig("/path/that/doesnot/exist/config.yaml", mockEnv)
	assert.NoError(t, err) // Should not error, just use defaults.
	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Port)
}

func TestLoadConfig_FilePermissions(t *testing.T) {
	mockEnv := func(key string) string { return "" }

	// Create temp config file.
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// No permissions.
	filePermission := os.FileMode(0000)

	// Create file with given permissions.
	err := os.WriteFile(configPath, []byte("port: \"3000\""), filePermission)
	require.NoError(t, err)

	_, err = loadConfig(configPath, mockEnv)
	assert.Error(t, err)
}

func TestLoadConfig_EnvOverrideDefaultConfig(t *testing.T) {
	mockEnv := func(key string) string { return "9090" }

	config, err := loadConfig("", mockEnv)
	assert.NoError(t, err)
	assert.Equal(t, "9090", config.Port)
}

func TestLoadConfig_EnvOverrideYAMLConfig(t *testing.T) {
	mockEnv := func(key string) string { return "9090" }

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Read and write permissions.
	filePermission := os.FileMode(0666)

	err := os.WriteFile(configPath, []byte("port: \"3000\""), filePermission)
	require.NoError(t, err)

	config, err := loadConfig(configPath, mockEnv)
	assert.NoError(t, err)
	assert.Equal(t, "9090", config.Port)
}

func TestLoadConfig_YAMLOverrideDefaultConfig(t *testing.T) {
	mockEnv := func(key string) string { return "" }

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	filePermission := os.FileMode(0666)

	err := os.WriteFile(configPath, []byte("port: \"9090\""), filePermission)
	require.NoError(t, err)

	config, err := loadConfig(configPath, mockEnv)
	assert.NoError(t, err)
	assert.Equal(t, "9090", config.Port)
}
