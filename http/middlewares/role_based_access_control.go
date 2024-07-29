package middlewares

import (
	"net/http"
	"relif/bff/entities"
)

type RoleBasedAccessControl struct {
}

func NewRoleBasedAccessControl() *RoleBasedAccessControl {
	return &RoleBasedAccessControl{}
}

func (middleware *RoleBasedAccessControl) Middleware(allowedRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user").(entities.User)

			for _, role := range allowedRoles {
				if user.Role == role {
					next.ServeHTTP(w, r)
				}
			}

			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
