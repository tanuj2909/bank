package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/tanuj2909/bank/api"
	db "github.com/tanuj2909/bank/db/sqlc"
	"github.com/tanuj2909/bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
