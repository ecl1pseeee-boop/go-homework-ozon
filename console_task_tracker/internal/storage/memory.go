package storage

import (
	"fmt"
	"slices"
	"task_tracker/internal/task"
)

type MemoryStorage struct {
	nextID task.ID
	tasks  []task.Task
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{tasks: make([]task.Task, 0), nextID: 1}
}

func (s *MemoryStorage) AddTask(newTask task.Task) (task.Task, error) {
	newTask.ID = s.nextID
	s.tasks = append(s.tasks, newTask)
	s.nextID++
	return newTask, nil
}

func (s *MemoryStorage) RemoveTask(taskID task.ID) error {
	for i, t := range s.tasks {
		if t.ID == taskID {
			s.tasks = slices.Delete(s.tasks, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("задача %d: %w", taskID, task.ErrTaskNotFound)
}

func (s *MemoryStorage) GetTask(taskID task.ID) (task.Task, error) {
	for _, t := range s.tasks {
		if t.ID == taskID {
			return t, nil
		}
	}

	return task.Task{}, fmt.Errorf("задача %d: %w", taskID, task.ErrTaskNotFound)
}

func (s *MemoryStorage) GetAllTasks() ([]task.Task, error) {
	tasks := slices.Clone(s.tasks)
	return tasks, nil
}

func (s *MemoryStorage) UpdateTask(updatedTask task.Task) error {
	for i, t := range s.tasks {
		if t.ID == updatedTask.ID {
			s.tasks[i] = updatedTask
			return nil
		}
	}
	return fmt.Errorf("задача %d: %w", updatedTask.ID, task.ErrTaskNotFound)
}
