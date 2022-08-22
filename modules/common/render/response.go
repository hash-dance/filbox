// Package render send data
package render

import (
	"io"
	"net/http"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
	"github.com/go-chi/render"
)

// Response defined Response body
type Response struct {
	Code       ErrNumber   `json:"code"`
	ErrMessage string      `json:"errMessage"`
	Data       interface{} `json:"data"`
}

// SendError wrap send error
func SendError(w http.ResponseWriter, r *http.Request, code ErrorCode, err error) {
	render.JSON(w, r, Response{code.ErrNumber, err.Error(), nil})
}

// SendJSON send json data
func SendJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.JSON(w, r, Response{Success, "", data})
}

// SendPaginationJSON send json data with pagination
func SendPaginationJSON(w http.ResponseWriter, r *http.Request, pagination *mysql.Pagination, data interface{}) {
	SendJSON(w, r, map[string]interface{}{
		"pagination": pagination,
		"data":       data,
	})
}

// DecodeJSON decode to json
func DecodeJSON(r io.Reader, v interface{}) error {
	return render.DecodeJSON(r, v)
}
