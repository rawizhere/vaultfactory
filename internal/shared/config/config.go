package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type SecurityConfig struct {
	JWTSecret                  string        `mapstructure:"jwt_secret"`
	EncryptionKey              string        `mapstructure:"encryption_key"`
	JWTExpireHours             int           `mapstructure:"jwt_expire_hours"`
	RefreshTokenExpireDays     int           `mapstructure:"refresh_token_expire_days"`
	JWTExpireDuration          time.Duration `mapstructure:"-"`
	RefreshTokenExpireDuration time.Duration `mapstructure:"-"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

func LoadConfig(configPath string) (*Config, error) {
	// Загружаем .env файл если он существует
	godotenv.Load()

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Включаем чтение из переменных окружения
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.dsn", "postgres://vaultfactory:password@localhost:5432/vaultfactory?sslmode=disable")
	viper.SetDefault("security.jwt_expire_hours", 24)
	viper.SetDefault("security.refresh_token_expire_days", 30)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	config.Security.JWTExpireDuration = time.Duration(config.Security.JWTExpireHours) * time.Hour
	config.Security.RefreshTokenExpireDuration = time.Duration(config.Security.RefreshTokenExpireDays) * 24 * time.Hour

	return &config, nil
}

func (c *Config) GetDSN() string {
	// Проверяем установлена ли переменная окружения DB_DSN
	if dsn := viper.GetString("DB_DSN"); dsn != "" {
		return dsn
	}

	return c.Database.DSN
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) GetJWTSecret() string {
	// Проверяем установлена ли переменная окружения JWT_SECRET
	if jwtSecret := viper.GetString("JWT_SECRET"); jwtSecret != "" {
		return jwtSecret
	}
	return c.Security.JWTSecret
}

func (c *Config) GetEncryptionKey() string {
	// Проверяем установлена ли переменная окружения ENCRYPTION_KEY
	if encryptionKey := viper.GetString("ENCRYPTION_KEY"); encryptionKey != "" {
		return encryptionKey
	}
	return c.Security.EncryptionKey
}
