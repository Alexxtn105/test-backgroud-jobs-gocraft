package main

import (
	"fmt"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"log"
)

// redisPool Пул очередей фоновых задач с использованием redis
var redisPool = &redis.Pool{
	MaxActive: 5, // максимум активных задач
	MaxIdle:   5, // максимум простаивающих задач
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

// enqueuer постановщик задач в очередь (пространство имен - demo_app)
var enqueuer = work.NewEnqueuer("demo_app", redisPool)

func main() {
	_, err := enqueuer.Enqueue(
		"email", // имя задачи
		//work.Q{"email": "test@mail.ru", "subject": "Testing!"}, //аргументы задачи
		work.Q{"userID": 10, "subject": "Testing!"}, //аргументы задачи - заменил адрес на ИД пользователя - для использования в контексте
	)
	fmt.Println("Задача помещена в очередь. Проверьте вывод пула воркеров (worker.go)")
	if err != nil {
		log.Fatal(err)
	}
}
