package database

import "github.com/gocraft/dbr"

var (
	Persontable = "person"
	Listtable   = "list"
	Movietable  = "movie"
	Tablename   = "user"
	seq         = 1
	conn, _     = dbr.Open("mysql", "root:@tcp(localhost:3306)/netflix", nil)
	Sess        = conn.NewSession(nil)
)
