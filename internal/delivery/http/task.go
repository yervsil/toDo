package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yervsil/toDo-microservice/internal/entity"
)


// @Summary Create todo item
// @Tags tasks
// @Description Create a new todo item
// @Accept json
// @Produce json
// @Param input body entity.Task true "Task information"
// @Success 200 {integer} integer 1
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /api/todo-list/tasks [post]

// Создать задачу
func (h *Handler) createTask(c *gin.Context) {
	var input entity.Task

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error(err)
		errorResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if !isValidDateFormat(input.ActiveAt){
		h.logger.Error("incorrect date format")
		errorResponse(c, http.StatusBadRequest, "incorrect date format")

		return
	}


	id, err := h.service.CreateTask(c.Request.Context(), input)
	if err != nil {
		h.logger.Error(err)
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
// @Summary Update todo item
// @Tags tasks
// @Description Update an existing todo item
// @ID update-task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param input body entity.Task true "Updated task information"
// @Success 201 {string} string "Successfully updated"
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /api/todo-list/tasks/{int} [put]

// Заменить задачу по id
func (h *Handler) updateTask(c *gin.Context) {
	taskId, err := parseIdFromPath(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")

		return
	}
	
	var input entity.Task

	if err := c.BindJSON(&input); err != nil {
		h.logger.Error(err)
		errorResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if !isValidDateFormat(input.ActiveAt){
		h.logger.Error("incorrect date format")
		errorResponse(c, http.StatusBadRequest, "incorrect date format")

		return
	}


	err = h.service.UpdateTask(c.Request.Context(), input, taskId)

	if err != nil {
		h.logger.Error(err)
		
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	c.JSON(http.StatusCreated, "successfully updated")
}

// @Summary Delete todo item
// @Tags tasks
// @Description Delete an existing todo item
// @ID delete-task
// @Param id path string true "Task ID"
// @Success 201 {string} string "Successfully deleted"
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /api/todo-list/tasks/{id} [delete]

// Удалить задачу по id
func (h *Handler) deleteTask(c *gin.Context) {
	taskId, err := parseIdFromPath(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")

		return
	}

	err = h.service.DeleteTask(c.Request.Context(), taskId)
	if err != nil {
		h.logger.Error(err)
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	c.JSON(http.StatusCreated, "successfully deleted")
}

// @Summary Update status of todo item
// @Tags tasks
// @Description Update status of an existing todo item
// @ID update-status
// @Param id path string true "Task ID"
// @Success 201 {string} string "Status has been changed"
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /api/todo-list/tasks/{id}/done [patch]

// Обновить статус задачи на выполнено по id
func (h *Handler) statusUpdate(c *gin.Context) {
	taskId, err := parseIdFromPath(c, "id")
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")

		return
	}

	err = h.service.StatusUpdate(c.Request.Context(), taskId)
	if err != nil {
		h.logger.Error(err)
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	c.JSON(http.StatusCreated, "status has been changed")
}

// @Summary Get todo items
// @Tags tasks
// @Description Get a list of todo items
// @Accept json
// @Produce json
// @Param status query string false "Status filter: active or done"
// @Success 200 {array} entity.Task "List of todo items"
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Router /api/todo-list/tasks [get]

// Получить все задачи взависимости от статуса
func (h *Handler) getTasks(c *gin.Context) {
	status := c.DefaultQuery("status", "active")

	tasks, err := h.service.GetTasks(c.Request.Context(), status)
	if err != nil {
		h.logger.Error(err)
		errorResponse(c, http.StatusNotFound, err.Error())

		return
	}

	if len(tasks) == 0 {
        tasks = []entity.Task{}
    }

	c.JSON(http.StatusOK, tasks)
}

func isValidDateFormat(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}