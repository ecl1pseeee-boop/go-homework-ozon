package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task_tracker/internal/service"
	"task_tracker/internal/task"
	"time"
)

type App struct {
	service *service.TaskService
}

func NewApp(service *service.TaskService) *App {
	return &App{service: service}
}

const DateLayout = "2006-01-02"

func (c *App) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		fields := strings.Fields(userInput)
		if len(fields) == 0 {
			fmt.Println("Введите команду")
			continue
		}

		var err error

		switch fields[0] {
		case "help":
			c.printCommandsInfo()
		case "add":
			err = c.addTask(fields)
		case "list":
			err = c.getTaskList()
		case "task":
			err = c.getTask(fields)
		case "status":
			err = c.changeStatus(fields)
		case "edit":
			err = c.editTask(fields)
		case "delete":
			err = c.deleteTask(fields)
		case "exit":
			return nil
		default:
			fmt.Println("Такой команды не существует. Введите другую.")
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка выполнения команды: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка: %w", err)
	}
	return nil
}

func (c *App) printCommandsInfo() {
	helpText := `Доступные команды Task Tracker:
  help                                              - Показать эту справку
  exit                                              - Выйти из программы
  add <title> <YYYY-MM-DD>                          - Добавить новую задачу
  list                                              - Показать список всех задач
  task <id>                                         - Посмотреть детали задачи по ID
  status <id> <status>                              - Изменить статус задачи
  edit <id> --title <title> --deadline <YYYY-MM-DD> - Отредактировать задачу
  delete <id>                                       - Удалить задачу по ID`
	fmt.Println(helpText)
}

func (c *App) addTask(fields []string) error {
	fmt.Println("Создание задачи")

	title, parsedDate, parseErr := c.parseAddCmd(fields)
	if parseErr != nil {
		return fmt.Errorf("ошибка парсинга данных из команды: %w", parseErr)
	}
	t, err := c.service.CreateTask(title, parsedDate)
	if err != nil {
		return fmt.Errorf("ошибка добавления задачи: %w", err)
	}
	c.printTask(t)
	return nil
}

func (c *App) parseAddCmd(fields []string) (string, time.Time, error) {
	if len(fields) < 3 {
		return "", time.Time{}, fmt.Errorf("неправильно введена команда, требуется формат 'add <title> <YYYY-MM-DD>'")
	}
	title := strings.Join(fields[1:len(fields)-1], " ")
	date := fields[len(fields)-1]
	parsedDate, parseErr := time.ParseInLocation(DateLayout, date, time.Local)
	if parseErr != nil {
		return "", time.Time{}, fmt.Errorf("ошибка конвертации даты: %w", parseErr)
	}
	return title, parsedDate, nil
}

func (c *App) getTaskList() error {
	fmt.Println("Получение списка задач")
	tasks, err := c.service.GetAllTasks()
	if err != nil {
		return fmt.Errorf("ошибка получения списка задач: %w", err)
	}
	if len(tasks) == 0 {
		fmt.Printf("Список задач пустой. \n")
		return nil
	}
	for _, t := range tasks {
		c.printTask(t)
	}
	return nil
}

func (c *App) getTask(fields []string) error {
	fmt.Println("Получение задачи")

	if len(fields) != 2 {
		return fmt.Errorf("неправильно введена команда, требуется формат task <id>")
	}
	taskID, idConvErr := strconv.Atoi(fields[1])
	if idConvErr != nil {
		return fmt.Errorf("ошибка конвертации id: %w", idConvErr)
	}
	t, err := c.service.GetTask(task.ID(taskID))
	if err != nil {
		return fmt.Errorf("ошибка получения задачи: %w", err)
	}
	c.printTask(t)
	return nil
}

