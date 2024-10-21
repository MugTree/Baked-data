package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed data.db
var embeddedDB []byte

func main() {

	// In serverless environments, use /tmp for the temporary database file
	//tempFile, err := os.Create("/tmp/embeddedDB.db")

	tempFile, err := os.CreateTemp("", "embeddedDB-*.db")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tempFile.Name()) // Clean up the temp file later

	if _, err := tempFile.Write(embeddedDB); err != nil {
		fmt.Println("Error writing to temp file:", err)
		return
	}

	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	type Product struct {
		Id   int
		Name string
	}

	prods := []Product{}
	prod := Product{}

	rows, err := db.Query("SELECT * from product;")
	if err != nil {
		fmt.Println("Error querying database:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {

		err := rows.Scan(&prod.Id, &prod.Name)
		if err != nil {
			log.Fatal(err)
		}
		prods = append(prods, prod)
	}

	//simple webserver

	// type Product struct {
	// 	Id   int
	// 	Name string
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		for _, v := range prods {
			fmt.Fprint(w, v)
		}

		//fmt.Fprintf(w, "hey!")

	})

	log.Fatal(http.ListenAndServe(":3000", nil))

}
