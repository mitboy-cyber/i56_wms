package main
import (
	"net/http"
	"github.com/i56/framework/core/auth"; "github.com/i56/framework/core/router"
	custSvc "github.com/i56/modules/customer/service"
	parcelSvc "github.com/i56/modules/parcel/service"
	orderSvc "github.com/i56/modules/order/service"
	whSvc "github.com/i56/modules/warehouse/service"
	pdaSvc "github.com/i56/modules/pda/service"
	// repos for printing/webhook/report
)

func registerAPIRoutes(api *router.Router, tm *auth.TokenManager, ps *parcelSvc.ParcelService, osvc *orderSvc.OrderService, cs *custSvc.ClientService, ws *whSvc.WarehouseService, rr interface{}, cour interface{}, sr interface{}, wor interface{}, ppr interface{}, whr interface{}, rpt interface{}, pdaOps *pdaSvc.PDAOperations) {
	// Minimal API — existing API routes work via the old registration in main.go
	_=ps;_=osvc;_=cs;_=ws;_=rr;_=cour;_=sr;_=wor;_=ppr;_=whr;_=rpt;_=pdaOps;_=tm
	// Health check
	api.GET("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json")
		w.Write([]byte(`{"data":{"name":"I56 Framework","version":"1.1.0","status":"ok","deps":"in-memory"}}`))
	})
}
