# API "Sport World"

## Описание проекта

Этот проект представляет собой API, разработанный для взаимодействия с базой данных (СУБД: PostgreSQL). API обеспечивает управление информацией о продуктах и пользователях.

# Начало работы

## Таблица Users

```sql
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(25) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Таблица Products

```sql
CREATE TABLE Products (
    product_id SERIAL PRIMARY KEY,
    seller_id INT REFERENCES Users(user_id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255),
    price DECIMAL NOT NULL,
    status_id INT REFERENCES ProductStatuses(status_id),
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES Users(user_id),
    deleted_by INT REFERENCES Users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Создание таблицы заказов

```sql
CREATE TABLE Orders (
    order_id SERIAL PRIMARY KEY,
    buyer_id INT REFERENCES Users(user_id),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES Users(user_id),
    deleted_by INT REFERENCES Users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

##Создание таблицы товаров в заказе

```sql
CREATE TABLE OrderItems (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES Orders(order_id),
    product_id INT REFERENCES Products(product_id),
    quantity INT NOT NULL,
    total_price DECIMAL NOT NULL,
    created_by INT REFERENCES Users(user_id),
    deleted_by INT REFERENCES Users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```
## Создание таблицы Обмены
```sql
CREATE TABLE Trades (
    trade_id SERIAL PRIMARY KEY,
    initiator_id INT REFERENCES Users(user_id) NOT NULL,
    receiver_id INT REFERENCES Users(user_id) NOT NULL,
    product_id INT REFERENCES Products(product_id) NOT NULL,
    status_id INT REFERENCES TradeStatuses(status_id),
    initiated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    rejected_at TIMESTAMP
);
```


2. После создания таблиц, вы можете взаимодействовать с API для управления данными в базе данных.

## Запуск проекта

Для запуска API выполните следующие шаги:

1. Перейдите в каталог с файлом `start.bat`.
2. Введите команду: `./start.bat`

Теперь вы можете начать работать с API и взаимодействовать с базой данных.

## Взаимодействие с API

Для взаимодействия с api вы можете использовать postman или API tester(в моем случае это расширение для браузера).
В случае если вам не интересно использовать api  в таком виде, вы можете попробовать телеграм бота который был написан на js. И уже там собственно можете опробовать api) 
