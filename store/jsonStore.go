package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type JsonStore struct {
	BaseStore
	FileName string
}

func NewJsonMemoryStore(ctx context.Context) *JsonStore {
	jsonStore := &JsonStore{
		BaseStore: BaseStore{
			ToDo:       make([]ToDoItem, 0),
			addChan:    make(chan addRequest, 1),
			deleteChan: make(chan DeleteRequest, 1),
			updateChan: make(chan updateRequest, 1),
			listChan:   make(chan listRequest, 1),
		},
		FileName: "Sample.json",
	}
	jsonStore.SetUUIDGenerator(&RealUUIDGenerator{}) // Set the default UUID generator
	jsonStore.CreateFile()
	jsonStore.LoadFile()
	go jsonStore.Run(ctx)
	return jsonStore
}

func (jsonStore *JsonStore) SetUUIDGenerator(generator UuidGenerator) {
	jsonStore.generator = generator
}

func (jsonStore *JsonStore) LoadFile() error {

	fmt.Println(jsonStore.FileName)
	content, ef := os.ReadFile(jsonStore.FileName)
	if ef != nil {
		log.Fatal("Read file error", ef)
		return nil
	}
	if len(content) == 0 {
		return nil
	}

	e := json.Unmarshal(content, &jsonStore.ToDo)
	if e != nil {
		log.Fatal("Unmarshal", ef)
		return nil
	}
	return nil
}
func (JsonStore *JsonStore) CreateFile() error {
	_, e := os.Stat(JsonStore.FileName)
	if e != nil {
		if os.IsNotExist(e) {
			_, err := os.Create(JsonStore.FileName)
			if err != nil {
				return errors.New("not able to create File...try again")
			} else {
				return nil
			}
		}
	}
	return nil
}

func (jsonStore *JsonStore) WriteToJsonFile() error {

	jObject, e := json.MarshalIndent(&jsonStore.ToDo, "", "   ")
	if e != nil {
		log.Fatal("Marshal error", e)
		return e
	}
	we := os.WriteFile(jsonStore.FileName, jObject, 0644)
	if we != nil {
		log.Fatal("Marshal error", we)
		return we
	}
	return nil
}
func (jsonStore *JsonStore) insertItem(req addRequest) {

	newToDo := ToDoItem{}
	jsonStore.mu.Lock()
	defer jsonStore.mu.Unlock()

	id := jsonStore.generator.GenerateUUID()
	newToDo.Id = id
	newToDo.Item = req.item
	newToDo.Status = req.status
	jsonStore.ToDo = append(jsonStore.ToDo, newToDo)
	err := jsonStore.WriteToJsonFile()
	if err != nil {
		req.ResponseChan <- newToDo
	}
	req.ResponseChan <- newToDo
}

func (jsonStore *JsonStore) InsertItem(item string, status string) ToDoItem {
	responseChan := make(chan ToDoItem, 1)
	jsonStore.addChan <- addRequest{item: item, status: status, ResponseChan: responseChan}
	return <-responseChan

}
func (jsonStore *JsonStore) updateItem(req updateRequest) {
	jsonStore.mu.Lock()
	defer jsonStore.mu.Unlock()

	for i, content := range jsonStore.ToDo {
		if content.Item == req.item {
			content.Status = req.status
			jsonStore.ToDo[i] = content
			req.ResponseChan <- nil
		}
	}
	err := jsonStore.WriteToJsonFile()
	if err != nil {
		req.ResponseChan <- err
	}
}

func (jsonStore *JsonStore) UpdateItem(item string, status string) error {
	responseChan := make(chan error, 1)
	jsonStore.updateChan <- updateRequest{item: item, status: status, ResponseChan: responseChan}
	return <-responseChan

}
func (jsonStore *JsonStore) deleteItem(req DeleteRequest) {
	jsonStore.mu.Lock()
	defer jsonStore.mu.Unlock()

	for i, content := range jsonStore.ToDo {
		if content.Item == req.item {
			jsonStore.ToDo = append(jsonStore.ToDo[:i], jsonStore.ToDo[(i+1):]...)
			req.ResponseChan <- nil
		}
	}
	err := jsonStore.WriteToJsonFile()
	if err != nil {
		req.ResponseChan <- err
	}
}

func (jsonStore *JsonStore) DeleteItem(deleteItem string) error {
	responseChan := make(chan error, 1)
	jsonStore.deleteChan <- DeleteRequest{item: deleteItem, ResponseChan: responseChan}
	return <-responseChan
}
func (jsonStore *JsonStore) listItems(req listRequest) {
	jsonStore.mu.Lock()
	defer jsonStore.mu.Unlock()
	ListOfItems := make([]Todolist, 0)
	err := jsonStore.LoadFile()
	if err != nil {
		req.ResponseChan <- ListOfItems
	}
	for _, item := range jsonStore.ToDo {
		ListOfItems = append(ListOfItems, Todolist{item.Item, item.Status})
	}
	req.ResponseChan <- ListOfItems
}

func (jsonStore *JsonStore) ListItems() []Todolist {
	responseChan := make(chan []Todolist, 1)
	jsonStore.listChan <- listRequest{ResponseChan: responseChan}
	return <-responseChan
}

func (jsonStore *JsonStore) Run(ctx context.Context) {
	for {
		select {
		case addReq, ok := <-jsonStore.addChan:
			if !ok {
				fmt.Print("closing add channel...")
				return
			}

			jsonStore.wg.Add(1)
			go func(addReq addRequest) {
				defer jsonStore.wg.Done()
				jsonStore.insertItem(addReq)
			}(addReq)
			jsonStore.wg.Wait()

		case updateReq, ok := <-jsonStore.updateChan:
			if !ok {
				fmt.Print("closing add channel...")
				return
			}

			jsonStore.wg.Add(1)
			go func(updateReq updateRequest) {
				defer jsonStore.wg.Done()
				jsonStore.updateItem(updateReq)
			}(updateReq)
			jsonStore.wg.Wait()

		case deleteReq, ok := <-jsonStore.deleteChan:
			if !ok {
				fmt.Print("closing add channel...")
				return
			}
			fmt.Println("DELETIONS....")
			jsonStore.wg.Add(1)
			go func(deleteReq DeleteRequest) {
				defer jsonStore.wg.Done()
				jsonStore.deleteItem(deleteReq)
			}(deleteReq)
			jsonStore.wg.Wait()

		case listReq, ok := <-jsonStore.listChan:
			if !ok {
				fmt.Print("closing add channel...")
				return
			}

			jsonStore.wg.Add(1)
			go func(listReq listRequest) {
				defer jsonStore.wg.Done()
				jsonStore.listItems(listReq)
			}(listReq)
			jsonStore.wg.Wait()

		case <-ctx.Done():
			jsonStore.wg.Wait()
		}
	}
}
