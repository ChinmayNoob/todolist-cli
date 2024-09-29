package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

type Todo struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}

var todoFile = "todos.csv"

func main() {
	var rootCmd = &cobra.Command{
		Use: "todos",
	}

	var addCmd = &cobra.Command{
		Use:   "add <description>",
		Short: "Add a new task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			addTodo(args[0])
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List todos",
		Run: func(cmd *cobra.Command, args []string) {
			showAll, _ := cmd.Flags().GetBool("all")
			listTodos(showAll)
		},
	}

	listCmd.Flags().BoolP("all", "a", false, "Show all todos")

	var completeCmd = &cobra.Command{
		Use:   "complete <todoid>",
		Short: "Complete a todo",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			completeTodo(args[0])
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete <todoid>",
		Short: "Delete a todo",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			deleteTodo(args[0])
		},
	}

	rootCmd.AddCommand(addCmd, listCmd, completeCmd, deleteCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func loadFile() (*os.File, error) {
	f, err := os.OpenFile(todoFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %v", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}
	return f, nil
}

func closeFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}

func addTodo(description string) {
	f, err := loadFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer closeFile(f)

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()
	var newID int
	if len(records) > 0 {
		lastRecord := records[len(records)-1]
		newID, _ = strconv.Atoi(lastRecord[0])
		newID++
	} else {
		newID = 1
	}

	todo := Todo{
		ID:          newID,
		Description: description,
		CreatedAt:   time.Now(),
		IsComplete:  false,
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()
	if err := writer.Write([]string{
		strconv.Itoa(todo.ID),
		todo.Description,
		todo.CreatedAt.Format(time.RFC3339),
		strconv.FormatBool(todo.IsComplete),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing todo: %v\n", err)
		return
	}
	fmt.Printf("Added task: %s\n", todo.Description)
}

func listTodos(showAll bool) {
	// Implementation for listing todos
}

func completeTodo(todoID string) {
	// Implementation for completing a todo
}

func deleteTodo(todoID string) {
	// Implementation for deleting a todo
}