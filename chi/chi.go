package chi

import (
	"github.com/go-chi/chi"
)

type Proxy interface {
	Charge(mux *chi.Mux)
}
