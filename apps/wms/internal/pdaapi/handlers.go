// Package pdaapi provides PDA (handheld device) JSON API handlers.
package pdaapi

import (
	"encoding/json"
	"net/http"

	"github.com/i56/framework/core/router"

	orderRepo "github.com/i56/modules/order/repository"
	parcelRepo "github.com/i56/modules/parcel/repository"
	pdaRepo "github.com/i56/modules/pda/repository"
	pdaSvc "github.com/i56/modules/pda/service"
)

// RegisterPDAAPI registers all PDA JSON API endpoints.
func RegisterPDAAPI(
	r *router.Router,
	pdaR *pdaRepo.MemPDARepo,
	ops *pdaSvc.PDAOperations,
	pr *parcelRepo.MemParcelRepo,
	or *orderRepo.MemOrderRepo,
) {
	svc := pdaSvc.NewPDAService(pdaR, pr, or)

	pdaAuth := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ck, err := r.Cookie("pda_token")
			if err != nil || pdaR.ValidateSession(ck.Value) == nil {
				apiJSON(w, 401, map[string]string{"error": "unauthorized"})
				return
			}
			next(w, r)
		}
	}

	getOpID := func(r *http.Request) int64 {
		ck, _ := r.Cookie("pda_token")
		if ck == nil {
			return 1
		}
		sess := pdaR.ValidateSession(ck.Value)
		if sess == nil {
			return 1
		}
		return sess.OperatorID
	}

	getToken := func(r *http.Request) string {
		ck, _ := r.Cookie("pda_token")
		if ck == nil {
			return ""
		}
		return ck.Value
	}

	// Login
	r.POST("/pda/api/login", func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			Code string `json:"code"`
			PIN  string `json:"pin"`
		}
		json.NewDecoder(req.Body).Decode(&body)
		sess, err := svc.Login(body.Code, body.PIN, req.RemoteAddr)
		if err != nil {
			apiJSON(w, 401, map[string]string{"error": err.Error()})
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "pda_token", Value: sess.Token, Path: "/pda", HttpOnly: true, MaxAge: 43200})
		apiJSON(w, 200, map[string]interface{}{"token": sess.Token, "operator_id": sess.OperatorID})
	})

	// Logout
	r.POST("/pda/api/logout", func(w http.ResponseWriter, req *http.Request) {
		if ck, err := req.Cookie("pda_token"); err == nil {
			pdaR.Logout(ck.Value)
		}
		http.SetCookie(w, &http.Cookie{Name: "pda_token", Value: "", Path: "/pda", MaxAge: -1})
		apiJSON(w, 200, map[string]string{"ok": "logged_out"})
	})

	// Me
	r.GET("/pda/api/me", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, map[string]int64{"operator_id": getOpID(req)})
	}))

	// Dashboard
	r.GET("/pda/api/dashboard", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		stats := ops.WarehouseStats(req.Context())
		logs := ops.RecentLogs(10)
		apiJSON(w, 200, map[string]interface{}{"stats": stats, "recent_logs": logs, "op_id": getOpID(req)})
	}))

	// Receive
	r.POST("/pda/api/receive", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			Scan   string  `json:"scan"`
			Weight float64 `json:"weight"`
			Length float64 `json:"length"`
			Width  float64 `json:"width"`
			Height float64 `json:"height"`
		}
		json.NewDecoder(req.Body).Decode(&body)
		token := getToken(req)
		p, err := svc.ReceiveParcel(req.Context(), token, body.Scan, body.Weight, body.Length, body.Width, body.Height)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, p)
	}))

	// Weigh
	r.POST("/pda/api/weigh", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Scan string; Weight float64 }
		json.NewDecoder(req.Body).Decode(&body)
		opID := getOpID(req)
		p, actual, err := ops.Weigh(req.Context(), opID, body.Scan, body.Weight)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, map[string]interface{}{"parcel": p, "actual_weight": actual})
	}))

	// Putaway
	r.POST("/pda/api/putaway", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Scan, LocationBarcode string }
		json.NewDecoder(req.Body).Decode(&body)
		opID := getOpID(req)
		p, err := ops.PutAway(req.Context(), opID, body.Scan, body.LocationBarcode)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, p)
	}))

	// Pick
	r.POST("/pda/api/pick", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ OrderNo string }
		json.NewDecoder(req.Body).Decode(&body)
		opID := getOpID(req)
		order, parcels, err := ops.Pick(req.Context(), opID, body.OrderNo)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, map[string]interface{}{"order": order, "parcels": parcels})
	}))

	// Pack
	r.POST("/pda/api/pack", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ OrderNo string }
		json.NewDecoder(req.Body).Decode(&body)
		opID := getOpID(req)
		order, err := ops.Pack(req.Context(), opID, body.OrderNo)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, order)
	}))

	// Load
	r.POST("/pda/api/load", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ ContainerNo, OrderNo string }
		json.NewDecoder(req.Body).Decode(&body)
		opID := getOpID(req)
		err := ops.LoadContainer(req.Context(), opID, body.ContainerNo, body.OrderNo)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, map[string]string{"ok": "loaded"})
	}))

	// Exception
	r.POST("/pda/api/exception", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Scan, Reason string }
		json.NewDecoder(req.Body).Decode(&body)
		opID := getOpID(req)
		err := ops.MarkException(req.Context(), opID, body.Scan, body.Reason)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, map[string]string{"ok": "marked"})
	}))

	// Query
	r.POST("/pda/api/query", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Scan string }
		json.NewDecoder(req.Body).Decode(&body)
		token := getToken(req)
		p, logs, err := svc.QueryParcel(req.Context(), token, body.Scan)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 200, map[string]interface{}{"parcel": p, "logs": logs})
	}))

	// Pending
	r.GET("/pda/api/pending", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		apiJSON(w, 200, map[string]interface{}{
			"receive": ops.PendingReceive(ctx), "putaway": ops.PendingPutAway(ctx),
			"weigh": ops.PendingWeigh(ctx), "pick": ops.PendingPick(ctx),
			"pack": ops.PendingPack(ctx),
		})
	}))
}

func apiJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
