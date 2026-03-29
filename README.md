# Todo API on Go

Небольшой учебный REST API для работы со списком задач.

Проект позволяет:
- регистрировать пользователя;
- выполнять вход по логину и паролю;
- создавать задачи;
- получать список активных задач;
- получать задачу по `id`;
- обновлять задачу;
- помечать задачу выполненной;
- отправлять задачу в архив и просматривать архив.

## Стек

- Go
- `net/http`
- [gorilla/mux](https://github.com/gorilla/mux)
- in-memory хранилище без базы данных

## Особенности реализации

- Аутентификация сделана через cookie `user_session`.
- Данные хранятся в памяти процесса.
- После остановки приложения все пользователи и задачи удаляются.
- Сервер запускается на порту `8080`.

## Структура проекта

```text
cmd/
  main.go                  - точка входа, настройка роутов и запуск сервера
internal/
  handler/                 - HTTP-обработчики
  middleware/              - middleware для проверки авторизации
  models/                  - модели User и Task
  repository/              - интерфейсы и in-memory реализация хранилища
```

## Требования

- Go 1.24 или выше

Проверка установленной версии:

```powershell
go version
```

## Установка зависимостей

Если Go уже установлен, зависимости подтянутся автоматически при первом запуске.

При необходимости можно отдельно выполнить:

```powershell
go mod download
```

## Запуск проекта

Из корня проекта:

```powershell
cd C:\code\GO\Work
go run ./cmd
```

Ожидаемый результат:

```text
Server is running on :8080
```

Адрес сервера:

```text
http://localhost:8080
```

Важно:

- маршрут `/` не реализован;
- при открытии `http://localhost:8080/` в браузере сервер вернет `404 page not found`;


## API эндпоинты

### Публичные маршруты

#### `POST /register`

Регистрация нового пользователя.

Пример тела запроса:

```json
{
  "username": "admin",
  "password": "1234"
}
```

#### `POST /login`

Авторизация пользователя. После успешного входа сервер устанавливает cookie `user_session`.

Пример тела запроса:

```json
{
  "username": "admin",
  "password": "1234"
}
```

### Защищенные маршруты

Для всех маршрутов ниже требуется cookie после логина.

#### `POST /add`

Создать задачу.

Пример тела запроса:

```json
{
  "headline": "test",
  "details": "first task"
}
```

#### `POST /update`

Обновить существующую задачу.

Пример тела запроса:

```json
{
  "id": 1,
  "headline": "updated task",
  "details": "new details"
}
```

#### `POST /resolve/{id}`

Пометить задачу выполненной и отправить в архив.

Пример:

```text
POST /resolve/1
```

#### `POST /delete/{id}`

Архивировать задачу.

Пример:

```text
POST /delete/1
```

#### `GET /get`

Получить список активных задач пользователя.

#### `GET /get/{id}`

Получить одну задачу по идентификатору.

#### `GET /archive`

Получить архив задач пользователя.

## Полная проверка через PowerShell

### 1. Регистрация пользователя

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/register -ContentType "application/json" -Body '{"username":"admin","password":"1234"}'
```

### 2. Вход и сохранение cookie-сессии

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

### 5. Получение задачи по `id`

```powershell
Invoke-RestMethod -Method GET -Uri http://localhost:8080/get/1 -WebSession $session | ConvertTo-Json -Depth 5
```

### 6. Обновление задачи

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/update -WebSession $session -ContentType "application/json" -Body '{"id":1,"headline":"updated task","details":"new details"}'
```

### 7. Завершение задачи

```powershell
Invoke-RestMethod -Method POST -Uri http://localhost:8080/resolve/1 -WebSession $session
```

### 8. Получение архива

```powershell
Invoke-RestMethod -Method GET -Uri http://localhost:8080/archive -WebSession $session | ConvertTo-Json -Depth 5
```

## Пример ответа сервера для списка задач

```json
[
  {
    "id": 1,
    "user": "admin",
    "headline": "test",
    "details": "first task",
    "completed": false,
    "archived": false,
    "created_on": "2026-03-30T00:08:05.881251+03:00",
    "updated_on": "2026-03-30T00:08:05.881251+03:00"
  }
]
```

## Возможные ограничения проекта

- Данные не сохраняются в базу данных.
- Нет разделения на роли пользователей.
- Нет logout-маршрута.
- Нет фронтенда, проект представляет собой только backend API.
- Нет автоматических тестов.

## Автор

Уваров Марк


