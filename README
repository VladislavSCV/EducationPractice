# API для взаимодействия с базой данных

## Описание проекта

Этот проект представляет собой API, разработанный для взаимодействия с базой данных (СУБД: PostgreSQL). API обеспечивает управление информацией о продуктах и пользователях.

# Начало работы

## Таблица Users
    ```sql
    CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
    ```

## Таблица Products

```sql
CREATE TABLE Products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255),
    price DECIMAL NOT NULL,
    status VARCHAR(255) DEFAULT 'В продаже',
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
    ```

2. После создания таблиц, вы можете взаимодействовать с API для управления данными в базе данных.

## Запуск проекта

Для запуска API выполните следующие шаги:

1. Перейдите в каталог с файлом `start.bat`.
2. Введите команду: `./start.bat`

Теперь вы можете начать работать с API и взаимодействовать с базой данных.

## Взаимодействие с API

API предоставляет следующие возможности:

- **Управление продуктами:**
  - Получение списка продуктов
  - Добавление нового продукта
  - Обновление информации о продукте
  - Удаление продукта

- **Управление пользователями:**
  - Получение списка пользователей
  - Добавление нового пользователя
  - Обновление информации о пользователе
  - Удаление пользователя
