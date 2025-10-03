package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// ConfigReader определяет интерфейс для чтения конфигурации.
type ConfigReader interface {
	GetDSN() string
	GetServerAddr() string
	GetJWTSecret() string
	GetEncryptionKey() string
	GetJWTExpireDuration() time.Duration
	GetRefreshTokenExpireDuration() time.Duration
	GetLoggingLevel() string
	GetLoggingFormat() string
	GetLoggingOutput() string
}

type config struct {
	server   serverConfig
	database databaseConfig
	security securityConfig
	logging  loggingConfig
}

type serverConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type databaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type securityConfig struct {
	JWTSecret                  string        `mapstructure:"jwt_secret"`
	EncryptionKey              string        `mapstructure:"encryption_key"`
	JWTExpireHours             int           `mapstructure:"jwt_expire_hours"`
	RefreshTokenExpireDays     int           `mapstructure:"refresh_token_expire_days"`
	JWTExpireDuration          time.Duration `mapstructure:"-"`
	RefreshTokenExpireDuration time.Duration `mapstructure:"-"`
}

type loggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// LoadConfig загружает конфигурацию из файла и возвращает ConfigReader.
func LoadConfig(configPath string) (ConfigReader, error) {
	_ = godotenv.Load()

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

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

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.security.JWTExpireDuration = time.Duration(cfg.security.JWTExpireHours) * time.Hour
	cfg.security.RefreshTokenExpireDuration = time.Duration(cfg.security.RefreshTokenExpireDays) * 24 * time.Hour

	return &cfg, nil
}

func (c *config) GetDSN() string {
	if dsn := viper.GetString("DB_DSN"); dsn != "" {
		return dsn
	}
	return c.database.DSN
}

func (c *config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.server.Host, c.server.Port)
}

func (c *config) GetJWTSecret() string {
	if jwtSecret := viper.GetString("JWT_SECRET"); jwtSecret != "" {
		return jwtSecret
	}
	return c.security.JWTSecret
}

func (c *config) GetEncryptionKey() string {
	if encryptionKey := viper.GetString("ENCRYPTION_KEY"); encryptionKey != "" {
		return encryptionKey
	}
	return c.security.EncryptionKey
}

func (c *config) GetJWTExpireDuration() time.Duration {
	return c.security.JWTExpireDuration
}

func (c *config) GetRefreshTokenExpireDuration() time.Duration {
	return c.security.RefreshTokenExpireDuration
}

func (c *config) GetLoggingLevel() string {
	return c.logging.Level
}

func (c *config) GetLoggingFormat() string {
	return c.logging.Format
}

func (c *config) GetLoggingOutput() string {
	return c.logging.Output
}
