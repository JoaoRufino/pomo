package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joao.rufino/pomo/pkg/core/models"
)

// PomodoroSave saves a pomodoro
func (s *RestServer) PomodoroSave() http.HandlerFunc {

	// swagger:operation POST /api/pomodoros/{id} PomodoroSave
	//
	// Create/Save Pomodoro
	//
	// Creates or saves a pomodoro. Omit the ID to auto generate.
	// Pass an existing ID to update.
	//
	// ---
	// parameters:
	// - name: pomodoro
	//   in: body
	//   description: Pomodoro to Save/Update
	//   required: true
	//   type: object
	//   schema:
	//     "$ref": "#/definitions/models_PomodoroExample"
	// responses:
	//   '200':
	//     description: Pomodoro Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/models_Pomodoro"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		id := chi.URLParam(r, "id")
		taskID, err := strconv.Atoi(id)
		if err != nil {
			RenderErrInvalidRequest(w, err)
		}

		var pomodoro = new(models.Pomodoro)
		if err := DecodeJSON(r.Body, pomodoro); err != nil {
			RenderErrInvalidRequest(w, err)
			return
		}

		err = s.store.PomodoroSave(ctx, taskID, pomodoro)
		if err != nil {
			if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpSave))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("PomodoroSave error", "error", err, "error_id", errID)
			}
			return
		}

		RenderJSON(w, http.StatusOK, pomodoro)
	}

}

// PomodoroGetByID returns the pomodoro
func (s *RestServer) PomodoroGetByID() http.HandlerFunc {

	// swagger:operation GET /api/pomodoros/{id} PomodoroGetByID
	//
	// Get a Pomodoro
	//
	// Fetches a Pomodoro
	//
	// ---
	// tags:
	// - WIDGETS
	// parameters:
	// - name: id
	//   in: path
	//   description: Pomodoro ID to fetch
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: Pomodoro Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/models_Pomodoro"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")
		taskID, _ := strconv.Atoi(id)

		pomodoro, err := s.store.PomodoroGetByTaskID(ctx, taskID)
		if err != nil {
			if err == models.ErrNotFound {
				RenderErrResourceNotFound(w, "pomodoro")
			} else if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpGet))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("PomodoroGetByID error", "error", err, "error_id", errID)
			}
			return
		}

		RenderJSON(w, http.StatusOK, pomodoro)
	}
}

// PomodoroDeleteByID deletes a pomodoro
func (s *RestServer) PomodoroDeleteByID() http.HandlerFunc {

	// swagger:operation DELETE /api/pomodoros/{id} PomodoroDeleteByID
	//
	// Delete a Pomodoro
	//
	// Deletes a Pomodoro
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   description: Pomodoro ID to delete
	//   type: string
	//   required: true
	// responses:
	//   '204':
	//     description: No Content
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")
		taskID, _ := strconv.Atoi(id)

		err := s.store.PomodoroDeleteByTaskID(ctx, taskID)
		if err != nil {
			if err == models.ErrNotFound {
				RenderErrResourceNotFound(w, "pomodoro")
			} else if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpDelete))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("PomodoroDeleteByID error", "error", err, "error_id", errID)
			}
			return
		}

		RenderNoContent(w)
	}
}

// GetStatus returns the server status
func (s *RestServer) StatusGet() http.HandlerFunc {
	// swagger:operation GET /api/status GetStatus
	//
	// Get the server status
	//
	// Fetches the server status
	//
	// ---
	// responses:
	//   '200':
	//     description: Status Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/models_Status"
	return func(w http.ResponseWriter, r *http.Request) {
		RenderJSON(w, http.StatusOK, s.status)
	}
}

// StatusSave saves the server status
func (s *RestServer) StatusSave() http.HandlerFunc {

	// swagger:operation POST /api/status StatusSave
	//
	// Save Status
	//
	// Saves the current server status
	//
	// ---
	// parameters:
	// - name: status
	//   in: body
	//   description: Status to Save/Update
	//   required: true
	//   type: object
	//   schema:
	//     "$ref": "#/definitions/models_Status"
	// responses:
	//   '200':
	//     description: Status Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/models_Status"
	return func(w http.ResponseWriter, r *http.Request) {

		var status = new(models.Status)
		if err := DecodeJSON(r.Body, status); err != nil {
			RenderErrInvalidRequest(w, err)
			return
		}
		s.status = *status

		RenderJSON(w, http.StatusOK, s.status)
	}

}
