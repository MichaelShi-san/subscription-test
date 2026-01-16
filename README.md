# Subscription Service


Простой REST-сервис для управления онлайн-подписками.


## 🔧 Требования


- Docker & Docker Compose
- Git
- (Опционально) Go 1.22+ если хотите собирать локально


---

## ⚙️ Настройка

1. Клонируем репозиторий:

```bash
git clone https://github.com/MichaelShi-san/subscription-test.git
cd subscription-test

2. Создаём .env на основе примера:

cp .env.example .env


---

Запуск сервиса через Docker:

docker compose up --build

---

REST API:

1. Создать подписку:

POST /subscriptions
Content-Type: application/json

{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "2025-07-01",
  "end_date": "2025-12-31" // опционально
}


2. Получить подписку по ID:

GET /subscriptions/{id}


3. Обновить подписку:

PUT /subscriptions/{id}
Content-Type: application/json
{
  "service_name": "Netflix",
  "price": 500,
  "start_date": "2025-07-01",
  "end_date": "2025-12-31"
}


4. Удалить подписку:

DELETE /subscriptions/{id}


5. Список подписок с фильтром:

GET /subscriptions?user_id=<uuid>&service_name=<название>


6. Суммарная стоимость подписок за период:

GET /subscriptions/total?user_id=<uuid>&service_name=<название>&start=YYYY-MM&end=YYYY-MM


---

Документация доступна в docs/swagger.yaml.


