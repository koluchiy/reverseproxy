package access

import "net/http"

type Manager interface {
	CheckAccess(r *http.Request)(bool, error)
}