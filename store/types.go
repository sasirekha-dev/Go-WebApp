package store

import "sync"

type ToDoItem struct {
	Id     string `json:"id"`
	Item   string `json:"item"`
	Status string `json:"status"`
}
type addRequest struct {
	item         string
	status       string
	ResponseChan chan ToDoItem
}
type updateRequest struct {
	item         string
	status       string
	ResponseChan chan error
}
type DeleteRequest struct {
	item         string
	ResponseChan chan error
}
type listRequest struct {
	ResponseChan chan []Todolist
}
type BaseStore struct {
	ToDo       []ToDoItem
	generator  UuidGenerator
	addChan    chan addRequest
	updateChan chan updateRequest
	deleteChan chan DeleteRequest
	listChan   chan listRequest
	mu         sync.Mutex
	wg         sync.WaitGroup
}
