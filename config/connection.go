package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func Connect() (*pgxpool.Pool, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal("Error while creating connection to the databse!!")
	}
	conn, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer conn.Release()

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping the database!!")
	}
	log.Println("Connected to the database!!")

	return connPool, err
}
