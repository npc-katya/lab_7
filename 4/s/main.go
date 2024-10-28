package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// структура для приема JSON данных
type Data struct {
	Message string `json:"message"`
}

const (
	port    = ":8080"
	timeout = 5 * time.Second
)

func main() {

	/*
			Создание HTTP-сервера:
		•	Реализуйте базовый HTTP-сервер с обработкой простейших GET и POST запросов.
		•	Сервер должен поддерживать два пути:
			•	GET /hello — возвращает приветственное сообщение.
			•	POST /data — принимает данные в формате JSON и выводит их содержимое в консоль.
	*/

	// создание контекста для обработки сигналов завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// запуск сервера
	if err := runServer(ctx); err != nil {
		log.Fatal(err)
	}
}

// функция запуска сервера
func runServer(ctx context.Context) error {
	mux := http.NewServeMux() // Создаем мультиплексор для маршрутизации запросов

	// регистрация обработчиков для маршрутов
	mux.HandleFunc("/hello", handleHello) // обработка GET /hello
	mux.HandleFunc("/data", handleData)   // обработка POST /data

	srv := &http.Server{
		Addr:    "127.0.0.1" + port, // адрес сервера
		Handler: mux,
	}

	// запуск сервера
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	log.Printf("сервер запущен на %s", srv.Addr)
	<-ctx.Done() // ожидание сигнала завершения

	log.Println("начало корректного завершения сервера")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout) // таймаут для завершения
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("ошибка завершения: %w", err)
	}

	log.Println("сервер завершен")
	return nil
}

// обработчик для GET /hello
func handleHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // установить статус ответа 200 OK
	w.Write([]byte("привет!!!")) // отправка приветственного сообщения
}

// обработчик для POST /data
func handleData(w http.ResponseWriter, r *http.Request) {
	var data Data

	// декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	log.Printf("полученные данные: %s", data.Message) // выводим полученные данные в консоль
	w.WriteHeader(http.StatusOK)                      // установка статуса ответа 200 OK
}
