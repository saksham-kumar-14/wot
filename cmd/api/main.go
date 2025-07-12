package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/saksham-kumar-14/wot/internal/auth"
	"github.com/saksham-kumar-14/wot/internal/db"
	"github.com/saksham-kumar-14/wot/internal/env"
	"github.com/saksham-kumar-14/wot/internal/mailer"
	ratelimiter "github.com/saksham-kumar-14/wot/internal/rateLimiter"
	"github.com/saksham-kumar-14/wot/internal/store"
	"github.com/saksham-kumar-14/wot/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.1"

// @title wot
// @description web app for discussing academics and related
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in 						   header
// @name 					   Authorization
// @description

func main() {

	cfg := config{
		addr:        env.GetString("PORT", ":8000"),
		apiURL:      env.GetString("EXTERNAL_URL", ":8000"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:strongpwd@localhost:5432/wot?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 48,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			token: tokenConfig{
				secret: env.GetString("AUTH_SECRET", "super_strong_secret"),
				exp:    time.Hour * 48,
				issuer: "wot",
			},
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", true),
		},
		ratelimiter: ratelimiter.Config{
			ReqPerTimeFrame: env.GetInt("RATELIMITER_REQ_COUNT", 20),
			TimeFrame:       time.Second * 5,
			Enabled:         env.GetBool("RATELIMITER_ENABLED", true),
		},
	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB connected")

	// Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
	}

	// Rate limiter
	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.ratelimiter.ReqPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	store := store.NewDbStorage(db)
	cacheStore := cache.NewRedisStorage(rdb)

	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	jwtAuth := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.issuer, cfg.auth.token.issuer)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStore,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuth,
		ratelimiter:   ratelimiter,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
