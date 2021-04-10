package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joao.rufino/pomo/pkg/core/models"
)

// TaskSave saves a task
func (s *RestServer) TaskSave() http.HandlerFunc {

	// swagger:operation POST /api/tasks TaskSave
	//
	// Create/Save Task
	//
	// Creates or saves a task. Omit the ID to auto generate.
	// Pass an existing ID to update.
	//
	// ---
	// parameters:
	// - name: task
	//   in: body
	//   description: Task to Save/Update
	//   required: true
	//   type: object
	//   schema:
	//     "$ref": "#/definitions/models_TaskExample"
	// responses:
	//   '200':
	//     description: User Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/models_Task"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var task = models.Task{}
		if err := DecodeJSON(r.Body, task); err != nil {
			RenderErrInvalidRequest(w, err)
			return
		}

		err := s.store.TaskSave(ctx, task)
		if err != nil {
			if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpSave))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("TaskSave error", "error", err, "error_id", errID)
			}
			return
		}

		RenderJSON(w, http.StatusOK, task)
	}

}

// TaskGetByID returns the task
func (s *RestServer) TaskGetByID() http.HandlerFunc {

	// swagger:operation GET /api/tasks/{id} TaskGetByID
	//
	// Get a Task
	//
	// Fetches a Task
	//
	// ---
	// tags:
	// - THINGS
	// parameters:
	// - name: id
	//   in: path
	//   description: Task ID to fetch
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: Task Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/models_Task"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		task, err := s.store.TaskGetByID(ctx, id)
		if err != nil {
			if err == models.ErrNotFound {
				RenderErrResourceNotFound(w, "task")
			} else if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpGet))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("TaskGetByID error", "error", err, "error_id", errID)
			}
			return
		}

		RenderJSON(w, http.StatusOK, task)
	}

}

// TaskDeleteByID deletes a task
func (s *RestServer) TaskDeleteByID() http.HandlerFunc {

	// swagger:operation DELETE /api/tasks/{id} TaskDeleteByID
	//
	// Delete a Task
	//
	// Deletes a Task
	//
	// ---
	// tags:
	// - THINGS
	// parameters:
	// - name: id
	//   in: path
	//   description: Task ID to delete
	//   type: string
	//   required: true
	// responses:
	//   '204':
	//     description: No Content
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		err := s.store.TaskDeleteByID(ctx, id)
		if err != nil {
			if err == models.ErrNotFound {
				RenderErrResourceNotFound(w, "task")
			} else if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpDelete))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("TaskDeleteByID error", "error", err, "error_id", errID)
			}
			return
		}

		RenderNoContent(w)

	}

}

// TasksFind finds tasks
func (s *RestServer) TasksFind() http.HandlerFunc {

	// swagger:operation GET /api/tasks TasksFind
	//
	// Find Tasks
	//
	// Gets a list of tasks
	//
	// ---
	// parameters:
	// - name: limit
	//   in: query
	//   description: Number of records to return
	//   type: int
	//   required: false
	// - name: offset
	//   in: query
	//   description: Offset of records to return
	//   type: int
	//   required: false
	// - name: id
	//   in: query
	//   description: Filter id
	//   type: string
	//   required: false
	// - name: name
	//   in: query
	//   description: Filter name
	//   type: string
	//   required: false
	// responses:
	//   '200':
	//     description: Task Objects
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/models_Task"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		tasks, count, err := s.store.TasksFind(ctx, qp)
		if err != nil {
			if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpFind))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("TasksFind error", "error", err, "error_id", errID)
			}
			return
		}

		RenderJSON(w, http.StatusOK, models.Results{Count: count, Results: tasks})

	}

}
