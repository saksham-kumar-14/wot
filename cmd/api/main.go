package main

import (
	"time"

	"github.com/saksham-kumar-14/wot/internal/db"
	"github.com/saksham-kumar-14/wot/internal/env"
	"github.com/saksham-kumar-14/wot/internal/mailer"
	"github.com/saksham-kumar-14/wot/internal/store"
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

	store := store.NewDbStorage(db)

	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
