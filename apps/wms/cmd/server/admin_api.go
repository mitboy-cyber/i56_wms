package main

import (
	"encoding/json"
	"net/http"

	"github.com/i56/i56-apps/i56-wms/internal/auth"
	"github.com/i56/framework/core/router"
	domain "github.com/i56/modules/rbac/domain"
	rbacRepo "github.com/i56/modules/rbac/repository"
)

// registerJSONAPI adds JSON API endpoints for the React frontend.
func registerJSONAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc, rbac *rbacRepo.MemRBACRepo, sessionMgr *auth.SessionManager) {
	const tenant int64 = 1

	// GET /admin/api/me — return current session user as JSON
	r.GET("/admin/api/me", a(func(w http.ResponseWriter, r *http.Request) {
		ck, err := r.Cookie("admin_session")
		if err != nil || ck.Value == "" {
			apiJSON(w, 401, map[string]string{"error": "unauthorized"})
			return
		}
		sess := sessionMgr.ValidateSession(ck.Value)
		if sess == nil {
			apiJSON(w, 401, map[string]string{"error": "invalid_session"})
			return
		}
		// Look up user by username from the session
		users, _, _ := rbac.ListUsers(r.Context(), tenant, 0, 200)
		var u *domain.User
		for i := range users {
			if users[i].Username == sess.Username {
				u = &users[i]
				break
			}
		}
		if u == nil {
			apiJSON(w, 401, map[string]string{"error": "user_not_found"})
			return
		}
		roles, _, _ := rbac.ListRoles(r.Context(), tenant, 0, 200)
		roleName := ""
		for _, ro := range roles {
			if ro.ID == u.RoleID {
				roleName = ro.Name
				break
			}
		}
		apiJSON(w, 200, map[string]interface{}{
			"id": u.ID, "username": u.Username, "real_name": u.RealName,
			"role_id": u.RoleID, "role_name": roleName,
			"email": u.Email, "phone": u.Phone, "is_active": u.IsActive,
		})
	}))

	// GET /admin/api/employees — list all employees as JSON
	r.GET("/admin/api/employees", a(func(w http.ResponseWriter, r *http.Request) {
		users, _, _ := rbac.ListUsers(r.Context(), tenant, 0, 200)
		roles, _, _ := rbac.ListRoles(r.Context(), tenant, 0, 200)
		roleNames := map[int64]string{}
		for _, ro := range roles {
			roleNames[ro.ID] = ro.Name
		}
		out := make([]map[string]interface{}, len(users))
		for i, u := range users {
			out[i] = map[string]interface{}{
				"id": u.ID, "username": u.Username, "real_name": u.RealName,
				"role_id": u.RoleID, "role_name": roleNames[u.RoleID],
				"email": u.Email, "phone": u.Phone, "is_active": u.IsActive,
			}
		}
		apiJSON(w, 200, out)
	}))

	// GET /admin/api/roles — list all roles as JSON
	r.GET("/admin/api/roles", a(func(w http.ResponseWriter, r *http.Request) {
		roles, _, _ := rbac.ListRoles(r.Context(), tenant, 0, 200)
		out := make([]map[string]interface{}, len(roles))
		for i, ro := range roles {
			out[i] = map[string]interface{}{
				"id": ro.ID, "name": ro.Name, "description": ro.Description, "is_active": ro.IsActive,
			}
		}
		apiJSON(w, 200, out)
	}))
}

func apiJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
