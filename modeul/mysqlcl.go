package modeul

import (
	"database/sql"
)

type MySQLCl struct {
	db *sql.DB
}

var db = MySQLCl{}
