package content

import (
	"errors"
	"net/http"
	"txp/restapistarter/internal/app/module/content/entity"
	"txp/restapistarter/internal/pkg/constant"
	"txp/restapistarter/pkg/adapter"
	datasql "txp/restapistarter/pkg/data/sql"
	"txp/restapistarter/pkg/response"
	"txp/restapistarter/pkg/time"

	"github.com/go-chi/chi"
)

type ContentService struct {
	repository *ContentRepository
}

func NewContentService(
	repository *ContentRepository,
) *ContentService {
	s := new(ContentService)
	s.repository = repository
	return s
}

func (s *ContentService) Create(p []byte, w http.ResponseWriter, r *http.Request) {
	d, err := adapter.BytesToValue[entity.Content](p)
	if err != nil {
		response.RespondError(http.StatusBadRequest, err, w)
		return
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	err = s.repository.Create(d)
	if err != nil {
		response.RespondError(
			http.StatusInternalServerError,
			errors.New(constant.InternalServerError),
			w,
		)
		return
	}
	response.Respond(http.StatusCreated, d, w)
}

func (s *ContentService) ReadMany(limit, page int, w http.ResponseWriter, r *http.Request) {
	rows, err := s.repository.ReadMany()
	if err != nil {
		response.RespondError(
			http.StatusInternalServerError,
			err,
			w,
		)
		return
	}
	var e entity.Content
	d, err := datasql.GetEntities(
		rows,
		&e,
		&e.Id,
		&e.Name,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		response.RespondError(
			http.StatusInternalServerError,
			err,
			w,
		)
		return
	}
	response.Respond(http.StatusOK, d, w)
}

func (s *ContentService) ReadOne(id string, w http.ResponseWriter, r *http.Request) {
	row := s.repository.ReadOne(id)
	if row == nil {
		response.RespondError(
			http.StatusInternalServerError,
			errors.New(constant.InternalServerError),
			w,
		)
		return
	}
	e := new(entity.Content)
	d, err := datasql.GetEntity(
		row,
		&e,
		&e.Id,
		&e.Name,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		response.RespondError(
			http.StatusInternalServerError,
			err,
			w,
		)
		return
	}
	response.Respond(http.StatusOK, d, w)
}

func (s *ContentService) Update(id string, p []byte, w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, constant.KeyId)
	d, err := adapter.BytesToValue[entity.Content](p)
	if err != nil {
		response.RespondError(http.StatusBadRequest, err, w)
		return
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	rowsAffected, err := s.repository.Update(userId, d)
	if err != nil || rowsAffected <= 0 {
		response.RespondError(
			http.StatusInternalServerError,
			errors.New(constant.InternalServerError),
			w,
		)
		return
	}
	response.Respond(http.StatusOK, d, w)
}

func (s *ContentService) Delete(id string, w http.ResponseWriter, r *http.Request) {
	rowsAffected, err := s.repository.Delete(id)
	if err != nil || rowsAffected <= 0 {
		response.RespondError(
			http.StatusInternalServerError,
			errors.New(constant.InternalServerError),
			w,
		)
		return
	}
	response.Respond(
		http.StatusOK,
		map[string]bool{"success": true},
		w,
	)
}
