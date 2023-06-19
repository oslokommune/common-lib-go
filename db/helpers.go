package db

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"reflect"
)

func PerformSelect(dbConn *PostgresConnection, result interface{}, query string, params ...interface{}) error {
	err := dbConn.DBPool.Select(result, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute sql query. ")
		return err
	}
	num := reflect.ValueOf(result).Elem().Len()
	if num == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func PerformGet(dbConn *PostgresConnection, result interface{}, query string, params ...interface{}) error {
	err := dbConn.DBPool.Get(result, query, params...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msg("Failed to execute sql query. ")
		}
		return err
	}
	return nil
}

func PerformExec(dbConn *PostgresConnection, query string, params ...interface{}) error {
	_, err := dbConn.DBPool.Exec(query, params...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msg("Failed to execute sql query. ")
		}
		return err
	}
	return nil
}
