package main
import (
	"encoding/json"; "fmt"; "net/http"; "strconv"; "time"
	"github.com/i56/framework/core/response"
	"github.com/i56/framework/core/router"
	pdaSvc "github.com/i56/modules/pda/service"
	tdRepo "github.com/i56/modules/taskdispatch/repository"
)

func registerPDAAPIRoutesOnRouter(r *router.Router, ops *pdaSvc.PDAOperations) {
	r.GET("/api/v1/pda/logs",func(w http.ResponseWriter,req *http.Request){response.JSON(w,200,ops.RecentLogs(20))})
	r.GET("/api/v1/pda/warehouse-stats",func(w http.ResponseWriter,req *http.Request){response.JSON(w,200,ops.WarehouseStats(req.Context()))})
	r.GET("/api/v1/pda/pending-receive",func(w http.ResponseWriter,req *http.Request){response.JSON(w,200,ops.PendingReceive(req.Context()))})
	r.GET("/api/v1/pda/pending-putaway",func(w http.ResponseWriter,req *http.Request){response.JSON(w,200,ops.PendingPutAway(req.Context()))})
	r.GET("/api/v1/pda/pending-pick",func(w http.ResponseWriter,req *http.Request){response.JSON(w,200,ops.PendingPick(req.Context()))})
	r.POST("/api/v1/pda/receive",func(w http.ResponseWriter,req *http.Request){
		var b struct{TrackingNo string `json:"tracking_no"`;Weight float64 `json:"weight"`;Length float64 `json:"length"`;Width float64 `json:"width"`;Height float64 `json:"height"`;Location string `json:"location_barcode"`;OpID int64 `json:"operator_id"`}
		json.NewDecoder(req.Body).Decode(&b);if b.OpID==0{b.OpID=1}
		p,err:=ops.Receive(req.Context(),b.OpID,b.TrackingNo,b.Weight,b.Length,b.Width,b.Height,b.Location)
		if err!=nil{response.Error(w,err);return};response.JSON(w,200,map[string]any{"parcel":p,"message":"入库成功"})
	})
	r.POST("/api/v1/pda/putaway",func(w http.ResponseWriter,req *http.Request){
		var b struct{TrackingNo string `json:"tracking_no"`;Location string `json:"location_barcode"`;OpID int64 `json:"operator_id"`}
		json.NewDecoder(req.Body).Decode(&b);if b.OpID==0{b.OpID=1}
		p,err:=ops.PutAway(req.Context(),b.OpID,b.TrackingNo,b.Location)
		if err!=nil{response.Error(w,err);return};response.JSON(w,200,map[string]any{"parcel":p,"message":"上架成功"})
	})
	r.POST("/api/v1/pda/pick",func(w http.ResponseWriter,req *http.Request){
		var b struct{OrderNo string `json:"order_no"`;OpID int64 `json:"operator_id"`}
		json.NewDecoder(req.Body).Decode(&b);if b.OpID==0{b.OpID=1}
		o,parcels,err:=ops.Pick(req.Context(),b.OpID,b.OrderNo)
		if err!=nil{response.Error(w,err);return};response.JSON(w,200,map[string]any{"order":o,"parcels":parcels,"message":"拣货完成"})
	})
	r.POST("/api/v1/pda/pack",func(w http.ResponseWriter,req *http.Request){
		var b struct{OrderNo string `json:"order_no"`;OpID int64 `json:"operator_id"`}
		json.NewDecoder(req.Body).Decode(&b);if b.OpID==0{b.OpID=1}
		o,err:=ops.Pack(req.Context(),b.OpID,b.OrderNo)
		if err!=nil{response.Error(w,err);return};response.JSON(w,200,map[string]any{"order":o,"message":"打包完成"})
	})
	r.POST("/api/v1/pda/load",func(w http.ResponseWriter,req *http.Request){
		var b struct{ContainerNo string `json:"container_no"`;OrderNo string `json:"order_no"`;OpID int64 `json:"operator_id"`}
		json.NewDecoder(req.Body).Decode(&b);if b.OpID==0{b.OpID=1}
		err:=ops.LoadContainer(req.Context(),b.OpID,b.ContainerNo,b.OrderNo)
		if err!=nil{response.Error(w,err);return};response.JSON(w,200,map[string]any{"message":"装柜完成"})
	})
	r.POST("/api/v1/pda/exception",func(w http.ResponseWriter,req *http.Request){
		var b struct{TrackingNo string `json:"tracking_no"`;Reason string `json:"reason"`;OpID int64 `json:"operator_id"`}
		json.NewDecoder(req.Body).Decode(&b);if b.OpID==0{b.OpID=1}
		err:=ops.MarkException(req.Context(),b.OpID,b.TrackingNo,b.Reason)
		if err!=nil{response.Error(w,err);return};response.JSON(w,200,map[string]any{"message":"已标记异常"})
	})
}

// ──────────────────────────────────────────────────────────────────────
// Task Dispatch Engine — 抢单池 API
// ──────────────────────────────────────────────────────────────────────

