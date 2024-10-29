package store

import (
	"WebApp/store"
	"context"
	"strconv"
	"testing"
)

// BenchmarkInMemoryStore_Insertions benchmarks parallel insertions.
func BenchmarkInMemoryStore_Insertions(b *testing.B) {
	ctx := context.Background()
	store := store.NewInMemoryStore(ctx)
	defer store.Shutdown()

	b.ResetTimer() // Reset timer to exclude setup time
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Insert items concurrently
			item := strconv.Itoa(b.N)
			store.InsertItem("item-"+item, "pending")
		}
	})
}

// BenchmarkInMemoryStore_Deletions benchmarks parallel deletions.
func BenchmarkInMemoryStore_Deletions(b *testing.B) {
	ctx := context.Background()
	store := store.NewInMemoryStore(ctx)
	defer store.Shutdown()

	// Prepopulate the store with items
	for i := 0; i < b.N; i++ {
		store.InsertItem("item-"+strconv.Itoa(i), "pending")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Attempt to delete items concurrently
			item := "item-" + strconv.Itoa(b.N)
			store.DeleteItem(item)
		}
	})
}

// BenchmarkInMemoryStore_Updates benchmarks parallel updates.
func BenchmarkInMemoryStore_Updates(b *testing.B) {
	ctx := context.Background()
	store := store.NewInMemoryStore(ctx)
	defer store.Shutdown()

	// Prepopulate store with items
	for i := 0; i < b.N; i++ {
		store.InsertItem("item-"+strconv.Itoa(i), "pending")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Update items concurrently
			itemID := "item-" + strconv.Itoa(b.N)
			store.UpdateItem(itemID, "complete")
		}
	})
}
