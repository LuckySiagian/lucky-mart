package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type AppConfig struct {
	Name    string `mapstructure:"NAME"`
	Env     string `mapstructure:"ENV"`
	Version string `mapstructure:"VERSION"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"PORT"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"DB_HOST"`
	Port            string        `mapstructure:"DB_PORT"`
	User            string        `mapstructure:"DB_USER"`
	Password        string        `mapstructure:"DB_PASSWORD"`
	Name            string        `mapstructure:"DB_NAME"`
	SSLMode         string        `mapstructure:"DB_SSLMODE"`
	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"JWT_ACCESS_SECRET"`
	RefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"`
	AccessTTL     time.Duration `mapstructure:"JWT_ACCESS_TTL"`
	RefreshTTL    time.Duration `mapstructure:"JWT_REFRESH_TTL"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"ALLOW_ORIGINS"`
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(fmt.Sprintf(".env.%s", env))
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set Default Configurations
	viper.SetDefault("APP.NAME", "LuckyMart API")
	viper.SetDefault("APP.ENV", env)
	viper.SetDefault("APP.VERSION", "1.0.0")

	viper.SetDefault("SERVER.PORT", "8080")
	viper.SetDefault("SERVER.READ_TIMEOUT", 15*time.Second)
	viper.SetDefault("SERVER.WRITE_TIMEOUT", 15*time.Second)

	viper.SetDefault("DATABASE.DB_HOST", "localhost")
	viper.SetDefault("DATABASE.DB_PORT", "5432")
	viper.SetDefault("DATABASE.DB_USER", "postgres")
	viper.SetDefault("DATABASE.DB_PASSWORD", "postgres")
	viper.SetDefault("DATABASE.DB_NAME", "sembako_db")
	viper.SetDefault("DATABASE.DB_SSLMODE", "disable")
	viper.SetDefault("DATABASE.DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DATABASE.DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DATABASE.DB_CONN_MAX_LIFETIME", 15*time.Minute)

	viper.SetDefault("JWT.JWT_ACCESS_SECRET", "super-secret-access-key")
	viper.SetDefault("JWT.JWT_REFRESH_SECRET", "super-secret-refresh-key")
	viper.SetDefault("JWT.JWT_ACCESS_TTL", 15*time.Minute)
	viper.SetDefault("JWT.JWT_REFRESH_TTL", 7*24*time.Hour)

	viper.SetDefault("CORS.ALLOW_ORIGINS", []string{"*"})

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		// If specific env file is missing, try fallback .env
		viper.SetConfigName(".env")
		_ = viper.ReadInConfig()
	}

	getStringVal := func(dotKey string) string {
		flatKey := strings.ReplaceAll(dotKey, ".", "_")
		if val := viper.GetString(flatKey); val != "" {
			return val
		}
		return viper.GetString(dotKey)
	}

	var cfg Config
	cfg.App.Name = getStringVal("APP.NAME")
	cfg.App.Env = getStringVal("APP.ENV")
	cfg.App.Version = getStringVal("APP.VERSION")

	cfg.Server.Port = getStringVal("SERVER.PORT")
	cfg.Server.ReadTimeout = viper.GetDuration("SERVER.READ_TIMEOUT")
	cfg.Server.WriteTimeout = viper.GetDuration("SERVER.WRITE_TIMEOUT")

	cfg.Database.Host = getStringVal("DATABASE.DB_HOST")
	cfg.Database.Port = getStringVal("DATABASE.DB_PORT")
	cfg.Database.User = getStringVal("DATABASE.DB_USER")
	cfg.Database.Password = getStringVal("DATABASE.DB_PASSWORD")
	cfg.Database.Name = getStringVal("DATABASE.DB_NAME")
	cfg.Database.SSLMode = getStringVal("DATABASE.DB_SSLMODE")
	cfg.Database.MaxOpenConns = viper.GetInt("DATABASE.DB_MAX_OPEN_CONNS")
	cfg.Database.MaxIdleConns = viper.GetInt("DATABASE.DB_MAX_IDLE_CONNS")
	cfg.Database.ConnMaxLifetime = viper.GetDuration("DATABASE.DB_CONN_MAX_LIFETIME")

	cfg.JWT.AccessSecret = getStringVal("JWT.JWT_ACCESS_SECRET")
	cfg.JWT.RefreshSecret = getStringVal("JWT.JWT_REFRESH_SECRET")
	cfg.JWT.AccessTTL = viper.GetDuration("JWT.JWT_ACCESS_TTL")
	cfg.JWT.RefreshTTL = viper.GetDuration("JWT.JWT_REFRESH_TTL")

	return &cfg, nil
}