func registerTaskDispatchRoutes(r *router.Router, td *tdRepo.MemTaskDispatchRepo) {
	// GET /api/v1/pda/task-pool — list all pending tasks (抢单池)
	r.GET("/api/v1/pda/task-pool", func(w http.ResponseWriter, req *http.Request) {
		tasks := td.TaskPool()
		response.JSON(w, 200, map[string]any{
			"count": len(tasks),
			"tasks": tasks,
		})
	})

	// POST /api/v1/pda/tasks/{id}/claim — claim a task from the pool
	r.POST("/api/v1/pda/tasks/{id}/claim", func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		taskID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.JSON(w, 400, map[string]any{"error": "invalid task id"})
			return
		}
		var b struct{ OperatorID int64 `json:"operator_id"` }
		json.NewDecoder(req.Body).Decode(&b)
		if b.OperatorID == 0 { b.OperatorID = 1 }
		task, err := td.ClaimTask(taskID, b.OperatorID)
		if err != nil {
			response.Error(w, err)
			return
		}
		response.JSON(w, 200, map[string]any{
			"task":    task,
			"message": fmt.Sprintf("任务 %s 已认领", task.TaskCode),
		})
	})

	// POST /api/v1/pda/tasks/{id}/start — start a claimed task
	r.POST("/api/v1/pda/tasks/{id}/start", func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		taskID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.JSON(w, 400, map[string]any{"error": "invalid task id"})
			return
		}
		task, err := td.StartTask(taskID)
		if err != nil {
			response.Error(w, err)
			return
		}
		response.JSON(w, 200, map[string]any{
			"task":    task,
			"message": fmt.Sprintf("任务 %s 已开始", task.TaskCode),
		})
	})

	// POST /api/v1/pda/tasks/{id}/complete — complete a task
	r.POST("/api/v1/pda/tasks/{id}/complete", func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		taskID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.JSON(w, 400, map[string]any{"error": "invalid task id"})
			return
		}
		task, err := td.CompleteTask(taskID)
		if err != nil {
			response.Error(w, err)
			return
		}
		response.JSON(w, 200, map[string]any{
			"task":    task,
			"message": fmt.Sprintf("任务 %s 已完成", task.TaskCode),
		})
	})

	// GET /api/v1/pda/my-tasks — operator's current tasks (claimed + in_progress)
	r.GET("/api/v1/pda/my-tasks", func(w http.ResponseWriter, req *http.Request) {
		opIDStr := req.URL.Query().Get("operator_id")
		opID, err := strconv.ParseInt(opIDStr, 10, 64)
		if err != nil || opID == 0 { opID = 1 }
		tasks := td.GetOperatorTasks(opID)
		response.JSON(w, 200, map[string]any{
			"operator_id": opID,
			"count":       len(tasks),
			"tasks":       tasks,
		})
	})

	// GET /api/v1/pda/tasks/{id} — get task details
	r.GET("/api/v1/pda/tasks/{id}", func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		taskID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.JSON(w, 400, map[string]any{"error": "invalid task id"})
			return
		}
		task, err := td.GetTaskByID(taskID)
		if err != nil {
			response.Error(w, err)
			return
		}
		response.JSON(w, 200, task)
	})

	// GET /api/v1/pda/operators — list all operators with capabilities
	r.GET("/api/v1/pda/operators", func(w http.ResponseWriter, req *http.Request) {
		ops := td.GetOperators()
		response.JSON(w, 200, map[string]any{
			"count":     len(ops),
			"operators": ops,
		})
	})

	// GET /api/v1/pda/matching/{id} — find best operator match for a task
	r.GET("/api/v1/pda/matching/{id}", func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		taskID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.JSON(w, 400, map[string]any{"error": "invalid task id"})
			return
		}
		task, err := td.GetTaskByID(taskID)
		if err != nil {
			response.Error(w, err)
			return
		}
		match := td.MatchOperator(task)
		if match == nil {
			response.JSON(w, 200, map[string]any{"task_id": taskID, "match": nil, "message": "无可用操作员"})
			return
		}
		response.JSON(w, 200, map[string]any{
			"task_id": taskID,
			"match":   match,
			"score":   match.MatchScore(task.RequiredCapabilities),
		})
	})

	// POST /api/v1/pda/check-timeouts — manually trigger timeout check
	r.POST("/api/v1/pda/check-timeouts", func(w http.ResponseWriter, req *http.Request) {
		timedOut := td.CheckTimeouts()
		response.JSON(w, 200, map[string]any{
			"timed_out": len(timedOut),
			"tasks":     timedOut,
			"message":   fmt.Sprintf("%d 个超时任务已重新分配", len(timedOut)),
		})
	})
}

// StartTimeoutChecker runs a background goroutine that checks for timeouts every 30 seconds.
func StartTimeoutChecker(td *tdRepo.MemTaskDispatchRepo) chan struct{} {
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				td.CheckTimeouts()
			case <-stop:
				return
			}
		}
	}()
	return stop
}
