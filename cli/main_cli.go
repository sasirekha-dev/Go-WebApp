package cli

import (
	"WebApp/store"
	"flag"
	"fmt"
	"os"
)

func main_cli(inMemory store.ToDoStore) {
	args := os.Args[1:]

	insert := flag.NewFlagSet("insert", flag.ExitOnError)
	item := insert.String("item", "Default Item", "Todo item to add")
	status := insert.String("status", "pending", "Status of the Todo item")

	delete := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteItem := delete.String("item", "", "item to delete")

	update := flag.NewFlagSet("update", flag.ExitOnError)
	updateItem := update.String("item", "", "Todo item to add")
	updateStatus := update.String("status", "", "Status of the Todo item")

	switch args[0] {
	case "insert":
		insert.Parse(args[1:])
		inMemory.InsertItem(*item, *status)
	case "update":
		update.Parse(args[1:])
		inMemory.UpdateItem(*updateItem)
		fmt.Printf("%s is updated to %s", *updateItem, *updateStatus)
	case "delete":
		delete.Parse(args[1:])
		inMemory.DeleteItem(*deleteItem)
		fmt.Printf("%s is deleted", *deleteItem)
	case "list":
		fmt.Print("list")
		inMemory.ListItems()
	default:
		fmt.Print("Not a valid option, Available options are insert, delete, update, list")
	}
}
