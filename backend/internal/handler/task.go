package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/IwantHappiness/todolist/internal/models"
	"github.com/IwantHappiness/todolist/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	scv service.TaskService
}

func NewTaskHandler(scv service.TaskService) TaskHandler {
	return TaskHandler{
		scv: scv,
	}
}

func (h *TaskHandler) RegisterRouter(r *gin.Engine) {
	tasks := r.Group("/tasks")
	{
		tasks.GET("", h.GetAll)
		tasks.GET("/:id", h.GetTaskHandler)
		tasks.POST("", h.CreateTaskHandler)
		tasks.DELETE("", h.DeleteAllTaskHandler)
		tasks.DELETE("/:id", h.DeleteTaskHandler)
		tasks.PUT("/:id", h.UpdateTaskHandler)
		tasks.PATCH("/:id", h.CompleteTaskHandler)
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
func (h *TaskHandler) GetAll(ctx *gin.Context) {
	query, ok := ctx.GetQuery("completed")

	if !ok {
		tasks, err := h.scv.GetAllTasks(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.ErrorDTO{Message: err.Error(), Time: time.Now()})
			return
		}
		ctx.JSON(http.StatusOK, tasks)
		return
	}

	completed, err := strconv.ParseBool(query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{Message: err.Error(), Time: time.Now()})
		return
	}

	tasks, err := h.scv.GetByCompleted(ctx.Request.Context(), completed)

	ctx.JSON(http.StatusOK, tasks)
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
func (h *TaskHandler) GetTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	taskId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	task, err := h.scv.GetTaskByID(ctx.Request.Context(), taskId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorDTO{
			Message: err.Error(),
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
func (h *TaskHandler) CreateTaskHandler(ctx *gin.Context) {
	var taskDTO models.TaskDTO
	if err := json.NewDecoder(ctx.Request.Body).Decode(&taskDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	if err := taskDTO.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	task, err := h.scv.CreateTask(ctx.Request.Context(), taskDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, task)
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
func (h *TaskHandler) UpdateTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	var taskDTO models.TaskDTO
	if err := json.NewDecoder(ctx.Request.Body).Decode(&taskDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	task, err := h.scv.UpdateTask(ctx.Request.Context(), taskId, taskDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	ctx.JSON(http.StatusOK, task)
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
func (h *TaskHandler) CompleteTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	var completeDTO models.CompleteTaskDTO
	if err := json.NewDecoder(ctx.Request.Body).Decode(&completeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	task, err := h.scv.CompleteTaskStatus(ctx.Request.Context(), taskId, completeDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	ctx.JSON(http.StatusOK, task)
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
func (h *TaskHandler) DeleteTaskHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	taskId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorDTO{Message: err.Error(), Time: time.Now()})
		return
	}

	if err := h.scv.DeleteTaskByID(ctx.Request.Context(), taskId); err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorDTO{Message: err.Error(), Time: time.Now()})
		return
	}

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
func (h *TaskHandler) DeleteAllTaskHandler(ctx *gin.Context) {
	if err := h.scv.DeleteAllTask(ctx.Request.Context()); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorDTO{Message: err.Error(), Time: time.Now()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
