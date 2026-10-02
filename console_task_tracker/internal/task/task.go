package task

import (
	"errors"
	"time"
)

type Task struct {
	ID       ID        `json:"task_id"`
	Status   Status    `json:"status"`
	Title    string    `json:"title"`
	Deadline time.Time `json:"deadline"`
}

type ID int
type Status string

const (
	Planned    Status = "planned"
	InProgress Status = "in_progress"
	Canceled   Status = "canceled"
	Done       Status = "done"
)

var ErrTaskNotFound = errors.New("задача не найдена")
var ErrTaskEmptyTitle = errors.New("заголовок не может быть пустым")
var ErrTaskDeadlineInPast = errors.New("время выполнения вышло")
var ErrTaskClosed = errors.New("задача завершена, операции изменения с ней выполнять нельзя")
var ErrStatusNotAllowable = errors.New("данный статус не разрешен")
var ErrStatusUnchanged = errors.New("статус уже такой")
var ErrTaskFieldsEmpty = errors.New("требуется отправить хотя бы одно из полей")

func (s Status) IsValid() bool {
	switch s {
	case Planned, InProgress, Canceled, Done:
		return true
	default:
		return false
	}
}
