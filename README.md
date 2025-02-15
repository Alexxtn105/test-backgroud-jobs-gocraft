# Тестовый проект по background jobs с использованием [gocraft/work](https://github.com/gocraft/work)
Особенности фреймворка:
- фоновые задачи используются в случае длительных операций;
- фоновые задачи отслеживают изменение в БД, в случае сбоя или перезапуска приложения есть возможность восстановить процесс обработки;
- контроль использования ресурсов путем изменения количества выполняемых фоновых задач;
- простота управления задачами.

## Установка
Установка фреймворка:
```bash
go get github.com/gocraft/work
```
Установка зависимостей пользовательского интерфейса:
```bash
go get github.com/gocraft/work/cmd/workwebui
```
Установка исполняемых файлов пользовательского интерфейса (папка `$GOPATH/bin`)
```bash
go install github.com/gocraft/work/cmd/workwebui
```


## Проверка работоспособности
Запуск пула воркеров
```bash
go run worker.go
```

Запуск фоновых задач
```bash
go run enqueuer/enqueuer.go
```

Запуск визуального интерфейса для просмотра текущего состояния задач:
```bash
workwebui -redis="redis://localhost:6379" -ns="demo_app" -listen=":5040"
```
Просмотр осуществляется в браузере по адресу http://localhost:5040.

Значения параметров команды:
```bash
-redis="redis://localhost:6379"     # адрес redis
-ns="demo_app"                      # пространство имен
-listen=":5040"                     # адрес сервера UI 
```
