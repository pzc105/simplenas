package task

type TaskStatus int

const (
	TaskStatusIniting = iota
	TaskStatusRunning
	TaskStatusInitDone
	TaskStatusInitFailed
)

type ITask interface {
	GetStatus() TaskStatus
	GetProgress() float32
	Delete() error
}

type RawTask struct {
	ITask
	status   TaskStatus
	progress float32
}
