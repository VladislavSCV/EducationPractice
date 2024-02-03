package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Product структура представляет собой модель данных для таблицы "Товары"
type Product struct {
	ProductID   int       `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	Append_Date 	time.Time `json:"registration_date"`
}


// GetProducts возвращает список всех продуктов
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Закрываем соединение с бд при закрытии функции
	defer db.Close()

	// Делаем запрос из бд и получаем данные в products
	products, err := db.Query("SELECT * FROM products;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer products.Close()

	// Создаем список из структур Product
	var productList []Product
	for products.Next() {
		var product Product
		if err := products.Scan(&product.ProductID, &product.Name, 
			&product.Description, &product.Category, &product.Price, 
			&product.Category, &product.Append_Date); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		productList = append(productList, product)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, productList)
}


func GetProductById(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	productID := params["id"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(productID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM products WHERE product_id = $1;", productID)

	var product Product
	/* Сканирование копирует столбцы из сопоставленной строки в значения.
	Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&product.ProductID, &product.Name, &product.Description, 
		&product.Category, &product.Price, &product.Category, &product.Append_Date)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Пользователь не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, product)
}


// Поиск товаров по категориям
func GetProductByCategory(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	category := params["category"]
	log.Printf("[DEBUG] GetProductByCategory")

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM products WHERE category = $1;", category)

	var product Product
	/* Сканирование копирует столбцы из сопоставленной строки в значения.
	Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&product.ProductID, &product.Name, &product.Description, 
		&product.Category, &product.Price, &product.Category, &product.Append_Date)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Пользователь не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, product)
}


func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при подключении к базе данных"})
		return
	}
	defer db.Close()

	// Получаем парраметры из запроса
	params := mux.Vars(r)
	name := params["name"]
	desc := params["description"]
	categ := params["category"]
	price := params["price"]
	status := params["status"]

	// Делаем запрос к бд
	_, err = db.Exec("INSERT INTO products (name, description, category, price, status) VALUES ($1, $2, $3, $4, $5);", name, desc, categ, price, status)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении пользователя"})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusCreated, map[string]string{"Код": "200"})
}


func PutProduct(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса 
	params := mux.Vars(r)
	productID := params["id"]
	what := params["what"]
	new := params["new"]
	log.Printf(productID)

	// Проверка на валидность ID
	if _, err := strconv.Atoi(productID); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}
	// Запрос к бд
	_, err = db.Exec("UPDATE products SET "+what+" = $1 WHERE product_id = $2", new, productID)

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error updating user", "error_detail": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product data updated"})
}


// DeleteProduct удаляет продукт по ID
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Инициализируем подключение к базе данных
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем идентификатор продукта из параметров запроса
	params := mux.Vars(r)
	productID := params["id"]

	// Проверяем действительность идентификатора продукта
	_, err = strconv.Atoi(productID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
		return
	}

	// Удаляем продукт из базы данных
	_, err = db.Exec("DELETE FROM products WHERE product_id = $1", productID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error deleting product"})
		return
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product successfully deleted"})
}
