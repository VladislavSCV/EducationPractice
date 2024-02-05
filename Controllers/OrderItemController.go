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
	users, err := db.Query("SELECT * FROM users;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer users.Close()

	// Создаем список из структур User
	var userList []User
	for users.Next() {
		var user User
		if err := users.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.RegDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userList = append(userList, user)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, userList)
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
	userID := params["id"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(userID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM users WHERE user_id = $1;", userID)

	var user User
	/* Сканирование копирует столбцы из сопоставленной строки в значения,
	на которые указывает dest. Смотрите документацию по строкам.Сканирование для
	получения подробной информации. Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.RegDate)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Пользователь не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении пользователя", "error1": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, user)
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
	UserName := params["username"]
	Email := params["email"]
	Password := params["password"]

	// Вставка нового пользователя в базу данных
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3);", UserName, Email, Password)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении пользователя"})
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
    userID := params["id"]
    what := params["what"]
    new := params["new"]

    // Проверка корректности значения параметра
    if _, err := strconv.Atoi(userID); err != nil {
        respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
        return
    }

	// Запрос к бд
    _, err = db.Exec("UPDATE users SET " + what + " = $1 WHERE user_id = $2", new, userID)
    if err != nil {
        respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении пользователя", "erro2": err.Error()})
        return
    }

    // Передаем в функцию преобразования в json
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Данные пользователя обновлены"})
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
	userID := params["id"]
	log.Printf(userID)

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(userID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

	// Затем удаляем пользователя
	_, err = db.Exec("DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении пользователя"})
		return
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Пользователь успешно удален"})
}
