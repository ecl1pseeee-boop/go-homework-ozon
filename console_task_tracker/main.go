package main

import (
	"fmt"
	"os"
	"task_tracker/internal/cli"
	"task_tracker/internal/service"
	"task_tracker/internal/storage"
)

func main() {
	// s := storage.NewFileStorage("data.json")
	s := storage.NewMemoryStorage()
	taskService := service.NewTaskService(s)
	app := cli.NewApp(taskService)
	err := app.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
