package db

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func PerformSelect(dbConn *sqlx.DB, result interface{}, query string, params ...interface{}) (empty bool, err error) {
	log.Debug().Msgf("Execution query: %s, %v", query, params)

	err = dbConn.Select(result, query, params...)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to execute sql query.")
		return
	}
	num := reflect.ValueOf(result).Elem().Len()
	if num == 0 {
		empty = true
	}
	return
}

func PerformGet(dbConn *sqlx.DB, result interface{}, query string, params ...interface{}) (empty bool, err error) {
	log.Debug().Msgf("Execution query: %s, %v", query, params)

	err = dbConn.Get(result, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			empty = true
			err = nil
		} else {
			log.Debug().Err(err).Msg("Failed to execute sql query.")
		}
	}
	return
}

func PerformExec(dbConn *sqlx.DB, query string, params ...interface{}) (result sql.Result, err error) {
	result, err = dbConn.Exec(query, params...)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to execute sql query.")
	}
	return
}
