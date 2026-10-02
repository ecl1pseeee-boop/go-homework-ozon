package main

import (
	"fmt"
	"os"
	"task_tracker/internal/cli"
	"task_tracker/internal/service"
	"task_tracker/internal/storage"
)

func main() {
	fileStorage := storage.NewFileStorage("data.json")
	taskService := service.NewTaskService(fileStorage)
	app := cli.NewApp(taskService)
	err := app.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
