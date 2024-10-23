package store

import (
	"context"
	"reflect"
	"slices"
	"testing"
)

type MockUUIDGenerator struct {
	MockedUUID string
}

func (g MockUUIDGenerator) generateUUID() string {
	return g.MockedUUID
}

func TestInsertInMem(t *testing.T) {
	t.Run("test insert", func(t *testing.T) {
		//Given
		ctx := context.Background()
		store := NewInMemoryStore(ctx)
		mockedUUID := "1"
		mockGenerator := MockUUIDGenerator{MockedUUID: mockedUUID}

		store.ToDo = []ToDoItem{}
		store.SetUUIDGenerator(mockGenerator)

		//when
		got := store.InsertItem("plant rose", "pending")
		// mockedUUID = "2"
		// mockGenerator2 := MockUUIDGenerator{MockedUUID: mockedUUID}
		// store.SetUUIDGenerator(mockGenerator2)
		// store.InsertItem("groceries", "pending")
		// fmt.Print("The Todo list is: ", ToDo)
		//Then
		want := ToDoItem{Id: "1", Item: "plant rose", Status: "pending"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", store.ToDo, want)
		}

	})
}

func TestDeleteInMem(t *testing.T) {
	t.Run("delete an item", func(t *testing.T) {
		store := InMemoryStore{}
		store.ToDo = []ToDoItem{{Id: "1", Item: "abc", Status: "pending"}}
		want := []ToDoItem{}
		store.DeleteItem("abc")
		if !slices.Equal(store.ToDo, want) {
			t.Errorf("got %v, want %v", store.ToDo, want)
		}
	})
}

func TestUpdateInMem(t *testing.T) {
	t.Run("update an item", func(t *testing.T) {
		store := InMemoryStore{}
		store.ToDo = []ToDoItem{{Id: "1", Item: "abc", Status: "pending"}}

		want := []ToDoItem{{Id: "1", Item: "abc", Status: "completed"}}

		store.UpdateItem("1")

		if !slices.Equal(store.ToDo, want) {
			t.Errorf("got %v, want %v", store.ToDo, want)
		}
	})
}

func TestListInMem(t *testing.T) {
	t.Run("List items", func(t *testing.T) {
		ctx := context.Background()
		store := NewInMemoryStore(ctx)
		defer store.Shutdown()
		// var buf bytes.Buffer

		store.ToDo = []ToDoItem{{Id: "1", Item: "a", Status: "pending"}, {Id: "2", Item: "b", Status: "done"}}
		got := store.ListItems()
		print(got)
		if got[0].Task != "a" && got[0].Status != "pending" {
			t.Errorf("got %q, want %q", got, Todolist{"a", "pending"})
		}
		// 		PrintToBuffer(&buf, inMem.ListItems())

		// 		got := buf.String()
		// 		want := `1.a: pending
		// 2.b: done
		// `
		// 		if !reflect.DeepEqual(got, want) {
		// 			t.Errorf("got %v, want %v", got, want)
		// 		}
	})
}

// // 	t.Run("Error when no items to list", func(t *testing.T) {
// // 		var buf bytes.Buffer
// // 		ToDo = []ToDoItem{}
// // 		PrintToBuffer(&buf, ListItems)

// // 		got := buf.String()
// // 		want := "No items to display"

// // 		if !reflect.DeepEqual(got, want) {
// // 			t.Errorf("got %q, want %q", got, want)
// // 		}
// // 	})
// // }

// // // Utility function to capture the output of a function
// // func PrintToBuffer(buf *bytes.Buffer, fn func()) {
// // 	// Save the original stdout
// // 	originalStdout := os.Stdout

// // 	// Set os.Stdout to the buffer to capture the output
// // 	r, w, _ := os.Pipe()
// // 	os.Stdout = w

// // 	// Execute the function
// // 	fn()

// // 	// Close writer and restore stdout
// // 	w.Close()
// // 	os.Stdout = originalStdout

// // 	// Read the captured output into the buffer
// // 	buf.ReadFrom(r)
// // }
