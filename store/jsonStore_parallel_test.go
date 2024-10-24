package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
)

func Cleanup() error {
	err := os.Remove("Sample.json")
	if err != nil {
		return errors.New("not able to remove File...try again")
	} else {
		return nil
	}
}

func TestJSONMemoryStore_ParallelOperations(t *testing.T) {
	ctx := context.Background()
	Cleanup()
	store := NewJsonMemoryStore(ctx)
	// defer store.Shutdown()

	t.Run("TestParallelInsertions", func(t *testing.T) {
		var wg sync.WaitGroup
		// ctx := context.Background()
		// store := NewInMemoryStore(ctx)
		// defer store.Shutdown()

		//Parallel insertions
		// for i := 0; i < 10; i++ {
		// 	wg.Add(1)
		// 	go func(item string) {
		// 		defer wg.Done()
		// 		fmt.Println("TODO APP ITEM ADDED", store.InsertItem(item, "pending"))
		// 	}(fmt.Sprintf("item-%d", i))
		// }

		//Continuous insertions
		for i := 0; i < 10; i++ {
			item := "ITEM-" + strconv.Itoa(i)
			store.InsertItem(item, "pending")
		}

		// Wait for all insertions to complete
		wg.Wait()

		items := store.ListItems()
		fmt.Printf("%+v\n", items)
		if len(items) != 10 {
			t.Errorf("expected 10 items, got %d", len(items))
		}

	})

	t.Run("TestParallelDeletions", func(t *testing.T) {

		for i := 0; i < 10; i++ {
			itemName := "ITEM-" + strconv.Itoa(i)
			err := store.DeleteItem(itemName)
			if err != nil && err.Error() != "item was not deleted" {
				t.Errorf("unexpected error during deletion: %v", err)
			}

		}

		items := store.ListItems()
		if len(items) != 0 {
			t.Errorf("expected 0 items after deletion, got %d", len(items))
		}
	})

	t.Run("TestParallelUpdates", func(t *testing.T) {

		var wg sync.WaitGroup

		// Re-insert items for updating
		for i := 0; i < 5; i++ {
			store.InsertItem("jsonItem-"+strconv.Itoa(i), "pending")
		}
		fmt.Println("STORE", store.ToDo)

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(i int) {
				itemID := store.ToDo[i].Item
				defer wg.Done()
				err := store.UpdateItem(itemID, "complete")
				if err != nil {
					t.Errorf("error updating item: %v", err)
				}
			}(i)
		}
		wg.Wait()

		items := store.ListItems()
		for _, item := range items {
			if item.Status != "complete" {
				t.Errorf("expected item status to be 'completed', got %s", item.Status)
			}
		}
	})
}
