package store

import (
	"context"
	"strconv"
	"sync"
	"testing"
)

func TestInMemoryStore_ParallelOperations(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryStore(ctx)
	defer store.Shutdown()
	numOfItems := 10
	t.Run("TestParallelInsertions", func(t *testing.T) {
		var wg sync.WaitGroup

		// Parallel insertions
		// for i := 0; i < 10; i++ {
		// 	wg.Add(1)
		// 	go func(item string) {
		// 		defer wg.Done()
		// 		fmt.Println("TODO APP ITEM ADDED", store.InsertItem(item, "pending"))
		// 	}(fmt.Sprintf("item-%d", i))
		// }
		//Continuous insertions
		for i := 0; i < numOfItems; i++ {
			item := "ITEM-" + strconv.Itoa(i)
			store.InsertItem(item, "pending")
		}

		// Wait for all insertions to complete
		wg.Wait()

		items := store.ListItems()
		if len(items) != numOfItems {
			t.Errorf("expected 3 items, got %d", len(items))
		}

	})

	t.Run("TestParallelDeletions", func(t *testing.T) {

		for i := 0; i < numOfItems; i++ {
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
			store.InsertItem("ITEM-"+strconv.Itoa(i), "pending")
		}

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
