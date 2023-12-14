package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func getAllTasks(res http.ResponseWriter, req *http.Request) {
	// сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента JSON
	res.Header().Set("Content-Type", "application/json")
	// статус OK
	res.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	res.Write(resp)
}
func postTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer
	// Читаем в буфер из тела запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Сериализуем данные из буфера
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Проверяем существование задачи
	if _, ok := tasks[task.ID]; ok {
		http.Error(res, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}
	// Сохраняем задачу в слайсе tasks
	tasks[task.ID] = task
	resp, err := json.Marshal(tasks[task.ID])
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}
func getTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	// Проверяем существование задачи
	task, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена", http.StatusNotFound)
		return
	}
	// Сериализуем данные
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(res, fmt.Sprintf("Ошибка сериализации: %v", err), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	// Записываем сериализованные в JSON данные в тело ответа
	res.Write(resp)
}
func deleteTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	// Проверяем существование задачи
	_, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена", http.StatusBadRequest)
		return
	}
	// Удаляем задачу из слайса tasks
	delete(tasks, id)
	res.WriteHeader(http.StatusOK)
}
func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
