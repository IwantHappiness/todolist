package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mtx sync.RWMutex
var tasks map[int]Task = make(map[int]Task)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func NewTask(title, description string) Task {
	var id int
	for i := range tasks {
		if i > id {
			id = i
		}
	}
	id++

	return Task{
		ID:          id,
		Title:       title,
		CreatedAt:   time.Now(),
		Description: description,
		Completed:   false,
		CompletedAt: nil,
	}
}

/*
 * pattern /task
 * query param: completed = true|false (optional)
 * method: GET
 * response body: Json representation of all tasks
 *
 * failed:
 * - status code: 400
 * - response body: Json with error + time
 */
func GetAllTaskHandler(ctx *gin.Context) {
	query, ok := ctx.GetQuery("completed")

	if !ok {
		mtx.RLock()
		defer mtx.RUnlock()

		ctx.JSON(http.StatusOK, tasks)
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	completed, err := strconv.ParseBool(query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorDTO{Message: err.Error(), Time: time.Now()})
		return
	}

	var CompletedTasks map[int]Task = make(map[int]Task)
	for _, task := range tasks {
		if completed == task.Completed {
			CompletedTasks[task.ID] = task
		}
	}

	ctx.JSON(http.StatusOK, CompletedTasks)
}

/*
 * pattern /task/{id}
 * method: GET
 * response body: Json representation of task
 *
 * failed:
 * - status code: 400, 404
 * - response body: Json with error + time
 */
func GetTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	taskId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	mtx.RLock()
	task, ok := tasks[taskId]
	mtx.RUnlock()

	if !ok {
		ctx.JSON(http.StatusNotFound, ErrorDTO{
			Message: "Task not found",
			Time:    time.Now(),
		})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

/*
 * pattern /task
 * method: POST
 * response body: Json representation of task creation
 *
 * failed:
 * - status code: 400
 * - response body: Json with error + time
 */
func CreateTaskHandler(ctx *gin.Context) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(ctx.Request.Body).Decode(&taskDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	if err := taskDTO.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	task := NewTask(taskDTO.Title, taskDTO.Description)

	mtx.Lock()
	defer mtx.Unlock()
	tasks[task.ID] = task

	ctx.JSON(http.StatusCreated, tasks[task.ID])
}

/*
 * pattern /task/{id}
 * method: PATCH
 * response body: Json representation of task
 *
 * failed:
 * - status code: 400, 404
 * - response body: Json with error + time
 */
func CompleteTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	task, ok := tasks[taskId]
	if !ok {
		ctx.JSON(http.StatusNotFound, ErrorDTO{
			Message: "Task not found",
			Time:    time.Now(),
		})
		return
	}

	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(ctx.Request.Body).Decode(&completeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	task.Completed = completeDTO.Completed

	if completeDTO.Completed {
		now := time.Now()
		task.CompletedAt = &now
	} else {
		task.CompletedAt = nil
	}

	tasks[taskId] = task
	ctx.JSON(http.StatusOK, tasks[taskId])
}

/*
 * pattern /task/{id}
 * method: DELETE
 * response body: Json representation of task deletion
 *
 * failed:
 * - status code: 400, 404
 * - response body: Json with error + time
 */
func DeleteTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	taskId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorDTO{Message: err.Error(), Time: time.Now()})
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	_, ok := tasks[taskId]
	if !ok {
		ctx.JSON(http.StatusNotFound, ErrorDTO{Message: "Task not found", Time: time.Now()})
		return
	}

	delete(tasks, taskId)

	ctx.Status(http.StatusNoContent)
}

/*
 * pattern /task
 * method: DELETE
 * info: -
 *
 * success:
 * - status code: 200
 *
 */
func DeleteAllTaskHandler(ctx *gin.Context) {
	mtx.Lock()
	for k := range tasks {
		delete(tasks, k)
	}
	mtx.Unlock()

	ctx.Status(http.StatusNoContent)
}
