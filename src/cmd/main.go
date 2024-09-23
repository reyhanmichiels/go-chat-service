package main

import (
	"context"
	"errors"

	"github.com/reyhanmichiels/go-pkg/auth"
	"github.com/reyhanmichiels/go-pkg/configreader"
	"github.com/reyhanmichiels/go-pkg/files"
	"github.com/reyhanmichiels/go-pkg/hash"
	"github.com/reyhanmichiels/go-pkg/log"
	"github.com/reyhanmichiels/go-pkg/parser"
	"github.com/reyhanmichiels/go-pkg/rate_limiter"
	"github.com/reyhanmichiels/go-pkg/redis"
	"github.com/reyhanmichiels/go-pkg/sql"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/domain"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/business/usecase"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/handler/rest"
	"github.com/reyhanmichies/go-rest-api-boiler-plate/src/utils/config"
)

// @contact.name   Reyhan Hafiz Rusyard
// @contact.email  michielsreyhan@gmail.com

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization

const (
	configfile   string = "./etc/cfg/conf.json"
	templatefile string = "./etc/tpl/conf.template.json"
	appnamespace string = ""
)

func main() {
	defaultLogger := log.DefaultLogger()

	// panic recovery
	defer func() {
		if err := recover(); err != nil {
			defaultLogger.Panic(err)
		}
	}()

	// TODO: need a way to build config file automatically, for now build the file manually
	if !files.IsExist(configfile) {
		defaultLogger.Fatal(context.Background(), errors.New("config file doesn't exist"))
	}

	// read config from config file
	cfg := config.Init()
	configReader := configreader.Init(configreader.Options{
		ConfigFile: configfile,
	})
	configReader.ReadConfig(&cfg)

	// init logger
	log := log.Init(cfg.Log)

	// init cache
	cache := redis.Init(cfg.Redis, log)

	// init db
	db := sql.Init(cfg.SQL, log)

	// init rate limiter
	rateLimiter := rate_limiter.Init(cfg.RateLimiter, log)

	// init parser
	parser := parser.InitParser(log, cfg.Parser)

	// init domain
	dom := domain.Init(domain.InitParam{Log: log, Db: db, Redis: cache, Json: parser.JSONParser()})

	// hash
	hash := hash.Init()

	// auth
	auth := auth.Init(cfg.Auth, log)

	// init usecase
	uc := usecase.Init(usecase.InitParam{Dom: dom, Log: log, Json: parser.JSONParser(), Hash: hash, Auth: auth})

	// init http server
	r := rest.Init(rest.InitParam{Uc: uc, GinConfig: cfg.Gin, Log: log, RateLimiter: rateLimiter, Json: parser.JSONParser(), Auth: auth})

	// run http server
	r.Run()
}
