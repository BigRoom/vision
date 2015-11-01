package models

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
)

func init() {
	dat.EnableInterpolation = true
	dat.Strict = false
}

// DB is the database instance used be the models
var DB *runner.DB

// Init creates a new connection to a database
func Init(user, pass, service, port, name string) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
			user,
			pass,
			service,
			port,
			name,
		),
	)

	if err != nil {
		fmt.Println("lol")
		panic(err)
	}

	runner.MustPing(db)

	runner.LogQueriesThreshold = 10 * time.Millisecond
	DB = runner.NewDB(db, "postgres")

	DB.DB.MustExec(query)
}

var query = `
CREATE TABLE IF NOT EXISTS pings (
	id serial PRIMARY KEY,
	message text NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	username text NOT NULL,
	password text NOT NULL,
	email text NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
	id serial PRIMARY KEY,
	message text NOT NULL,
	username text NOT NULL,
	channel_key text NOT NULL,
	time bigserial NOT NULL
);
`
