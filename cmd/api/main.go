package main

import (
	"log"

	"github.com/saksham-kumar-14/wot/internal/db"
	"github.com/saksham-kumar-14/wot/internal/env"
	"github.com/saksham-kumar-14/wot/internal/store"
)

func main() {

	cfg := config{
		addr: env.GetString("PORT", ":8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:strongpwd@localhost:5432/wot?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("DB connected")

	store := store.NewDbStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
