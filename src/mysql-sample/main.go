package main

import(
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//http://go-database-sql.org/accessing.html
func main(){
	db, err := sql.Open("mysql",
		"root:myPassw0rd@tcp(127.0.0.1:3306)/testgo")
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//Test get db data
	var (
		id int
		name string
	)

	rows, err := db.Query("select id, name from new_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}