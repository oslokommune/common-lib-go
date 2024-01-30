package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

const (
	_defaultMaxPoolSize = 1
)

type PostgresConnection struct {
	sqlDB *sql.DB

	maxPoolSize int
}

func New(conf DbConf, opts ...Option) *PostgresConnection {
	p := &PostgresConnection{
		maxPoolSize: _defaultMaxPoolSize,
	}

	dbUserName := conf.Username
	dbPassword := conf.Password
	dbHost := conf.Host
	dbName := conf.Database
	dbPort := conf.Port

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUserName, dbPassword, dbName)
	sqlDB, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open postgres connection")
	}

	// Custom options
	for _, opt := range opts {
		opt(p)
	}

	sqlDB.SetMaxOpenConns(p.maxPoolSize)

	p.sqlDB = sqlDB
	return p
}

func (p *PostgresConnection) Connection() *sql.DB {
	return p.sqlDB
}

func (p *PostgresConnection) CloseConnection() {
	err := p.sqlDB.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close postgres connection")
	}
}
