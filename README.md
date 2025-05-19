## Вимоги для запуску

*   Go (рекомендовано версію 1.18 або новішу)
*   MySQL сервер
*   Дійсний API ключ від [WeatherAPI.com](https://www.weatherapi.com/)

## Налаштування та запуск сервера локально

1.  **Клонуйте репозиторій:**
    ```
    git clone https://github.com/Kibarashka/weather_api.git
    cd weather_api
    ```

2.  **Встановіть залежності Go:**
    Перебуваючи в кореневій папці проекту, виконайте:
    ```
    go mod tidy
    ```

3.  **Налаштуйте MySQL базу даних:**
    *   Переконайтеся, що ваш MySQL сервер запущений.
    *   Підключіться до MySQL з правами адміністратора (наприклад, користувачем `root`).
    *   Створіть базу даних та користувача для додатку. Приклад SQL запитів:
        ```sql
        CREATE DATABASE IF NOT EXISTS weather_app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
        CREATE USER IF NOT EXISTS 'weather_user'@'localhost' IDENTIFIED BY 'your_actual_password'; -- тут потрібно буде поставити свій пароль
        GRANT ALL PRIVILEGES ON weather_app_db.* TO 'weather_user'@'localhost';
        FLUSH PRIVILEGES;
        ```

4.  **Налаштуйте змінні середовища:**
    *   Зробіть файл `.env` та заповніть його данними:
        ```
        # Database Configuration (MySQL)
        DB_HOST=localhost
        DB_USER=weather_user # ім'я користувача 
        DB_PASSWORD=your_actual_password # пароль
        DB_NAME=weather_app_db # назва бд
        DB_PORT=3306

        # Application Configuration
        APP_PORT=8080 # Порт, на якому буде працювати API
        APP_BASE_URL=http://localhost:8080 # Для генерації посилань в email

        # External Services API Keys
        WEATHER_API_KEY=your_actual_weatherapi_com_key # API ключ
        ```
       *Також важливо:* Файл `.env` містить секретні дані і вже доданий до `.gitignore`, тому він не потрапить у репозиторій.
         **Запустіть сервер:**
    Перебуваючи в кореневій папці проекту, виконайте:
    ```
    go run cmd/api/main.go 
    # Або, якщо ваш головний файл знаходиться в weather_api/api/main.go:
    # go run api/main.go
    ```
    Ви маєте побачити в консолі логи про успішний запуск сервера на вказаному порту (за замовчуванням `8080`).


## Основні API Ендпоінти

Базовий URL для всіх запитів: `http://localhost:PORT/api`

*   **Отримати поточну погоду:**
    *   `GET /weather?city={cityName}`
    *   Приклад: `GET http://localhost:8080/api/weather?city=Kyiv`
*   **Підписатися на оновлення:**
    *   `POST /subscribe`
    *   Тіло запиту (`application/json` або `application/x-www-form-urlencoded`):
        ```json
        {
            "email": "user@example.com",
            "city": "Lviv",
            "frequency": "daily" // "daily" або "hourly"
        }
        ```
*   **Підтвердити підписку:**
    *   `GET /confirm/{token}` (токен надсилається на email після запиту на підписку)
*   **Відписатися від оновлень:**
    *   `GET /unsubscribe/{token}` (токен для відписки надається після підтвердження або в листах з оновленнями)

Планувалося додати підтримку Docker для спрощення розгортання та забезпечення консистентного середовища. Однак, у процесі виникли певні технічні складнощі з налаштуванням Dockerfile та Docker Compose, які потребували додаткового часу на вирішення.
У поточній версії проект запускається локально без Docker, як описано в розділі "Налаштування та запуск сервера локально". Додавання повноцінної Docker-підтримки розглядається як один з наступних кроків у розвитку проекту.
