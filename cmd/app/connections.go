package main

import (
	"context"
	"flag"
	"time"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v4/pgxpool"
)
type AppJSON struct {
	SERVER_ADDR string

	POSTGRES_ADDR string
	POSTGRES_USER_NAME string
	POSTGRES_PASSWORD string
}
func (a AppJSON) Validate() error {
	if a.SERVER_ADDR == "" {
		return errors.New("Config Error: field `SERVER_ADDR` is not specified in app.json")
	}else if a.POSTGRES_ADDR == "" {
		return errors.New("Config Error: field `POSTGRES_ADDR` is not specified in app.json")
	}else if a.POSTGRES_USER_NAME == "" {
		return errors.New("Config Error: field `POSTGRES_USER_NAME` is not specified in app.json")
	}else if a.POSTGRES_PASSWORD == "" {
		return errors.New("Config Error: field `POSTGRES_PASSWORD` is not specified in app.json")
	}
	return nil
}
func ParseFlags() (string, error) {
	// Redis getting host/port and create connection
	host := flag.String("host", "", "Redis host")
	port := flag.String("port", "", "Redis port")
	flag.Parse()

	if *host == "" || *port == "" {
		return "", errors.New("Empty fields -host/-port for Redis connection")
	}

	return *host + ":" + *port, nil
}
func RedisClient(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:	  addr,
		Password: "",
		DB:		  0,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return client, nil
}
func PostgresPool(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	db, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}
	return db, nil
}