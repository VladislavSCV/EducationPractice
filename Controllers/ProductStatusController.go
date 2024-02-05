package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
)

// ProductStatus структура представляет собой модель данных для таблицы "Статусы товаров"
type ProductStatus struct {
	StatusID   int    `json:"status_id"`
	StatusName string `json:"status_name"`
}

// GetProducts возвращает список всех продуктов
func GetProductsStatus(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Закрываем соединение с бд при закрытии функции
	defer db.Close()

	// Делаем запрос из бд и получаем данные в products
	statuses, err := db.Query("SELECT * FROM productstatuses;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statuses.Close()

	// Создаем список из структур Product
	var statusList []ProductStatus
	for statuses.Next() {
		var status ProductStatus
		if err := statuses.Scan(&status.StatusID, &status.StatusName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statusList = append(statusList, status)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, statusList)
}


func GetProductStatus(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	StatusID := params["id"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(StatusID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID роли"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM productstatuses WHERE status_id = $1;", StatusID)

	var status TradeStatus
	/* Сканирование копирует столбцы из сопоставленной строки в значения.
	Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&status.StatusID, &status.StatusName)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "статус не найдена"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, status)
}
