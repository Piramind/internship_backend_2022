# Инструкция по использованию
## База данных
Предполагается, что сервис подключается к базе данных PostgreSQL, которая имеет:
- Две таблицы: users и reserves. Для хранения актуальных балансов пользователей и зарезервированных транзауций соответственно
- Таблица users имеет поля id (varchar) и balance (int). Для хранения идентификатора и баланса пользователя соответственно. Предполагается, что id заранее известен тому, кто подаёт запрос. Будем считать, что баланс в неких виртуальных единицах, поэтому имеет целочисленное значение.
- Таблица reserves имеет поля id (varchar), service_id (varchar), order_id (varchar), balance (int).
## Запросы
Запросы можно разделить на 2 типа:
- Которые обращаются к таблице users. Они отправляют json соответствующий формату таблицы.
```
{
  id: "...",
  balance: ...
}
```
- Которые обращаются к таблице users. Они тоже отправляют json соответствующий формату таблицы.
```
{
  id: "...",
  service_id: "...",
  order_id: "...",
  balance: ...
}
```
Важно: balance в данном случае довольно противоречивое название, потому что в это поле также может записываться сумма пополнения (а не сам баланс)

Программа может принимать следующие запросы:
- Получение баланса пользователя по id. Возвращает json c id и балансом. Формат запроса:
```
curl host:port/balances/id
```
- Пополнение баланса пользователя:
```
curl host:port/transfer --include --header "Content-Type: application/json" -d @file.json --request "POST"
```
- Создание нового пользователя с балансом:
```
curl host:port/create_balance --include --header "Content-Type: application/json" -d @file.json --request "POST"
```
- Создание резерва баланса пользователя:
```
curl host:port/create_reserve --include --header "Content-Type: application/json" -d @file.json --request "POST"
```
- Пополнение резерва баланса пользователя:
```
curl host:port/reserve --include --header "Content-Type: application/json" -d @file.json --request "POST"
```
- Признание резерва:
```
curl host:port/admit --include --header "Content-Type: application/json" -d @file.json --request "POST"
```

Принимаемые данные для запросов как указано в задании.
## Как запустить
```
docker-compose up
```
