package main

import (
	"os"
	"log"
	"net/http"

	"github.com/tkanos/gonfig"

	"github.com/Toor3-14/testProject/internal/app"
	"github.com/Toor3-14/testProject/internal/repo"
	"github.com/Toor3-14/testProject/internal/handler"
)


func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/redis/incr", handler.Incr )
	mux.HandleFunc("/sign/hmacsha512", handler.Hmacsha512)
	mux.HandleFunc("/postgres/users", handler.Users)
	return mux
}

// Create configure fields application struct (internal/app) and starting it
func main() {

	// Loggers - fields: infoLog, errLog
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	check := func (err error) {
		if err != nil {
			errLog.Fatal(err)
		}
	} 

	// Redis getting host/port from flags
	raddr, err := ParseFlags()
	check(err)

	// Create redis client - field: redis
	client, err := RedisClient(raddr)
	check(err)

	defer client.Close()


	// Getting configuration info (servAddr, postgre Addr,Username,Password) 
	// from resources/app.json
	json := &AppJSON{}
	err = gonfig.GetConf("../../resources/app.json", json)
	check(err)

	err = json.Validate()
	check(err)

	// Create pool of connections for PostgreSQL - field: userRepo
	dsn := "postgres://" + json.POSTGRES_USER_NAME + ":" + json.POSTGRES_PASSWORD + "@" + json.POSTGRES_ADDR
	db, err := PostgresPool(dsn)
	check(err)
	defer db.Close()

	// Configure application and check db's tables for exist
	app.Configure(
		json.SERVER_ADDR, 
		routes(),
		errLog, infoLog, 
		&repo.UserRepo{DB: db}, 
		client,
	)
	err = app.UserRepo().InitDB()
	check(err)

	infoLog.Printf("Statring with addr: %s", json.SERVER_ADDR)

	app.Start()
}