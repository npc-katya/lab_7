package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// структура для json данных
type Data struct {
	Message string `json:"message"`
}

const port = ":8080"

func main() {
	// GET запрос
	resp, err := http.Get("http://localhost" + port + "/hello")
	if err != nil {
		log.Fatalf("ошибка при GET запросе: %v", err)
	}
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	body.ReadFrom(resp.Body)
	fmt.Printf("GET /hello:\nстатус: %s\nответ: %s\n", resp.Status, body.String())

	// POST запрос
	data := Data{Message: "привет от клиента!"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("ошибка при сериализации JSON: %v", err)
	}

	resp, err = http.Post("http://localhost"+port+"/data", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("ошибка при POST запросе: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("POST /data:\nстатус: %s\n", resp.Status)
}
