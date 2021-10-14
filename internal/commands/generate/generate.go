package generate

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/db-struct-generator/pkg/verbosity"
)

// setup ping command
func Generate(job *cli.Cmd) {

	var (
		user     = job.StringArg("USER", "", "DB username")
		password string
		host     = job.StringArg("HOST", "", "DB host")
		port     = job.StringArg("PORT", "", "DB port")
		table    = job.StringArg("TABLE", "", "DB table")
	)

	// function to execute

	job.Action = func() {

		verbosity.Debug(*user, *host, *port, *table)

		fmt.Printf("%s@%s:%s password : ", *user, *host, *port)
		fmt.Scanln(&password)

		dbConn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *user, password, *host, *port, ""))

		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		rows, err := dbConn.Query("SELECT * FROM " + *table)
		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		verbosity.Info(rows.Columns())

	}
}
