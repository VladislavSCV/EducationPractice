package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// User структура представляет собой модель данных для таблицы "Пользователи"
type User struct {
	UserID           int       `json:"user_id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	RegistrationDate time.Time `json:"registration_date"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("28012006")))

// GetUsers возвращает список всех пользователей
func GetUsers(w http.ResponseWriter, r *http.Request) {
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
		if err := users.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.RegistrationDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userList = append(userList, user)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, userList)
}

// GetUser возвращает пользователя по ID
func GetUser(w http.ResponseWriter, r *http.Request) {
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
	err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.RegistrationDate)
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при подключении к базе данных"})
		return
	}
	defer db.Close()

	session, _ := store.Get(r, "CreateUserSession")

	// Получаем парраметры из запроса
	params := mux.Vars(r)
	UserName := params["username"]
	Email := params["email"]
	Password := params["password"]

	session.Values["username"] = UserName
	session.Values["email"] = Email
	session.Values["password"] = Password
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Вставка нового пользователя в базу данных
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3);", UserName, Email, Password)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении пользователя"})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusCreated, map[string]string{"Код": "200"})
}


func AuthUser(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	status := false

	// Получаем параметры из запроса
	params := mux.Vars(r)
	email := params["email"]
	password := params["password"]

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM users WHERE email = $1 AND password = $2;", email, password)

	var user User
	/* Сканирование копирует столбцы из сопоставленной строки в значения,
	на которые указывает dest. Смотрите документацию по строкам.Сканирование для
	получения подробной информации. Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.RegistrationDate)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusOK, status)
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusOK, status)
		return
	}
	status = true

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, status)
}


// UpdateUser обновляет данные пользователя по ID
func PutUser(w http.ResponseWriter, r *http.Request) {
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
        respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении пользователя"})
        return
    }

    // Передаем в функцию преобразования в json
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Данные пользователя обновлены"})
}


// DeleteUser удаляет пользователя по ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Инициализируем подключение к базе данных
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем идентификатор продукта из параметров запроса
	params := mux.Vars(r)
	userID := params["id"]

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
