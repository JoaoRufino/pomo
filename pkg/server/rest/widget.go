package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joao.rufino/pomo/pkg/core/models"
)

// PomodoroSave saves a pomodoro
func (s *RestServer) PomodoroSave() http.HandlerFunc {

	// swagger:operation POST /api/pomodoros PomodoroSave
	//
	// Create/Save Pomodoro
	//
	// Creates or saves a pomodoro. Omit the ID to auto generate.
	// Pass an existing ID to update.
	//
	// ---
	// tags:
	// - WIDGETS
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
	//     description: User Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/models_Pomodoro"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var pomodoro = new(models.Pomodoro)
		if err := DecodeJSON(r.Body, pomodoro); err != nil {
			RenderErrInvalidRequest(w, err)
			return
		}

		err := s.grStore.PomodoroSave(ctx, pomodoro)
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

		pomodoro, err := s.grStore.PomodoroGetByID(ctx, id)
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
	// tags:
	// - WIDGETS
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

		err := s.grStore.PomodoroDeleteByID(ctx, id)
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

// PomodorosFind finds pomodoros
func (s *RestServer) PomodorosFind() http.HandlerFunc {

	// swagger:operation GET /api/pomodoros PomodorosFind
	//
	// Find Pomodoros
	//
	// Gets a list of pomodoros
	//
	// ---
	// tags:
	// - WIDGETS
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
	//     description: Pomodoro Objects
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/models_Pomodoro"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		pomodoros, count, err := s.grStore.PomodorosFind(ctx, qp)
		if err != nil {
			if serr, ok := err.(*models.Error); ok {
				RenderErrInvalidRequest(w, serr.ErrorForOp(models.ErrorOpFind))
			} else {
				errID := RenderErrInternalWithID(w, nil)
				s.logger.Errorw("PomodorosFind error", "error", err, "error_id", errID)
			}
			return
		}

		RenderJSON(w, http.StatusOK, models.Results{Count: count, Results: pomodoros})
	}
}
