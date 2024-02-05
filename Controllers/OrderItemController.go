package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// OrderItem структура представляет собой модель данных для таблицы "Товары в Заказе"
type OrderItem struct {
	OrderItemID  int            `json:"order_item_id"`
	OrderID      int            `json:"order_id"`
	ProductID    int            `json:"product_id"`
	Quantity     int            `json:"quantity"`
	TotalPrice   float64        `json:"total_price"`
	CreatedBy    int            `json:"created_by"`
	DeletedBy    sql.NullInt64  `json:"deleted_by"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    sql.NullTime   `json:"deleted_at"`
}


// GetUsers возвращает список всех пользователей
func GetItems(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	// Закрываем соединение с бд при закрытии функции
	defer db.Close()

	// Делаем запрос из бд и получаем данные в users
	items, err := db.Query("SELECT * FROM orderitems;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer items.Close()

	// Создаем список из структур User
	var itemList []OrderItem
	for items.Next() {
		var item OrderItem
		if err := items.Scan(&item.OrderItemID, &item.OrderID, &item.ProductID, 
			&item.Quantity, &item.TotalPrice, &item.CreatedBy, &item.DeletedBy, 
			&item.CreatedAt, &item.DeletedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		itemList = append(itemList, item)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, itemList)
}

// GetUser возвращает пользователя по ID
func GetItem(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	ID := params["id"]

	// Проверка корректности значения параметра
	id, err := strconv.Atoi(ID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM orderitems WHERE order_item_id = $1;", id)

	var item OrderItem
	err = row.Scan(&item.OrderItemID, &item.OrderID, &item.ProductID, 
		&item.Quantity, &item.TotalPrice, &item.CreatedBy, &item.DeletedBy, 
		&item.CreatedAt, &item.DeletedAt)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Товар в заказе не найден не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении пользователя", "error1": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, item)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при подключении к базе данных"})
		return
	}
	defer db.Close()

	// Получаем парраметры из запроса
	params := mux.Vars(r)
	Order_id := params["order_id"]
	Product_id := params["product_id"]
	Quantity := params["quantity"]
	Total_price := params["total_price"]
	Created_by := params["created_by"]

	// Вставка нового пользователя в базу данных
	_, err = db.Exec("INSERT INTO OrderItems (order_id, product_id, quantity, total_price, created_by) VALUES ($1, $2, $3, $4, $5);", Order_id, Product_id, Quantity, Total_price, Created_by)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении товара в заказ"})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusCreated, map[string]string{"Код": "200"})
}


// UpdateUser обновляет данные пользователя по ID
func PutItem(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
    db, err := initDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	// Получаем параметры из запроса 
    params := mux.Vars(r)
    ID := params["id"]
    what := params["what"]
    new := params["new"]

    // Проверка корректности значения параметра
    if _, err := strconv.Atoi(ID); err != nil {
        respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
        return
    }

	// Запрос к бд
    _, err = db.Exec("UPDATE orderitems SET " + what + " = $1 WHERE order_item_id = $2", new, ID)
    if err != nil {
        respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении товара в заказе"})
        return
    }

    // Передаем в функцию преобразования в json
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Данные заказа обновлены"})
}

// DeleteUser удаляет пользователя по ID
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Инициализируем подключение к базе данных
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем идентификатор продукта из параметров запроса
	params := mux.Vars(r)
	ID := params["id"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(ID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID товара"})
		return
	}

	// Затем удаляем пользователя
	_, err = db.Exec("DELETE FROM orderitems WHERE order_item_id = $1", ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении товара"})
		return
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Товар успешно удален"})
}
