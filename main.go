package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoPayLoad struct {
	Desc string `json:"description"`
}

var todosStore = TodoStore{
	todos: []Todo{{1, "Read a book", false},
		{2, "Exercise", true}},
}

func main() {
	fmt.Print("Server is Running")
	r := mux.NewRouter()
	r.HandleFunc("/get", rootHandler).Methods("GET")
	r.HandleFunc("/get/todos", getTodosHandler).Methods("GET")
	r.HandleFunc("/get/todos/{id}", getTodoDetailHandler).Methods("GET")
	r.HandleFunc("/get/todos", addTodoHandler).Methods("POST")
	r.HandleFunc("/get/todos/{id}/mark_completed", markCompleteHandler).Methods("PUT")
	r.HandleFunc("/get/todos/{id}/mark_incompleted", markInCompleteHandler).Methods("PUT")
	r.HandleFunc("/get/todos/{id}", updateTodoHandler).Methods("PUT")
	r.HandleFunc("/get/todos/{id}", deleteTodoHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
func writeMessage(w http.ResponseWriter, status int, message string) {
	writeJson(w, status, map[string]string{
		"message": message,
	})
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, map[string]string{"message": "Hello World!"})
}
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("completed") == "true" {
		writeJson(w, http.StatusOK, todosStore.GetByStatus(true))
	} else if r.URL.Query().Get("completed") == "false" {
		writeJson(w, http.StatusOK, todosStore.GetByStatus(false))
	} else {
		writeJson(w, http.StatusOK, todosStore.GetAll())
	}

}
func getTodoDetailHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	todo, err := todosStore.GetTodoDetail(id)
	if err != nil {
		writeMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJson(w, http.StatusOK, todo)
}
func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	var payload TodoPayLoad
	json.NewDecoder(r.Body).Decode(&payload)
	todosStore.AddTodo(payload.Desc)
	writeMessage(w, http.StatusOK, "Successfully Added")
}
func markCompleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	err1 := todosStore.SetComplete(id, true)
	if err1 != nil {
		writeMessage(w, http.StatusNotFound, err1.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully Mark Completed")
}
func markInCompleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	err1 := todosStore.SetComplete(id, false)
	if err1 != nil {
		writeMessage(w, http.StatusNotFound, err1.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully Mark InCompleted")
}
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	var payload TodoPayLoad
	json.NewDecoder(r.Body).Decode(&payload)
	err1 := todosStore.Update(id, payload.Desc)
	if err1 != nil {
		writeMessage(w, http.StatusNotFound, err1.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully updated")
}
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid Id")
		return
	}
	err1 := todosStore.Delete(id)
	if err1 != nil {
		writeMessage(w, http.StatusNotFound, err1.Error())
		return
	}
	writeMessage(w, http.StatusOK, "successfully Deleted")
}
