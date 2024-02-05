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


	/* Endpoints для таблицы с пользователями */

	// Получение всех пользователей
	router.HandleFunc("/api/users", 
	Controllers.GetUsers).Methods("GET")

	// Получить пользователя по id
	router.HandleFunc("/api/user/{id}", 
	Controllers.AuthUser).Methods("GET")

	// Проверка на существование в бд пользователя (Авторизация пользователя)
	router.HandleFunc("/api/user/{email}/{password}", 
	Controllers.AuthUser).Methods("GET")

	// Создать нового пользователя (Регистсрация пользователя)
	router.HandleFunc("/api/user/{username}/{email}/{password}", 
	Controllers.CreateUser).Methods("POST")
	
	// Обновить данные пользователя
	router.HandleFunc("/api/user/{id}/{what}/{new}", 
	Controllers.PutUser).Methods("PUT")

	// Удалить пользователя
	router.HandleFunc("/api/user/{id}", 
	Controllers.DeleteUser).Methods("DELETE")



	/* Endpoints для таблицы с товарами */

	// Получить все товары
	router.HandleFunc("/api/products", 
	Controllers.GetProducts).Methods("GET")

	// Получить товар по категории
	router.HandleFunc("/api/productCategory/{category}", 
	Controllers.GetProductByCategory).Methods("GET")

	// Получить товар по id
	router.HandleFunc("/api/productId/{id}", 
	Controllers.GetProductById).Methods("GET")

	// Создать новый товар (Регистрация товара)
	router.HandleFunc("/api/product/{name}/{description}/{category}/{price}/{status}", 
	Controllers.CreateProduct).Methods("POST")

	// Обновить данные о товаре
	router.HandleFunc("/api/product/{id}/{what}/{new}", 
	Controllers.PutProduct).Methods("PUT")

	// Удалить товар
	router.HandleFunc("/api/product/{id}", 
	Controllers.DeleteProduct).Methods("DELETE")


	/* Endpoints для таблицы с заказами */

	// Получить все товары
	router.HandleFunc("/api/orders", 
	Controllers.GetOrders).Methods("GET")

	// Получить заказ по id
	router.HandleFunc("/api/order/{id}", 
	Controllers.GetOrder).Methods("GET")

	// Создать новый заказ (Регистрация заказа)
	router.HandleFunc("/api/order/{buyer_id}/{product_id}/{quantity}/{total_price}", 
	Controllers.CreateOrder).Methods("POST")

	// Обновить данные о товаре
	router.HandleFunc("/api/order/{id}/{what}/{new}", 
	Controllers.PutOrder).Methods("PUT")

	// Удалить товар
	router.HandleFunc("/api/order/{id}", 
	Controllers.DeleteOrder).Methods("DELETE")


	/* Endpoints для таблицы с ролями */

	// Получить все товары
	router.HandleFunc("/api/roles", 
	Controllers.GetRoles).Methods("GET")

	// Получить заказ по id
	router.HandleFunc("/api/role/{id}", 
	Controllers.GetRole).Methods("GET")

	// Запуск сервера на локальном порту 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}