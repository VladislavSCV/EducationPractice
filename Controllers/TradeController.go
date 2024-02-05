package Controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Trade структура представляет собой модель данных для таблицы "Обмены"
type Trade struct {
	TradeID     int          `json:"trade_id"`
	InitiatorID int          `json:"initiator_id"`
	ReceiverID  int          `json:"receiver_id"`
	ProductID   int          `json:"product_id"`
	StatusID    int          `json:"status_id"`
	InitiatedAt time.Time    `json:"initiated_at"`
	AcceptedAt  sql.NullTime `json:"accepted_at"`
	RejectedAt  sql.NullTime `json:"rejected_at"`
}

// GetUsers возвращает список всех пользователей
func GetTrades(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	// Закрываем соединение с бд при закрытии функции
	defer db.Close()

	// Делаем запрос из бд и получаем данные в users
	trades, err := db.Query("SELECT * FROM trades;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer trades.Close()

	// Создаем список из структур User
	var tradeList []Trade
	for trades.Next() {
		var trade Trade
		if err := trades.Scan(&trade.TradeID, &trade.InitiatorID, &trade.ReceiverID,
			&trade.ProductID, &trade.StatusID, &trade.InitiatedAt,
			&trade.AcceptedAt, &trade.RejectedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tradeList = append(tradeList, trade)
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, tradeList)
}

// GetUser возвращает пользователя по ID
func GetTrade(w http.ResponseWriter, r *http.Request) {
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
	_, err = strconv.Atoi(ID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

	// Запрос к бд
	row := db.QueryRow("SELECT * FROM trades WHERE trade_id = $1;", ID)

	var trade Trade
	/* Сканирование копирует столбцы из сопоставленной строки в значения,
	на которые указывает dest. Смотрите документацию по строкам.Сканирование для
	получения подробной информации. Если более одной строки соответствует запросу,
	сканирование использует первую строку и отбрасывает остальные. Если ни одна строка не
	соответствует запросу, Scan возвращает ErrNoRows. */
	err = row.Scan(&trade.TradeID, &trade.InitiatorID, &trade.ReceiverID, &trade.ProductID,
		&trade.StatusID, &trade.InitiatedAt, &trade.AcceptedAt, &trade.RejectedAt)
	if err == sql.ErrNoRows {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Торг не найден"})
		return
	} else if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении торга"})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, trade)
}

func CreateTrade(w http.ResponseWriter, r *http.Request) {
	// Создаем подключение к бд и обрабатываем ошибки
	db, err := initDB()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при подключении к базе данных"})
		return
	}
	defer db.Close()

	// Получаем парраметры из запроса
	params := mux.Vars(r)
	session, _ := store.Get(r, "session-name")
	receiver_id := params["receiverID"]
	product_id := params["productID"]
	status_id := params["statusID"]

	log.Printf("sessions: %v", session)

	_, err = db.Exec("INSERT INTO trades (initiator_id, receiver_id, product_id, status_id, initiated_at) VALUES ($1, $2, $3, $4, $5);",
		session, receiver_id, product_id, status_id)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении записи о торговле"})
		return
	}
	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusCreated, map[string]string{"Код": "200"})
}

// UpdateUser обновляет данные пользователя по ID
func PutTrade(w http.ResponseWriter, r *http.Request) {
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
	_, err = db.Exec("UPDATE trades SET "+what+" = $1 WHERE trade_id = $2", new, ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении пользователя", "erro2": err.Error()})
		return
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Данные пользователя обновлены"})
}

// DeleteUser удаляет пользователя по ID
func DeleteTrade(w http.ResponseWriter, r *http.Request) {
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
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный ID пользователя"})
		return
	}

	// Затем удаляем пользователя
	_, err = db.Exec("DELETE FROM trades WHERE trades_id = $1", ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении пользователя"})
		return
	}

	// Передаем в функцию преобразования в json
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Пользователь успешно удален"})
}
