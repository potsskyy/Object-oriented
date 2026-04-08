# Todo API on Go

Учебный backend-проект на Go для управления списком задач.

Приложение реализует:
- регистрацию пользователя;
- вход и выход из системы;
- серверные сессии через cookie;
- создание, просмотр, обновление, завершение и архивацию задач;
- автоматические тесты для основных сценариев работы.

## Назначение проекта

Проект демонстрирует базовую архитектуру REST API на Go:
- маршрутизацию через `gorilla/mux`;
- работу с HTTP-обработчиками и middleware;
- cookie-аутентификацию;
- хранение данных в памяти процесса;
- покрытие ключевой логики автоматическими тестами.

## Технологии

- Go
- `net/http`
- [`github.com/gorilla/mux`](https://github.com/gorilla/mux)
- `golang.org/x/crypto/bcrypt`

## Особенности реализации

- Пользователи и задачи хранятся в памяти, без базы данных.
- Пароли не сохраняются в открытом виде, а хешируются через `bcrypt`.
- После логина в cookie сохраняется не `username`, а случайный session ID.
- Проверка доступа идет по серверной сессии, которая хранится в памяти приложения.
- Cookie помечается как `HttpOnly` и `SameSite=Strict`.
- Сессия имеет ограниченный срок жизни.

Важно:

- проект является API, а не веб-сайтом;
- маршрут `/` не реализован, поэтому `http://localhost:8080/` возвращает `404 page not found`;
- при перезапуске приложения все данные и активные сессии очищаются.

## Структура проекта

```text
cmd/
  main.go                  точка входа и настройка маршрутов

internal/
  handler/                 HTTP-обработчики
  middleware/              проверка авторизации
  models/                  модели пользователя и задачи
  repository/              интерфейсы и in-memory хранилище
  session/                 настройки cookie и сессий
```

## Требования

- Go 1.24 или выше

Проверка версии:

```powershell
go version
```

## Установка зависимостей

Если Go установлен, зависимости можно скачать отдельно:

```powershell
go mod download
```

Либо просто запустить проект, и Go подтянет их автоматически.

## Запуск проекта

Запуск выполняется из корня репозитория:

```powershell
go run ./cmd
```

После старта в консоли появится сообщение:

```text
Server is running on :8080
```

Сервер будет доступен по адресу:

```text
http://localhost:8080
```

## Автоматические тесты

Для запуска всех тестов из корня репозитория:

```powershell
go test ./...
```

Тесты покрывают:
- регистрацию и логин;
- создание защищенной серверной сессии;
- отказ в доступе при поддельной cookie;
- logout и инвалидирование сессии;
- основной lifecycle задачи через HTTP API;
- жизненный цикл данных в in-memory репозитории.

## API

### Публичные маршруты

#### `POST /register`

Регистрирует нового пользователя.

Пример тела запроса:

```json
{
  "username": "admin",
  "password": "1234"
}
```

#### `POST /login`

Выполняет вход пользователя и устанавливает cookie с session ID.

Пример тела запроса:

```json
{
  "username": "admin",
  "password": "1234"
}
```

#### `POST /logout`

Удаляет активную серверную сессию и очищает cookie.

### Защищенные маршруты

Для вызова маршрутов ниже требуется успешный вход и cookie-сессия.

#### `POST /add`

Создает новую задачу.

```json
{
  "headline": "test",
  "details": "first task"
}
```

#### `POST /update`

Обновляет задачу по `id`.

```json
{
  "id": 1,
  "headline": "updated task",
  "details": "new details"
}
```

#### `POST /resolve/{id}`

Помечает задачу выполненной и переносит ее в архив.

#### `POST /delete/{id}`

Архивирует задачу.

#### `GET /get`

Возвращает все активные задачи пользователя.

#### `GET /get/{id}`

Возвращает одну задачу по идентификатору.

#### `GET /archive`

Возвращает архив задач пользователя.

## Пример ручной проверки через PowerShell

### 1. Регистрация

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/register -ContentType "application/json" -Body '{"username":"admin","password":"1234"}'
```

### 2. Логин и сохранение сессии

```powershell
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
Invoke-RestMethod -Method POST -Uri http://localhost:8080/login -WebSession $session -ContentType "application/json" -Body '{"username":"admin","password":"1234"}'
```

### 3. Создание задачи

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/add -WebSession $session -ContentType "application/json" -Body '{"headline":"test","details":"first task"}'
```

### 4. Получение списка задач

```powershell
Invoke-RestMethod -Method GET -Uri http://localhost:8080/get -WebSession $session | ConvertTo-Json -Depth 5
```

### 5. Обновление задачи

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/update -WebSession $session -ContentType "application/json" -Body '{"id":1,"headline":"updated task","details":"new details"}'
```

### 6. Завершение задачи

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/resolve/1 -WebSession $session
```

### 7. Проверка архива

```powershell
Invoke-RestMethod -Method GET -Uri http://localhost:8080/archive -WebSession $session | ConvertTo-Json -Depth 5
```

### 8. Выход из системы

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/logout -WebSession $session
```

## Ограничения проекта

- Используется in-memory хранилище без постоянного сохранения.
- Нет базы данных и миграций.
- Нет frontend-интерфейса.
- Проект ориентирован на учебный REST API и демонстрацию базовых backend-подходов.

## Автор

Уваров Марк