func (c *App) changeStatus(fields []string) error {
	fmt.Println("Изменение статуса задачи")

	id, status, parseErr := c.parseStatusCmd(fields)
	if parseErr != nil {
		return fmt.Errorf("ошибка парсинга команды: %w", parseErr)
	}
	t, err := c.service.ChangeTaskStatus(task.ID(id), task.Status(status))
	if err != nil {
		return fmt.Errorf("ошибка изменения статуса задачи: %w", err)
	}
	c.printTask(t)
	return nil
}

func (c *App) parseStatusCmd(fields []string) (int, string, error) {
	if len(fields) != 3 {
		return 0, "", fmt.Errorf("неправильно введена команда, требуется формат 'status <id> <status>'")
	}
	taskId, idConvErr := strconv.Atoi(fields[1])
	if idConvErr != nil {
		return 0, "", fmt.Errorf("ошибка конвертации id: %w", idConvErr)
	}
	status := fields[2]
	return taskId, status, nil
}

func (c *App) editTask(fields []string) error {
	fmt.Println("Изменение задачи")

	id, title, deadline, parseErr := c.parseUpdateCmd(fields)
	if parseErr != nil {
		return fmt.Errorf("ошибка получения данных из команды или конвертации: %w", parseErr)
	}
	t, updateErr := c.service.UpdateTask(task.ID(id), title, deadline)
	if updateErr != nil {
		return fmt.Errorf("ошибка обновления данных: %w", updateErr)
	}
	c.printTask(t)
	return nil
}

func (c *App) parseUpdateCmd(fields []string) (int, *string, *time.Time, error) {
	if len(fields) < 4 {
		return 0, nil, nil, fmt.Errorf("неправильно введена команда, требуется формат edit <id> --title <title> --deadline <YYYY-MM-DD>")
	}
	id, errConv := strconv.Atoi(fields[1])
	if errConv != nil {
		return 0, nil, nil, fmt.Errorf("неправильно введен id: %w", errConv)
	}

	var title *string
	var deadline *time.Time

	for i := 2; i < len(fields); i++ {
		if fields[i] == "--deadline" {
			if i >= len(fields)-1 {
				return 0, nil, nil, fmt.Errorf("необходимо ввести дату")
			}
			deadlineStr := fields[i+1]
			parsedDeadline, parseErr := time.ParseInLocation(DateLayout, deadlineStr, time.Local)
			if parseErr != nil {
				return 0, nil, nil, fmt.Errorf("ошибка парсинга даты: %w", parseErr)
			}
			deadline = &parsedDeadline
		}
		if fields[i] == "--title" {
			startIndex := i + 1
			finalIndex := -1
			for i = i + 1; i < len(fields); i++ {
				value := fields[i]
				if strings.HasPrefix(value, "--deadline") {
					break
				}
				finalIndex = i + 1
			}
			if finalIndex == -1 {
				return 0, nil, nil, fmt.Errorf("после --title нужно указать заголовок")
			}
			parsedTitle := strings.Join(fields[startIndex:finalIndex], " ")
			i = finalIndex - 1
			title = &parsedTitle
		}
	}

	return id, title, deadline, nil
}

func (c *App) deleteTask(fields []string) error {
	fmt.Println("Удаление задачи")

	if len(fields) != 2 {
		return fmt.Errorf("неправильно введена команда, требуется формат 'delete <id>'")
	}
	taskID, idConvErr := strconv.Atoi(fields[1])
	if idConvErr != nil {
		return fmt.Errorf("ошибка конвертации id: %w", idConvErr)
	}
	err := c.service.RemoveTask(task.ID(taskID))
	if err != nil {
		return fmt.Errorf("ошибка удаление задачи: %w", err)
	}
	fmt.Printf("Задача с id %d удалена успешно. \n", taskID)
	return nil
}

func (c *App) printTask(task task.Task) {
	fmt.Printf("ID: %d. Title: %s. Deadline: %s. Status: %s\n", task.ID, task.Title, task.Deadline.Format(DateLayout), task.Status)
}
