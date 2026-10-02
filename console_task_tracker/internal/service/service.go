package service

import (
	"fmt"
	"slices"
	"strings"
	"task_tracker/internal/task"
	"time"
)

type Storage interface {
	AddTask(newTask task.Task) (task.Task, error)
	RemoveTask(taskID task.ID) error
	GetTask(taskID task.ID) (task.Task, error)
	GetAllTasks() ([]task.Task, error)
	UpdateTask(updatedTask task.Task) error
}

type TaskService struct {
	storage Storage
}

func NewTaskService(storage Storage) *TaskService {
	return &TaskService{storage: storage}
}

func (s *TaskService) CreateTask(title string, deadline time.Time) (task.Task, error) {
	if errTitle := s.validateTitle(title); errTitle != nil {
		return task.Task{}, errTitle
	}

	if errDeadline := s.validateDeadline(deadline); errDeadline != nil {
		return task.Task{}, errDeadline
	}

	newTask := task.Task{
		Title:    strings.TrimSpace(title),
		Deadline: deadline,
		Status:   task.Planned,
	}

	savedTask, err := s.storage.AddTask(newTask)
	if err != nil {
		return task.Task{}, fmt.Errorf("сохранение задачи: %w", err)
	}
	return savedTask, nil
}

func (s *TaskService) RemoveTask(taskID task.ID) error {
	err := s.storage.RemoveTask(taskID)
	if err != nil {
		return fmt.Errorf("удаление задачи: %w", err)
	}
	return nil
}

func (s *TaskService) GetTask(taskID task.ID) (task.Task, error) {
	t, err := s.storage.GetTask(taskID)
	if err != nil {
		return task.Task{}, fmt.Errorf("получение задачи: %w", err)
	}
	return t, nil
}

func (s *TaskService) GetAllTasks() ([]task.Task, error) {
	allTasks, err := s.storage.GetAllTasks()
	if err != nil {
		return nil, fmt.Errorf("получение задач: %w", err)
	}
	slices.SortFunc(allTasks, func(t1, t2 task.Task) int {
		if t1.Deadline.Equal(t2.Deadline) {
			return int(t1.ID - t2.ID)
		}
		return t1.Deadline.Compare(t2.Deadline)
	})
	return allTasks, nil
}

func (s *TaskService) UpdateTask(taskID task.ID, title *string, deadline *time.Time) (task.Task, error) {
	if title == nil && deadline == nil {
		return task.Task{}, task.ErrTaskFieldsEmpty
	}

	t, err := s.storage.GetTask(taskID)
	if err != nil {
		return task.Task{}, fmt.Errorf("получение задачи: %w", err)
	}

	if errStatus := s.checkNotClosed(t.Status); errStatus != nil {
		return task.Task{}, errStatus
	}

	if title != nil {
		if errTitle := s.validateTitle(*title); errTitle != nil {
			return task.Task{}, errTitle
		}
		t.Title = strings.TrimSpace(*title)
	}

	if deadline != nil {
		if errDeadline := s.validateDeadline(*deadline); errDeadline != nil {
			return task.Task{}, errDeadline
		}
		t.Deadline = *deadline
	}

	updateTaskErr := s.storage.UpdateTask(t)
	if updateTaskErr != nil {
		return task.Task{}, fmt.Errorf("обновление задачи: %w", updateTaskErr)
	}
	return t, nil
}

func (s *TaskService) ChangeTaskStatus(taskID task.ID, newStatus task.Status) (task.Task, error) {
	if !newStatus.IsValid() {
		return task.Task{}, task.ErrStatusNotAllowable
	}

	t, err := s.storage.GetTask(taskID)
	if err != nil {
		return task.Task{}, fmt.Errorf("получение задачи: %w", err)
	}

	if errStatus := s.checkNotClosed(t.Status); errStatus != nil {
		return task.Task{}, errStatus
	}

	if t.Status == newStatus {
		return task.Task{}, task.ErrStatusUnchanged
	}

	t.Status = newStatus
	updateTaskErr := s.storage.UpdateTask(t)
	if updateTaskErr != nil {
		return task.Task{}, fmt.Errorf("обновление задачи: %w", updateTaskErr)
	}
	return t, nil
}

func (s *TaskService) validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return task.ErrTaskEmptyTitle
	}
	return nil
}

func (s *TaskService) validateDeadline(deadline time.Time) error {
	year, month, day := time.Now().Date()
	startDayTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	if deadline.Before(startDayTime) {
		return task.ErrTaskDeadlineInPast
	}
	return nil
}

func (s *TaskService) checkNotClosed(status task.Status) error {
	if status == task.Canceled || status == task.Done {
		return task.ErrTaskClosed
	}
	return nil
}
