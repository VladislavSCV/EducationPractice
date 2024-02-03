package Controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

var db *sql.DB

// initDB инициализирует базу данных ивозвращает объект DB, err 
func initDB() (*sql.DB, error) {
	/* Строка используемая для подключения к бд. 
	В ней сожержатся все основные параметры для подключения к бд извне */
	connString := "user=postgres password=31415926 dbname=SportMaster sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	// Проверяем подключение
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Вспомогательная функция для отправки JSON-ответа
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Преобразовываем входные данные в json
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Устанавливем заголовок ответа и записываем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}