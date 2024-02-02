package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


// Product структура представляет собой модель данных для таблицы "Товары"
type Product struct {
	ProductID   int       `json:"product_id"`
	SellerID    int       `json:"seller_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
}

// GetUsers returns a list of all users
func GetProduct(w http.ResponseWriter, r *http.Request) {
    users, err := db.Query("SELECT * FROM users;")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer users.Close()

    var userList []User
    for users.Next() {
        var user User
        if err := users.Scan(&user.UserID, &user.Username, &user.Email, &user.Password); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        userList = append(userList, user)
    }

    respondWithJSON(w, http.StatusOK, userList)
}

// GetUser возвращает пользователя по ID
func GetProducts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]
	log.Printf("user: %v", userID)

	// Проверка корректности значения параметра
	_, err := strconv.Atoi(userID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE user_id = $1;", userID)

	var user User
	err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Пользователь не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении пользователя"})
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}