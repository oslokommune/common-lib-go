package db

import (
	"database/sql"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func PerformSelect(dbConn *sqlx.DB, result interface{}, query string, params ...interface{}) error {
	err := dbConn.Select(result, query, params...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute sql query.")
		return err
	}
	num := reflect.ValueOf(result).Elem().Len()
	if num == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func PerformGet(dbConn *sqlx.DB, result interface{}, query string, params ...interface{}) error {
	err := dbConn.Get(result, query, params...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msg("Failed to execute sql query.")
		}
		return err
	}
	return nil
}

func PerformExec(dbConn *sqlx.DB, query string, params ...interface{}) error {
	_, err := dbConn.Exec(query, params...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msg("Failed to execute sql query.")
		}
		return err
	}
	return nil
}
