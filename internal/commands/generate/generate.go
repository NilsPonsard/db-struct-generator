package generate

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/db-struct-generator/pkg/verbosity"
)

// setup ping command
func Generate(job *cli.Cmd) {

	var (
		user         = job.StringArg("USER", "", "DB username")
		password     string
		host         = job.StringArg("HOST", "", "DB host")
		port         = job.StringArg("PORT", "", "DB port")
		table        = job.StringArg("TABLE", "", "DB table")
		goName       = job.BoolOpt("g go-name", false, "Show only go compatible names")
		originalName = job.BoolOpt("n original-name", false, "Show only db original name")
		jsonTag      = job.BoolOpt("j json-tag", false, "Add json tagging")
		generateFunc = job.BoolOpt("f func", false, "Generate the function to retrieve data")
		fileOutput   = job.StringOpt("o output", "", "set output file")
	)

	// function to execute

	job.Action = func() {

		elements := strings.Split(*table, ".")

		tableName := elements[len(elements)-1]

		verbosity.Debug(*user, *host, *port, *table)

		// Ask for password so it’s not logged in password history

		fmt.Printf("%s@%s:%s password : ", *user, *host, *port)
		fmt.Scanln(&password)

		// connect to DB

		dbConn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *user, password, *host, *port, ""))

		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		// Get table

		rows, err := dbConn.Query("SELECT * FROM " + *table)
		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		// Get columns names

		columns, err := rows.Columns()
		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		// get columns types

		types, err := rows.ColumnTypes()
		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		// translate to go names (capitalise)

		var (
			goNames []string
		)
		for _, column := range columns {

			if len(column) > 1 {

				w := strings.ToUpper(string(column[0])) + column[1:]

				goNames = append(goNames, w)
			} else {

				goNames = append(goNames, strings.ToUpper(column))
			}

		}

		// check options

		if *goName {
			output(arrayToString(goNames), fileOutput)
			cli.Exit(0)
		}
		if *originalName {
			output(arrayToString(columns), fileOutput)

			cli.Exit(0)
		}

		// output struct

		out := "\ntype " + tableName + " struct {\n"

		for i := range columns {

			// translate type

			t := types[i].ScanType().Name()

			if strings.Contains(t, "int") {
				t = "int"
			} else if t == "RawBytes" {
				t = "string"
			}

			// generate line

			line := goNames[i] + " " + t

			if *jsonTag {
				line = line + " `json:\"" + columns[i] + "\"`"
			}

			line = line + "\n"

			out = out + line

		}

		out = out + "}\n"

		if *generateFunc {

			scanVars := ""

			for i, n := range columns {

				scanVars = scanVars + "&" + tableName + "." + n

				if i != len(columns)-1 {
					scanVars = scanVars + ", "
				}
			}

			out = `
import (
	"fmt"
	"database/sql"
)
` + out + `
func Get` + tableName + `(dbConn *sql.DB) (result []` + tableName + `, err error) {
	var row ` + tableName + `
	rows, err := dbConn.Query("SELECT ` + arrayToString(columns) + ` FROM ` + *table + `")
	if err != nil {
		return result, err
	}
			
	for rows.Next() {
							
		err := rows.Scan(` + scanVars + `)
			
		if err == nil {
			dnsDomains = append(dnsDomains, domain)
		} else {
			fmt.Println(err)
		}
	}
			
	return result, nil
}`

		}

		output(out, fileOutput)
	}
}

// print an array to go-like format
func arrayToString(arr []string) string {
	out := ""
	for i, n := range arr {

		out = out + n

		if i != len(arr)-1 {
			out = out + ","
		}
	}
	return out
}

// output to file/term
func output(content string, filePath *string) {

	verbosity.Info(content)

	if len(*filePath) > 0 {
		os.WriteFile(*filePath, []byte(content), 0600)
	}

}
