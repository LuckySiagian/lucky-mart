package config

import (
	"errors"
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
	Storage  StorageConfig  `mapstructure:"storage"`
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
	URL             string        `mapstructure:"DB_URL"`
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

// StorageConfig says where uploaded files (product photos, transfer proofs)
// are kept on disk. They are served back under /uploads/images.
type StorageConfig struct {
	UploadDir string `mapstructure:"UPLOAD_DIR"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"ALLOW_ORIGINS"`
}

const (
	defaultAccessSecret  = "dev-only-access-secret-change-me"
	defaultRefreshSecret = "dev-only-refresh-secret-change-me"
)

// LoadConfig reads .env.<APP_ENV> (falling back to .env) and then real
// environment variables, which always win. Variable names are listed in
// .env.example. On hosting platforms the conventional PORT and DATABASE_URL
// variables are honoured too.
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(fmt.Sprintf(".env.%s", env))
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

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

	viper.SetDefault("JWT.JWT_ACCESS_SECRET", defaultAccessSecret)
	viper.SetDefault("JWT.JWT_REFRESH_SECRET", defaultRefreshSecret)
	viper.SetDefault("JWT.JWT_ACCESS_TTL", 15*time.Minute)
	viper.SetDefault("JWT.JWT_REFRESH_TTL", 7*24*time.Hour)

	viper.SetDefault("CORS.ALLOW_ORIGINS", "*")

	viper.SetDefault("STORAGE.UPLOAD_DIR", "uploads")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigName(".env")
		_ = viper.ReadInConfig()
	}

	getString := func(dotKey string) string {
		flatKey := strings.ReplaceAll(dotKey, ".", "_")
		if val := viper.GetString(flatKey); val != "" {
			return val
		}
		return viper.GetString(dotKey)
	}
	getDuration := func(dotKey string) time.Duration {
		if raw := getString(dotKey); raw != "" {
			if d, err := time.ParseDuration(raw); err == nil {
				return d
			}
		}
		return viper.GetDuration(dotKey)
	}
	getInt := func(dotKey string) int {
		flatKey := strings.ReplaceAll(dotKey, ".", "_")
		if viper.IsSet(flatKey) {
			return viper.GetInt(flatKey)
		}
		return viper.GetInt(dotKey)
	}

	var cfg Config
	cfg.App.Name = getString("APP.NAME")
	cfg.App.Env = getString("APP.ENV")
	cfg.App.Version = getString("APP.VERSION")

	cfg.Server.Port = getString("SERVER.PORT")
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	cfg.Server.ReadTimeout = getDuration("SERVER.READ_TIMEOUT")
	cfg.Server.WriteTimeout = getDuration("SERVER.WRITE_TIMEOUT")

	cfg.Database.URL = os.Getenv("DATABASE_URL")
	cfg.Database.Host = getString("DATABASE.DB_HOST")
	cfg.Database.Port = getString("DATABASE.DB_PORT")
	cfg.Database.User = getString("DATABASE.DB_USER")
	cfg.Database.Password = getString("DATABASE.DB_PASSWORD")
	cfg.Database.Name = getString("DATABASE.DB_NAME")
	cfg.Database.SSLMode = getString("DATABASE.DB_SSLMODE")
	cfg.Database.MaxOpenConns = getInt("DATABASE.DB_MAX_OPEN_CONNS")
	cfg.Database.MaxIdleConns = getInt("DATABASE.DB_MAX_IDLE_CONNS")
	cfg.Database.ConnMaxLifetime = getDuration("DATABASE.DB_CONN_MAX_LIFETIME")

	cfg.JWT.AccessSecret = getString("JWT.JWT_ACCESS_SECRET")
	cfg.JWT.RefreshSecret = getString("JWT.JWT_REFRESH_SECRET")
	cfg.JWT.AccessTTL = getDuration("JWT.JWT_ACCESS_TTL")
	cfg.JWT.RefreshTTL = getDuration("JWT.JWT_REFRESH_TTL")

	cfg.CORS.AllowOrigins = splitList(getString("CORS.ALLOW_ORIGINS"))

	cfg.Storage.UploadDir = getString("STORAGE.UPLOAD_DIR")

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func splitList(raw string) []string {
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if p := strings.TrimSpace(strings.Trim(part, `"`)); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// validate refuses to start production with settings that would be unsafe
// on the public internet.
func (c *Config) validate() error {
	if c.App.Env != "production" {
		return nil
	}

	var problems []string
	if c.JWT.AccessSecret == defaultAccessSecret || len(c.JWT.AccessSecret) < 32 {
		problems = append(problems, "JWT_JWT_ACCESS_SECRET must be set to a random value of at least 32 characters")
	}
	if c.JWT.RefreshSecret == defaultRefreshSecret || len(c.JWT.RefreshSecret) < 32 {
		problems = append(problems, "JWT_JWT_REFRESH_SECRET must be set to a random value of at least 32 characters")
	}
	if c.JWT.AccessSecret == c.JWT.RefreshSecret {
		problems = append(problems, "access and refresh JWT secrets must differ")
	}
	for _, o := range c.CORS.AllowOrigins {
		if o == "*" {
			problems = append(problems, "CORS_ALLOW_ORIGINS must list the real frontend URLs, not *")
			break
		}
	}
	if len(problems) > 0 {
		return errors.New("unsafe production configuration: " + strings.Join(problems, "; "))
	}
	return nil
}
