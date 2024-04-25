package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME    string = "root"
	PASSWORD    string = "wenliang75"
	SERVER      string = "127.0.0.1"
	PORT        int    = 3306
	DATABASE    string = "testsql"
	createDB           = "CREATE DATABASE IF NOT EXISTS testmysqlDB"
	createTable        = "CREATE TABLE IF NOT EXISTS testtable (id INT AUTO_INCREMENT PRIMARY KEY, item VARCHAR(255), createdatetime DATETIME)"
)

func main() {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", USERNAME, PASSWORD, SERVER, PORT)

	db, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println("MySQL connectin has some problem, the reson is: ", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("DataBase connection has some problems, the error is: ", err.Error())
	}
	defer db.Close()
	_, err = db.Exec(createDB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database created successfully or already exist.")

	_, err = db.Exec("USE testmysqlDB")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created successfully or already exists.")

	_, err = db.Exec("INSERT INTO testtable (item,createdatetime) VALUES (?,?)", "example item", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}

	// 查询数据
	rows, err := db.Query("SELECT item FROM testtable")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Item:", item)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
