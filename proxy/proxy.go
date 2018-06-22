package proxy

import (
	"net/http"
	"bitbucket.org/carprice/reverseproxy/access"
	"strings"
	"net/http/httputil"
	"net/url"
	"github.com/go-chi/chi"
)

type PrefixProxy interface {
	GetPrefix() string
	Handle(w http.ResponseWriter, r *http.Request)
}

func NewPrefixProxy(target string, prefix string, access access.Manager) (*prefixProxy, error) {
	u, err := url.Parse(target)

	if err != nil {
		return &prefixProxy{}, err
	}

	return &prefixProxy{
		prefix: prefix,
		host: u.Host,
		proxy: httputil.NewSingleHostReverseProxy(u),
		access: access,
	}, nil
}

type prefixProxy struct {
	access access.Manager
	prefix string
	host string
	proxy *httputil.ReverseProxy
}

func (p *prefixProxy) GetPrefix() string {
	return p.prefix
}

func (p *prefixProxy) Handle(w http.ResponseWriter, r *http.Request) {
	res, err := p.access.CheckAccess(r)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !res {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("X-GoProxy", "GoProxy")
	r.Host = p.host

	r.URL.Path = strings.Replace(r.URL.Path, "/" + p.prefix, "/", 1)

	p.proxy.ServeHTTP(w, r)
}

func (p *prefixProxy) Charge(mux *chi.Mux) {
	mux.HandleFunc("/" + p.GetPrefix() + "/*", p.Handle)
}

