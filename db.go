package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type RDSClient struct {
	db *sql.DB
}

// NewRDSStorage crea una conexión a una base de datos MySQL en AWS RDS
func NewRDSStorage(user, password, host, port, dbName string) (*RDSClient, error) {
	// Cadena de conexión: usuario:contraseña@tcp(host:puerto)/nombreBD
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	// Conecta a la base de datos
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %v", err)
	}

	// Establece parámetros de conexión
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Verifica la conexión
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error al hacer ping: %v", err)
	}

	log.Println("RDS: Successfully connected")
	return &RDSClient{db: db}, nil
}

func (r *RDSClient) GetDB() *sql.DB {
	return r.db
}
