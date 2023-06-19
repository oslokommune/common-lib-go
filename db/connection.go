package db

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"os"
)

type PostgresConnection struct {
	conf   *DbConf
	DBPool *sqlx.DB
}

func NewPostgresConnection(conf *DbConf) *PostgresConnection {
	return &PostgresConnection{conf: conf}
}

func (p *PostgresConnection) ConnectToDB() {
	dbUserName := p.conf.Username
	dbPassword := p.conf.Password
	dbHost := p.conf.Host
	dbName := p.conf.Database
	dbPort := p.conf.Port

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUserName, dbPassword, dbName)

	dbpool, err := sqlx.Connect("pgx", psqlInfo)

	if err != nil {
		log.Error().Msgf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	dbpool.SetMaxOpenConns(10)
	dbpool.SetMaxIdleConns(10)

	p.DBPool = dbpool
}

func (p *PostgresConnection) CloseConnection() {
	err := p.DBPool.Close()
	if err != nil {
		log.Error().Msgf("Failed to close database connection %v", err)
	}
}
