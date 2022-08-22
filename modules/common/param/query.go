// Package param parse rawQuery
package param

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

// QueryString will get a query parameter by key.
func QueryString(r *http.Request, key string) string {
	params, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		logrus.Errorf("QueryUnescape error [%s]", err.Error())
		return ""
	}
	qr, err := url.ParseQuery(params)
	if err != nil {
		return ""
	}
	return qr.Get(key)
}

// QueryStringInt will get a query parameter by key and convert it to an int or return an error.
func QueryStringInt(r *http.Request, key string) (int, error) {
	val, err := strconv.Atoi(QueryString(r, key))
	if err != nil {
		return 0, err
	}
	return val, nil
}
