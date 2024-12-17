package main

import (
	"fmt"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
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

// Context эта структура требуется для пакета gocraft/work
type Context struct{}

func main() {
	// создаем пул воркеров, который может исполнять несколько задач одновременно
	pool := work.NewWorkerPool(
		Context{},
		10,         // количество одновременно выполняемых задач
		"demo_app", // пространство имен очереди
		redisPool,  // пул задач redis
	)

	// создаем маппинг имен задачь к соответствующим им функциям выполнения
	pool.Job("email", SendEmail)

	// стартуем пул задач
	pool.Start()

	// ожидаем сигнала выхода (graceful stop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	// останавливаем пул
	pool.Stop()
}

func SendEmail(job *work.Job) error {
	addr := job.ArgString("email")
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Sending mail to ", addr, "with subject", subject)
	time.Sleep(2 * time.Second)
	return nil
}
