/*Package host handler request for ui
 */
package host

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// Router handler /index request
// this is a static file server
func Router() http.Handler {
	root := "./static/ui"
	fs := http.FileServer(http.Dir(root))

	r := chi.NewRouter()
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("host request uri %s", r.RequestURI)
		if _, err := os.Stat(root + r.RequestURI); !os.IsNotExist(err) {
			fs.ServeHTTP(w, r)
			return
		}
		r.RequestURI = "/"
		http.Redirect(w, r, "/", http.StatusFound)
	}))
	return r
}
