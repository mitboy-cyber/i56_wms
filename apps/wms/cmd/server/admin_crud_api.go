package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/i56/framework/core/router"
	rbacDomain "github.com/i56/modules/rbac/domain"
	rbacRepo "github.com/i56/modules/rbac/repository"
)

const defaultTenant int64 = 1

func registerAdminCRUDAPI(r *router.Router, repo *rbacRepo.MemRBACRepo) {
	r.GET("/api/v1/admin/crud", func(w http.ResponseWriter, req *http.Request) {
		entity := req.URL.Query().Get("entity")
		ctx := req.Context()
		offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
		if limit <= 0 {
			limit = 50
		}

		switch entity {
		case "permissions":
			list, total, err := repo.ListPermissions(ctx, offset, limit)
			if err != nil {
				writeError(w, err.Error(), 500)
				return
			}
			writeJSON(w, map[string]any{"data": list, "total": total})
		case "roles":
			list, total, err := repo.ListRoles(ctx, defaultTenant, offset, limit)
			if err != nil {
				writeError(w, err.Error(), 500)
				return
			}
			writeJSON(w, map[string]any{"data": list, "total": total})
		case "users":
			list, total, err := repo.ListUsers(ctx, defaultTenant, offset, limit)
			if err != nil {
				writeError(w, err.Error(), 500)
				return
			}
			writeJSON(w, map[string]any{"data": list, "total": total})
		case "client_permissions":
			list, total, err := repo.ListClientPermissions(ctx, defaultTenant, offset, limit)
			if err != nil {
				writeError(w, err.Error(), 500)
				return
			}
			writeJSON(w, map[string]any{"data": list, "total": total})
		default:
			writeError(w, "invalid entity: "+entity, 400)
		}
	})

	r.POST("/api/v1/admin/crud", func(w http.ResponseWriter, req *http.Request) {
		entity := req.URL.Query().Get("entity")
		action := req.URL.Query().Get("action")
		idStr := req.URL.Query().Get("id")

		if action == "" {
			writeError(w, "action required", 400)
			return
		}

		switch entity {
		case "permissions":
			handlePermissionAction(w, req, repo, action, idStr)
		case "roles":
			handleRoleAction(w, req, repo, action, idStr)
		case "users":
			handleUserAction(w, req, repo, action, idStr)
		case "client_permissions":
			handleClientPermissionAction(w, req, repo, action, idStr)
		default:
			writeError(w, "invalid entity: "+entity, 400)
		}
	})
}

func handlePermissionAction(w http.ResponseWriter, req *http.Request, repo *rbacRepo.MemRBACRepo, action, idStr string) {
	ctx := req.Context()

	switch action {
	case "create":
		var p rbacDomain.Permission
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if p.Slug == "" {
			writeError(w, "slug required", 400)
			return
		}
		if err := repo.CreatePermission(ctx, &p); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"data": p})

	case "update":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		var p rbacDomain.Permission
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if err := repo.UpdatePermission(ctx, id, &p); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"data": p})

	case "delete":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		if err := repo.DeletePermission(ctx, id); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"message": "ok"})

	default:
		writeError(w, "invalid action: "+action, 400)
	}
}

func handleRoleAction(w http.ResponseWriter, req *http.Request, repo *rbacRepo.MemRBACRepo, action, idStr string) {
	ctx := req.Context()

	switch action {
	case "create":
		var r rbacDomain.Role
		if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if r.Name == "" {
			writeError(w, "name required", 400)
			return
		}
		if err := repo.CreateRole(ctx, defaultTenant, &r); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"data": r})

	case "update":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		var r rbacDomain.Role
		if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if err := repo.UpdateRole(ctx, defaultTenant, id, &r); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"data": r})

	case "delete":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		if err := repo.DeleteRole(ctx, defaultTenant, id); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"message": "ok"})

	default:
		writeError(w, "invalid action: "+action, 400)
	}
}

func handleUserAction(w http.ResponseWriter, req *http.Request, repo *rbacRepo.MemRBACRepo, action, idStr string) {
	ctx := req.Context()

	switch action {
	case "create":
		var u rbacDomain.User
		if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if u.Username == "" {
			writeError(w, "username required", 400)
			return
		}
		if err := repo.CreateUser(ctx, defaultTenant, &u); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"data": u})

	case "update":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		var u rbacDomain.User
		if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if err := repo.UpdateUser(ctx, defaultTenant, id, &u); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"data": u})

	case "delete":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		if err := repo.DeleteUser(ctx, defaultTenant, id); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"message": "ok"})

	default:
		writeError(w, "invalid action: "+action, 400)
	}
}

func handleClientPermissionAction(w http.ResponseWriter, req *http.Request, repo *rbacRepo.MemRBACRepo, action, idStr string) {
	ctx := req.Context()

	switch action {
	case "create":
		var cp rbacDomain.ClientPermission
		if err := json.NewDecoder(req.Body).Decode(&cp); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if cp.ClientID <= 0 {
			writeError(w, "client_id required", 400)
			return
		}
		if err := repo.CreateClientPermission(ctx, defaultTenant, &cp); err != nil {
			writeError(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"data": cp})

	case "update":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		var cp rbacDomain.ClientPermission
		if err := json.NewDecoder(req.Body).Decode(&cp); err != nil {
			writeError(w, "invalid json body", 400)
			return
		}
		if err := repo.UpdateClientPermission(ctx, defaultTenant, id, &cp); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"data": cp})

	case "delete":
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeError(w, "invalid id", 400)
			return
		}
		if err := repo.DeleteClientPermission(ctx, defaultTenant, id); err != nil {
			writeError(w, err.Error(), 404)
			return
		}
		writeJSON(w, map[string]any{"message": "ok"})

	default:
		writeError(w, "invalid action: "+action, 400)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{"error": msg})
}
