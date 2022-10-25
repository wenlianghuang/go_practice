package mysqlselect

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type mySalary struct {
	Id     int
	Salary int
}

func mysqlselect() {
	db, err := sql.Open("mysql", "Matt:matt042275@tcp(127.0.0.1:3306)/testdb")
	defer db.Close()

}
