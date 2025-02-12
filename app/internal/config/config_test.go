package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchConfigPath(t *testing.T) {
	tests := []struct {
		name          string
		envConfigPath string
		expectedPath  string
	}{
		{
			name:          "config path from env",
			envConfigPath: "/path/to/config.yaml",
			expectedPath:  "/path/to/config.yaml",
		},
		{
			name:          "default config path",
			envConfigPath: "",
			expectedPath:  "config/config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменную окружения
			if tt.envConfigPath != "" {
				os.Setenv("CONFIG_FILE", tt.envConfigPath)
				defer os.Unsetenv("CONFIG_FILE")
			}

			// Вызываем функцию и проверяем результат
			path := fetchConfigPath()
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}

func TestMustLoadPath(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		expectError bool
	}{
		{
			name:        "valid config file",
			configPath:  "config_test.yaml",
			expectError: false,
		},
		{
			name:        "non-existent config file",
			configPath:  "non_existent_config.yaml",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectError {
				// Проверяем, что функция вызывает panic
				assert.Panics(t, func() {
					MustLoadPath(tt.configPath)
				})
			} else {
				// Проверяем, что функция не вызывает panic
				assert.NotPanics(t, func() {
					cfg := MustLoadPath(tt.configPath)
					assert.NotNil(t, cfg)
					assert.Equal(t, "local", cfg.Env)
					assert.Equal(t, "postgres", cfg.StorageType)
					assert.Equal(t, 8080, cfg.GRPC.Port)
					assert.Equal(t, 5*time.Second, cfg.GRPC.Timeout)
					assert.Equal(t, "postgres_db", cfg.Postgres.Host)
					assert.Equal(t, 5432, cfg.Postgres.Port)
					assert.Equal(t, "admin", cfg.Postgres.User)
					assert.Equal(t, "admin", cfg.Postgres.Password)
					assert.Equal(t, "url_shortener_db", cfg.Postgres.DBName)
					assert.Equal(t, "redis_d", cfg.Redis.Host)
					assert.Equal(t, "6379", cfg.Redis.Port)
					assert.Equal(t, "", cfg.Redis.Password)
					assert.Equal(t, 0, cfg.Redis.DB)
				})
			}
		})
	}
}

func TestMustLoad(t *testing.T) {
	tests := []struct {
		name          string
		envConfigPath string
		expectError   bool
	}{
		{
			name:          "valid config path",
			envConfigPath: "config_test.yaml",
			expectError:   false,
		},
		{
			name:          "empty config path",
			envConfigPath: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envConfigPath != "" {
				os.Setenv("CONFIG_FILE", tt.envConfigPath)
				defer os.Unsetenv("CONFIG_FILE")
			}

			if tt.expectError {
				assert.Panics(t, func() {
					MustLoad()
				})
			} else {
				assert.NotPanics(t, func() {
					cfg := MustLoad()
					assert.NotNil(t, cfg)
				})
			}
		})
	}
}

func TestStorageTypeFromEnv(t *testing.T) {
	os.Setenv("STORAGE_TYPE", "redis")
	defer os.Unsetenv("STORAGE_TYPE")

	cfg := MustLoadPath("config_test.yaml")

	// Проверяем, что значение переопределено
	assert.Equal(t, "redis", cfg.StorageType)
}
