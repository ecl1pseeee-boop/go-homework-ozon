package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"task_tracker/internal/task"
)

type FileStorage struct {
	filePath string
}

type fileData struct {
	NextID task.ID     `json:"next_id"`
	Tasks  []task.Task `json:"tasks"`
}

func NewFileStorage(filePath string) *FileStorage {
	return &FileStorage{filePath: filePath}
}

func (f *FileStorage) open() (fileData, error) {
	bytes, err := os.ReadFile(f.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fileData{NextID: 1, Tasks: make([]task.Task, 0)}, nil
		}
		return fileData{}, fmt.Errorf("чтение файла: %w", err)
	}
	data := fileData{}
	parseErr := json.Unmarshal(bytes, &data)
	if parseErr != nil {
		return fileData{}, fmt.Errorf("парсинг из JSON: %w", parseErr)
	}
	return data, nil
}

func (f *FileStorage) save(data fileData) error {
	bytes, jsonErr := json.MarshalIndent(data, "", "  ")
	if jsonErr != nil {
		return fmt.Errorf("конвертация в JSON: %w", jsonErr)
	}
	err := os.WriteFile(f.filePath, bytes, 0644)
	if err != nil {
		return fmt.Errorf("запись в файл: %w", err)
	}
	return nil
}

func (f *FileStorage) AddTask(newTask task.Task) (task.Task, error) {
	storageData, connectionErr := f.open()
	if connectionErr != nil {
		return task.Task{}, fmt.Errorf("создание задачи: %w", connectionErr)
	}
	newTask.ID = storageData.NextID
	storageData.NextID++
	storageData.Tasks = append(storageData.Tasks, newTask)
	saveErr := f.save(storageData)
	if saveErr != nil {
		return task.Task{}, fmt.Errorf("сохранение данных: %w", saveErr)
	}
	return newTask, nil
}

func (f *FileStorage) RemoveTask(taskID task.ID) error {
	storageData, connectionErr := f.open()
	if connectionErr != nil {
		return fmt.Errorf("удаление задачи %d: %w", taskID, connectionErr)
	}

	index := -1

	for i, dbTask := range storageData.Tasks {
		if dbTask.ID == taskID {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("задача %d: %w", taskID, task.ErrTaskNotFound)
	}
	storageData.Tasks = append(storageData.Tasks[:index], storageData.Tasks[index+1:]...)
	saveErr := f.save(storageData)
	if saveErr != nil {
		return fmt.Errorf("сохранение данных: %w", saveErr)
	}
	return nil
}

func (f *FileStorage) GetTask(taskID task.ID) (task.Task, error) {
	storageData, connectionErr := f.open()
	if connectionErr != nil {
		return task.Task{}, fmt.Errorf("получение задачи %d: %w", taskID, connectionErr)
	}

	for _, dbTask := range storageData.Tasks {
		if dbTask.ID == taskID {
			return dbTask, nil
		}
	}

	return task.Task{}, fmt.Errorf("задача %d: %w", taskID, task.ErrTaskNotFound)
}

func (f *FileStorage) GetAllTasks() ([]task.Task, error) {
	storageData, connectionErr := f.open()
	if connectionErr != nil {
		return nil, fmt.Errorf("получение задач: %w", connectionErr)
	}

	return storageData.Tasks, nil
}

func (f *FileStorage) UpdateTask(updatedTask task.Task) error {
	storageData, connectionErr := f.open()
	if connectionErr != nil {
		return fmt.Errorf("обновление задачи %d: %w", updatedTask.ID, connectionErr)
	}

	for i, dbTask := range storageData.Tasks {
		if dbTask.ID == updatedTask.ID {
			storageData.Tasks[i] = updatedTask

			saveErr := f.save(storageData)
			if saveErr != nil {
				return fmt.Errorf("сохранение данных: %w", saveErr)
			}
			return nil
		}
	}

	return fmt.Errorf("задача %d: %w", updatedTask.ID, task.ErrTaskNotFound)
}
