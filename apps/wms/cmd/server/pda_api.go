package main
import ("encoding/json";"fmt";"net/http";"github.com/i56/framework/core/response";"github.com/i56/framework/core/router";"github.com/i56/framework/core/sse";pdaSvc "github.com/i56/modules/pda/service")

func registerPDAAPIRoutes(api *router.Router, ops *pdaSvc.PDAOperations, hub *sse.Hub) {
	api.GET("/pda/logs", func(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, ops.RecentLogs(20)) })
	api.GET("/pda/warehouse-stats", func(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, ops.WarehouseStats(r.Context())) })
	api.GET("/pda/pending-receive", func(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, ops.PendingReceive(r.Context())) })
	api.GET("/pda/pending-putaway", func(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, ops.PendingPutAway(r.Context())) })
	api.GET("/pda/pending-pick", func(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, ops.PendingPick(r.Context())) })

	api.POST("/pda/receive", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ TrackingNo string `json:"tracking_no"`; Weight,Length,Width,Height float64; OpID int64 `json:"operator_id"` }
		json.NewDecoder(r.Body).Decode(&req); if req.OpID==0{req.OpID=1}
		p,err:=ops.Receive(r.Context(),req.OpID,req.TrackingNo,req.Weight,req.Length,req.Width,req.Height,"")
		if err!=nil{response.Error(w,err);return}
		hub.Publish("inbound",sse.Event{Type:"parcel_received",Data:fmt.Sprintf(`{"tracking_no":"%s"}`,req.TrackingNo)})
		response.JSON(w,200,map[string]any{"parcel":p,"message":"收货成功"})
	})
	api.POST("/pda/putaway", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ TrackingNo,LocationBarcode string `json:"tracking_no" json:"location_barcode"`; OpID int64 `json:"operator_id"` }
		json.NewDecoder(r.Body).Decode(&req); if req.OpID==0{req.OpID=1}
		pr,err:=ops.PutAway(r.Context(),req.OpID,req.TrackingNo,req.LocationBarcode)
		if err!=nil{response.Error(w,err);return}
		hub.Publish("inbound",sse.Event{Type:"parcel_stored",Data:fmt.Sprintf(`{"tracking_no":"%s"}`,req.TrackingNo)})
		response.JSON(w,200,map[string]any{"parcel":pr,"message":"上架成功"})
	})
	api.POST("/pda/pick", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ OrderNo string `json:"order_no"`; OpID int64 `json:"operator_id"` }
		json.NewDecoder(r.Body).Decode(&req); if req.OpID==0{req.OpID=1}
		o,parcels,err:=ops.Pick(r.Context(),req.OpID,req.OrderNo)
		if err!=nil{response.Error(w,err);return}
		hub.Publish("inbound",sse.Event{Type:"order_picked",Data:fmt.Sprintf(`{"order_no":"%s","parcels":%d}`,req.OrderNo,len(parcels))})
		response.JSON(w,200,map[string]any{"order":o,"message":"拣货完成"})
	})
	api.POST("/pda/pack", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ OrderNo string `json:"order_no"`; OpID int64 `json:"operator_id"` }
		json.NewDecoder(r.Body).Decode(&req); if req.OpID==0{req.OpID=1}
		o,err:=ops.Pack(r.Context(),req.OpID,req.OrderNo)
		if err!=nil{response.Error(w,err);return}
		hub.Publish("inbound",sse.Event{Type:"order_packed",Data:fmt.Sprintf(`{"order_no":"%s"}`,req.OrderNo)})
		response.JSON(w,200,map[string]any{"order":o,"message":"打包完成"})
	})
	api.POST("/pda/load", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ ContainerNo,OrderNo string `json:"container_no" json:"order_no"`; OpID int64 `json:"operator_id"` }
		json.NewDecoder(r.Body).Decode(&req); if req.OpID==0{req.OpID=1}
		err:=ops.LoadContainer(r.Context(),req.OpID,req.ContainerNo,req.OrderNo)
		if err!=nil{response.Error(w,err);return}
		hub.Publish("inbound",sse.Event{Type:"container_loaded",Data:fmt.Sprintf(`{"container":"%s"}`,req.ContainerNo)})
		response.JSON(w,200,map[string]any{"message":"装柜完成"})
	})
	api.POST("/pda/exception", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ TrackingNo,Reason string `json:"tracking_no" json:"reason"`; OpID int64 `json:"operator_id"` }
		json.NewDecoder(r.Body).Decode(&req); if req.OpID==0{req.OpID=1}
		err:=ops.MarkException(r.Context(),req.OpID,req.TrackingNo,req.Reason)
		if err!=nil{response.Error(w,err);return}
		hub.Publish("inbound",sse.Event{Type:"parcel_abnormal",Data:fmt.Sprintf(`{"tracking_no":"%s"}`,req.TrackingNo)})
		response.JSON(w,200,map[string]any{"message":"标记异常"})
	})
	api.POST("/pda/send-outbound", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ OrderNo string `json:"order_no"`; OpID int64 `json:"operator_id"` }
		json.NewDecoder(r.Body).Decode(&req); if req.OpID==0{req.OpID=1}
		err:=ops.SendToOutbound(r.Context(),req.OpID,req.OrderNo)
		if err!=nil{response.Error(w,err);return}
		response.JSON(w,200,map[string]any{"message":"已送出货区"})
	})
}
