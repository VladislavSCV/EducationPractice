package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Замените следующие значения на свои реальные данные
	host := "your_host"
	port := 5432
	user := "your_username"
	password := "your_password"
	dbname := "your_database_name"

	// Строка подключения
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL")

	// Пример выполнения запроса
	rows, err := db.Query("SELECT * FROM YourTable")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Обработка результатов запроса
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
}
