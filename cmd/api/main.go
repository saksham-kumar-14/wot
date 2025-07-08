package main

import (
	"log"

	"github.com/saksham-kumar-14/wot/internal/env"
)

func main() {

	cfg := config{
		addr: env.GetString("PORT", ":8000"),
	}
	app := &application{
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
