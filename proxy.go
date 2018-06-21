package reverseproxy

import (
	"github.com/casbin/casbin"
	"net/http/httputil"
	"net/url"
	"net/http"
	"strings"
)

type Proxy struct {
	Target string
	Proxy *httputil.ReverseProxy
	Host string
	Prefix string
	enforcer *casbin.Enforcer
}

func NewProxy(target string, prefix string, enforcer *casbin.Enforcer) *Proxy {
	u, _ := url.Parse(target)

	return &Proxy{
		Target: target,
		Proxy: httputil.NewSingleHostReverseProxy(u),
		Host: u.Host,
		Prefix: prefix,
		enforcer: enforcer,
	}
}

func (p *Proxy) Handle(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys) != 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	res := p.enforcer.Enforce(keys[0], r.URL.Path)

	if !res {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("X-GoProxy", "GoProxy")
	r.Host = p.Host

	r.URL.Path = strings.Replace(r.URL.Path, "/" + p.Prefix, "/", 1)

	p.Proxy.ServeHTTP(w, r)
}
