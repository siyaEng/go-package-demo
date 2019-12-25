package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {

	dsn := "root:root@tcp(127.0.0.1:33067)/siya?timeout=90s&collation=utf8mb4_general_ci"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(db)
}
