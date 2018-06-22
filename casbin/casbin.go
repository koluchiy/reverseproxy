package casbin

import (
	"net/http"
	"github.com/casbin/casbin"
	"github.com/koluchiy/reverseproxy/access"
	"github.com/koluchiy/reverseproxy/extractor"
)

type manager struct {
	enforcer *casbin.Enforcer
	keyExtractor extractor.KeyExtractor
}

func (m *manager) CheckAccess(r *http.Request)(bool, error) {
	key, err := m.keyExtractor.GetKey(r)

	if err != nil {
		return false, err
	}

	res := m.enforcer.Enforce(key, r.URL.Path)

	return res, nil
}

func NewCasbinAccessManager(enforcer *casbin.Enforcer, keyExtractor extractor.KeyExtractor) access.Manager {
	return &manager{enforcer: enforcer, keyExtractor: keyExtractor}
}