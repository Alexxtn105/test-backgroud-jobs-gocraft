package main

import (
	"fmt"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// redisPool Пул очередей фоновых задач с использованием redis. Используется для того, чтобы брать задачи из очереди
var redisPool = &redis.Pool{
	MaxActive: 5, // максимум активных задач
	MaxIdle:   5, // максимум простаивающих задач
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

type User struct {
	ID    int64
	Email string
	Name  string
}

// Context эта структура требуется для пакета gocraft/work
type Context struct {
	currentUser *User
}

// FindCurrentUser Middleware для получения данных о пользователе
func (c *Context) FindCurrentUser(job *work.Job, next work.NextMiddlewareFunc) error {
	// если аргумент задачи содержит ИД пользователя
	if _, ok := job.Args["userID"]; ok {
		userID := job.ArgInt64("userID")

		// как будто берем данные пользователя из БД (упрощенно - генерим)
		c.currentUser = &User{
			ID:    userID,
			Email: "test" + strconv.Itoa(int(userID)) + "@mail.ru",
			Name:  "User" + strconv.Itoa(int(userID)),
		}
	}
	return next()
}

// Log Middleware для логирования начала выполнения задачи
func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Старт новой задачи:", job.Name, " ИД:", job.ID)
	return next()
}

func main() {
	// создаем пул воркеров, который может исполнять несколько задач одновременно
	pool := work.NewWorkerPool(
		Context{},
		10,         // количество одновременно выполняемых задач
		"demo_app", // пространство имен очереди
		redisPool,  // пул задач redis
	)

	// Добавим middleware в пул воркеров. В данном случае это лог начала задачи
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCurrentUser)

	// создаем маппинг имен задач к соответствующим им функциям выполнения
	// вариант БЕЗ ОПЦИЙ
	//pool.Job("email",  SendEmail)

	// Вариант, если нужно установить приоритет задачи, а также настроить перезапуск в случае сбоя - указать опции work.JobOptions:
	pool.JobWithOptions(
		"email", // имя задачи
		work.JobOptions{
			Priority: 10, // приоритет
			MaxFails: 1,  // максимальное количество повторных выполнений задачи в случае сбоя
		},
		(*Context).SendEmail, // выполняемая функция
	)

	// стартуем пул задач
	pool.Start()

	// ожидаем сигнала выхода (graceful stop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	// останавливаем пул
	pool.Stop()
}

func (c *Context) SendEmail(job *work.Job) error {
	// жестко указываем аргумент
	//addr := job.ArgString("email")
	// ИЛИ берем данные пользователя из контекста
	addr := c.currentUser.Email
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Sending mail to", addr, "with subject", subject)
	time.Sleep(2 * time.Second)
	return nil
}
