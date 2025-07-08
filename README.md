## Демонстрационный сервис с простейшим интерфейсом для получения данных по некоторому uid-заказа

--- 

### Требования ⚙️
Для запуска сервиса необходимо в корневой папке проекта создать `.env` файл и заполнить по шаблону:
```
HTTP_USER=admin
HTTP_PASSWORD=admin777

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=wb-tech

KAFKA_ADDRESS=kafka:29092
TOPIC=wb-tech-topic
GROUP=wb-tech-consumer-group

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=
```
--- 

### Запуск 🔧
1. Выполните в терминале команду:
`docker-compose build`

2. Затем поднимите собранный dc командой:
`docker-compose up`

3. После сборки следует подождать некоторое время, пока `kafka` полностью заработает и запуститься сам сервер, затем можно приступить к запуску скрипта следующей командой:
`go run cmd/emulator-script/main.go`

После сгенерированных uid можно перейти по адресу `http://localhost:8081/order/{здесь вставляете ваш uid}`

--- 
### Видеоинтерфейс запуска
https://github.com/user-attachments/assets/b48d325c-b3ef-4610-b912-69f32821bda2
