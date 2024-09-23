package config

import (
	"time"

	"github.com/reyhanmichiels/go-pkg/auth"
	"github.com/reyhanmichiels/go-pkg/log"
	"github.com/reyhanmichiels/go-pkg/parser"
	"github.com/reyhanmichiels/go-pkg/rate_limiter"
	"github.com/reyhanmichiels/go-pkg/redis"
	"github.com/reyhanmichiels/go-pkg/sql"
	"github.com/reyhanmichiels/go-pkg/translator"
)

type Application struct {
	Meta        ApplicationMeta
	Gin         GinConfig
	Log         log.Config
	SQL         sql.Config
	Auth        auth.Config
	Redis       redis.Config
	Translator  translator.Config
	RateLimiter rate_limiter.Config
	Parser      parser.Options
}

type ApplicationMeta struct {
	Title       string
	Description string
	Host        string
	BasePath    string
	Version     string
}

type GinConfig struct {
	Port            string
	Mode            string
	LogRequest      bool
	LogResponse     bool
	Timeout         time.Duration
	ShutdownTimeout time.Duration
	CORS            CORSConfig
	Meta            ApplicationMeta
	Swagger         SwaggerConfig
	Dummy           DummyConfig
}

type CORSConfig struct {
	Mode string
}
type SwaggerConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type PlatformConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type DummyConfig struct {
	Enabled bool
	Path    string
}

type BasicAuthConf struct {
	Username string
	Password string
}

func Init() Application {
	return Application{}
}
