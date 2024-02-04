package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User структура представляет собой модель данных для таблицы "Пользователи"
type Order struct {
	OrderId		int 	`json:"order_id"`
	BuyerId 	int 	`json:"buyer_id"`
	ProductId 	int 	`json:"product_id"`
	Quantity 	int 	`json:"quantity"`
	TotalPrice 	int 	`json:"total_price"`
	OrderDate 	int 	`json:"order_date"`
}


// GetUsers возвращает список всех пользователей
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
		if err := orders.Scan(&order.OrderId, &order.BuyerId, &order.ProductId, &order.Quantity, &order.TotalPrice, &order.OrderDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orderList = append(orderList, order)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, orderList)
}

// GetUser возвращает пользователя по ID
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
	err = row.Scan(&order.OrderId, &order.BuyerId, &order.ProductId, &order.Quantity, &order.TotalPrice, &order.OrderDate)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Заказ не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении заказа", "error1": err.Error()})
		return
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
	Quantity := params["quantity"]
	Total_price := params["total_price"]


	// Вставка нового пользователя в базу данных
	_, err = db.Exec("INSERT INTO orders (quantity, total_price) VALUES ($1, $2);", Quantity, Total_price)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении заказа"})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusCreated, map[string]string{"Код": "200"})
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
    _, err = db.Exec("UPDATE users SET " + what + " = $1 WHERE user_id = $2", new, orderID)
    if err != nil {
        respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении заказа"})
        return
    }

    // Передаем в функцию преобразования в json
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Данные пользователя обновлены"})
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
	log.Printf(orderID)

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(orderID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

// Затем удаляем пользователя
_, err = db.Exec("DELETE FROM users WHERE user_id = $1", orderID)
if err != nil {
    respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении пользователя"})
    return
}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Пользователь успешно удален"})
}
