package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

const queryGetColAndType = `select column_name, data_type from information_schema.columns where table_schema = $1 AND table_name = $2`

func (dbStructures *dbStructures) Select(db *sqlx.DB, dbStruct dbStruct) (err error) {
	if err = db.Select(dbStructures, queryGetColAndType, dbStruct.Schema, dbStruct.TableName); err != nil {
		log.Fatalln(err)
	}
	return
}
