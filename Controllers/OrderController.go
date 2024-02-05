package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Order структура представляет собой модель данных для таблицы "Заказы"
type Order struct {
	OrderID   int           `json:"order_id"`
	BuyerID   int           `json:"buyer_id"`
	OrderDate time.Time     `json:"order_date"`
	CreatedBy int           `json:"created_by"`
	DeletedBy sql.NullInt64 `json:"deleted_by"`
	CreatedAt time.Time     `json:"created_at"`
	DeletedAt sql.NullTime  `json:"deleted_at"`
}

// GetUsers возвращает список всех заказов
func GetOrders(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	// Закрываем соединение с бд при закрытии функции
	defer db.Close()

	// Делаем запрос из бд и получаем данные в orders
	orders, err := db.Query("SELECT * FROM orders;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer orders.Close()

	// Создаем список из структур User
	var orderList []Order
	for orders.Next() {
		var order Order
		if err := orders.Scan(&order.OrderID, &order.BuyerID, &order.OrderDate, &order.CreatedBy,
			&order.DeletedBy, &order.CreatedAt, &order.DeletedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orderList = append(orderList, order)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, orderList)
}

// GetOrder возвращает заказа по ID
func GetOrder(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	orderID := params["id"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(orderID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID заказа"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM orders WHERE order_id = $1;", orderID)

	var order Order
	/* Сканирование копирует столбцы из сопоставленной строки в значения,
	на которые указывает dest. Смотрите документацию по строкам.Сканирование для
	получения подробной информации. Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	if err := row.Scan(&order.OrderID, &order.BuyerID, &order.OrderDate, &order.CreatedBy,
		&order.DeletedBy, &order.CreatedAt, &order.DeletedAt); err != nil {
		if err == sql.ErrNoRows {
			respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Заказ не найден"})
			return
		} else if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении заказа", "error1": err.Error()})
			return
		}
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, order)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при подключении к базе данных"})
		return
	}
	defer db.Close()

	// Получаем парраметры из запроса
	params := mux.Vars(r)
	Buyer_id := params["buyer_id"]
	Order_date := time.Now()
	Created_by := params["created_by"]

	// Вставка нового пользователя в базу данных
	_, err = db.Exec("INSERT INTO orders (buyer_id, order_date, created_by) VALUES ($1, $2, $3);", Buyer_id, Order_date, Created_by)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении заказа"})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Заказ создан"})
}

// UpdateUser обновляет данные пользователя по ID
func PutOrder(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	orderID := params["id"]
	what := params["what"]
	new := params["new"]

	// Проверка корректности значения параметра
	if _, err := strconv.Atoi(orderID); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID заказа"})
		return
	}

	// Запрос к бд
	_, err = db.Exec("UPDATE orders SET "+what+" = $1 WHERE order_id = $2", new, orderID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении заказа"})
		return
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Данные заказа обновлены"})
}

// DeleteUser удаляет пользователя по ID
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Инициализируем подключение к базе данных
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем идентификатор продукта из параметров запроса
	params := mux.Vars(r)
	orderID := params["id"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(orderID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID заказа"})
		return
	}

	// Затем удаляем пользователя
	_, err = db.Exec("DELETE FROM orders WHERE order_id = $1", orderID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении заказа"})
		return
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Заказ успешно удален"})
}
