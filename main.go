package main

import (
	"log"
	"net/http"

	"github.com/VladislavSCV/EducationPractice/Controllers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)



func main() {
	log.Print("Server starting")

	router := mux.NewRouter()

	// Endpoints для таблицы с пользователями

	// Получение всех пользователей
	router.HandleFunc("/api/users", Controllers.GetUsers).Methods("GET")

	// Получить пользователя по id
	router.HandleFunc("/api/user/{id}", Controllers.GetUser).Methods("GET")

	// Создать нового пользователя
	router.HandleFunc("/api/user/{username}/{email}/{password}", Controllers.CreateUser).Methods("POST")

	// Обновить данные пользователя
	router.HandleFunc("/api/user/{id}/{what}/{new}", Controllers.PutUser).Methods("PUT")

	// Удалить пользователя
	router.HandleFunc("/api/user/{id}", Controllers.DeleteUser).Methods("DELETE")

	// Endpoints для таблицы с товарами

	// Получить все товары
	router.HandleFunc("/api/products", Controllers.GetProducts).Methods("GET")

	// Получить товар по id
	router.HandleFunc("/api/product/{id}", Controllers.GetProduct).Methods("GET")

	// Создать новый товар
	router.HandleFunc("/api/product/{name}/{description}/{category}/{price}/{status}", Controllers.CreateProduct).Methods("POST")

	// Обновить данные о товаре
	router.HandleFunc("/api/product/{id}/{what}/{new}", Controllers.PutProduct).Methods("PUT")

	// Удалить товар
	router.HandleFunc("/api/product/{id}", Controllers.DeleteProduct).Methods("DELETE")

	// Запустить сервер на порту 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}
}
