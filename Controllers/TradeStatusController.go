package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
)

// TradeStatus структура представляет собой модель данных для таблицы "Статусы обмена"
type TradeStatus struct {
	StatusID   int    `json:"status_id"`
	StatusName string `json:"status_name"`
}

// GetProducts возвращает список всех продуктов
func GetTradesStatus(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Закрываем соединение с бд при закрытии функции
	defer db.Close()

	// Делаем запрос из бд и получаем данные в products
	statuses, err := db.Query("SELECT * FROM tradestatuses;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statuses.Close()

	// Создаем список из структур Product
	var statusesList []TradeStatus
	for statuses.Next() {
		var status TradeStatus
		if err := statuses.Scan(&status.StatusID, &status.StatusName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statusesList = append(statusesList, status)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, statusesList)
}


func GetTradeStatus(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	ID := params["status"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(ID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID статуса"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT status_name FROM tradestatuses WHERE status_name = $1;", ID)

	var status TradeStatus
	/* Сканирование копирует столбцы из сопоставленной строки в значения.
	Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&status.StatusID, &status.StatusName)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "статус не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, status)
}
