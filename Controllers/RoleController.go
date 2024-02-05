package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
)

// Role структура представляет собой модель данных для таблицы "Ролей"
type Role struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	// Дополнительные поля по необходимости
}

// GetProducts возвращает список всех продуктов
func GetRoles(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Закрываем соединение с бд при закрытии функции
	defer db.Close()

	// Делаем запрос из бд и получаем данные в products
	roles, err := db.Query("SELECT * FROM roles;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer roles.Close()

	// Создаем список из структур Product
	var rolesList []Role
	for roles.Next() {
		var role Role
		if err := roles.Scan(&role.RoleID, &role.RoleName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rolesList = append(rolesList, role)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, rolesList)
}


func GetRole(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем параметры из запроса
	params := mux.Vars(r)
	roleID := params["id"]

	// Проверка корректности значения параметра
	_, err = strconv.Atoi(roleID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID роли"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM roles WHERE roles_id = $1;", roleID)

	var role Role
	/* Сканирование копирует столбцы из сопоставленной строки в значения.
	Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&role.RoleID, &role.RoleName)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "роль не найдена"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, role)
}
