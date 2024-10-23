package store

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/google/uuid"
)

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
type InMemoryStore struct {
	BaseStore
}

func NewInMemoryStore(ctx context.Context) *InMemoryStore {
	store := &InMemoryStore{
		BaseStore: BaseStore{
			ToDo:       make([]ToDoItem, 0), // Initialize the ToDo slice
			addChan:    make(chan addRequest, 1),
			deleteChan: make(chan DeleteRequest, 1),
			updateChan: make(chan updateRequest, 1),
			listChan:   make(chan listRequest, 1),
		},
	}
	store.SetUUIDGenerator(&RealUUIDGenerator{}) // Set the default UUID generator
	go store.Run(ctx)
	return store
}

type ToDoStore interface {
	InsertItem(string, string) ToDoItem
	DeleteItem(string) error
	UpdateItem(string) error
	ListItems() []Todolist
}

type UuidGenerator interface {
	generateUUID() string
}

func (inMem *InMemoryStore) SetUUIDGenerator(generator UuidGenerator) {
	inMem.generator = generator
}

type RealUUIDGenerator struct{}

// NewUUID generates a real UUID
func (g RealUUIDGenerator) generateUUID() string {
	return uuid.New().String()
}

func (store *InMemoryStore) insertItem(req addRequest) {
	store.mu.Lock() // Lock before modifying the shared ToDo slice
	defer store.mu.Unlock()

	if req.status == "" {
		req.status = "pending"
	}
	todo := ToDoItem{Id: store.generator.generateUUID(), Item: req.item, Status: req.status}
	store.ToDo = append(store.ToDo, todo)
	fmt.Printf("Adding item %+v\n", store.ToDo)

	req.ResponseChan <- todo
}

func (store *InMemoryStore) InsertItem(item, status string) ToDoItem {
	responseChan := make(chan ToDoItem, 1)
	store.addChan <- addRequest{item: item, status: status, ResponseChan: responseChan}
	return <-responseChan
}

func (store *InMemoryStore) deleteItem(req DeleteRequest) {
	store.mu.Lock() // Lock before modifying the shared ToDo slice
	defer store.mu.Unlock()

	fmt.Println("Item to delete----->>>>>", req.item, store.ToDo)
	if len(req.item) == 0 {
		req.ResponseChan <- errors.New("no item to delete")
		return
	}
	for i, task := range store.ToDo {
		if task.Item == req.item {
			store.ToDo = append(store.ToDo[:i], store.ToDo[i+1:]...)
			req.ResponseChan <- nil
			return
		}
	}
	req.ResponseChan <- errors.New("item was not deleted")

}
func (store *InMemoryStore) DeleteItem(taskName string) error {
	// fmt.Print("received delete request")
	responseChan := make(chan error)
	store.deleteChan <- DeleteRequest{item: taskName, ResponseChan: responseChan}
	return <-responseChan
}

func (store *InMemoryStore) updateItem(req updateRequest) {
	store.mu.Lock() // Lock before modifying the shared ToDo slice
	defer store.mu.Unlock()

	for i, task := range store.ToDo {
		if task.Item == req.item {
			store.ToDo[i].Status = "completed"
			req.ResponseChan <- nil
		}
	}
	req.ResponseChan <- errors.New("no item found")
}

func (store *InMemoryStore) UpdateItem(task string) error {
	responseChan := make(chan error, 1)
	store.updateChan <- updateRequest{item: task, ResponseChan: responseChan}
	return <-responseChan

}

type Todolist struct {
	Task   string `json:"item"`
	Status string `json:"status"`
}

func (store *InMemoryStore) listItems(req listRequest) {
	store.mu.Lock() // Lock before modifying the shared ToDo slice
	defer store.mu.Unlock()
	var ListOfItems []Todolist
	if len(store.ToDo) == 0 {
		req.ResponseChan <- make([]Todolist, 0)
	}
	for _, item := range store.ToDo {
		// fmt.Printf("%d.%s: %s\n", i+1, item.Item, item.Status)
		ListOfItems = append(ListOfItems, Todolist{item.Item, item.Status})
	}
	req.ResponseChan <- ListOfItems
}

func (store *InMemoryStore) ListItems() []Todolist {
	listresponseChan := make(chan []Todolist, 1)
	store.listChan <- listRequest{ResponseChan: listresponseChan}
	list := <-listresponseChan
	fmt.Println("List of items=====>>>", list)
	return list
}

func (store *InMemoryStore) Shutdown() {
	store.wg.Wait()
}

// Listens to request and process them sequentially
func (store *InMemoryStore) Run(ctx context.Context) {
	fmt.Println("Go routine started....listening to request")
	for {
		select {
		case addReq, ok := <-store.addChan:
			fmt.Println("received req--> add")
			if !ok {
				fmt.Print("closing add channel...")
				return
			}
			store.wg.Add(1) // Increment the WaitGroup counter
			go func(req addRequest) {
				defer store.wg.Done()
				fmt.Println("Starting goroutine for adding item")
				store.insertItem(req)
				fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())
			}(addReq)
			store.wg.Wait()

		case updateReq, ok := <-store.updateChan:
			fmt.Println("received req--> update")
			if !ok {
				fmt.Print("closing update channel...")
				return
			}

			store.wg.Add(1) // Increment the WaitGroup counter
			go func(req updateRequest) {
				defer store.wg.Done()
				fmt.Println("Starting goroutine for updating item")
				store.updateItem(req)
				fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())
			}(updateReq)
			store.wg.Wait()
		case deleteReq, ok := <-store.deleteChan:
			fmt.Println("received req--> delete")
			if !ok {
				fmt.Print("closing delete channel...")
				return
			}

			store.wg.Add(1) // Increment the WaitGroup counter
			go func(req DeleteRequest) {
				defer store.wg.Done()
				fmt.Println("Starting goroutine for deleting item")
				store.deleteItem(req)

				fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())
			}(deleteReq)
			store.wg.Wait()
		case listReq, ok := <-store.listChan:
			fmt.Println("received req--> list")
			if !ok {
				fmt.Print("closing list channel...")
				return
			}

			store.wg.Add(1) // Increment the WaitGroup counter

			go func(req listRequest) {
				defer store.wg.Done()
				fmt.Println("Starting goroutine for list item")
				store.listItems(req)
				fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())
			}(listReq)
			store.wg.Wait()
		case <-ctx.Done():
			fmt.Println("received shutdown signal, waiting for all operations to finish")
			store.Shutdown()
			fmt.Println("workgroup ends.....")
			fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
			return
		}

	}
}