package service
import (
	"context"; "fmt"; "time"
	"github.com/i56/modules/pda/domain"
	"github.com/i56/modules/pda/repository"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelRepo "github.com/i56/modules/parcel/repository"
	orderRepo "github.com/i56/modules/order/repository"
)
type PDAService struct {
	pdaRepo *repository.MemPDARepo
	parcelRepo *parcelRepo.MemParcelRepo
	orderRepo *orderRepo.MemOrderRepo
}
func NewPDAService(pdaR *repository.MemPDARepo,parcelR *parcelRepo.MemParcelRepo,orderR *orderRepo.MemOrderRepo)*PDAService{
	return &PDAService{pdaRepo:pdaR,parcelRepo:parcelR,orderRepo:orderR}
}
func (s *PDAService) Login(code,pin,deviceID string)(*domain.Session,error){
	op:=s.pdaRepo.GetOperatorByCode(code)
	if op==nil{return nil,fmt.Errorf("操作员不存在")}
	if !s.pdaRepo.VerifyPin(op.ID,pin){return nil,fmt.Errorf("PIN码错误")}
	sess:=s.pdaRepo.CreateSession(op.ID,deviceID)
	s.pdaRepo.LogScan(&domain.ScanLog{TenantID:1,WarehouseID:op.WarehouseID,OperatorID:op.ID,Action:"login",Barcode:code,Success:true,Message:"登录成功"})
	return sess,nil
}
func (s *PDAService) ReceiveParcel(ctx context.Context,sessionToken,trackingNo string,weight,length,width,height float64)(*parcelDomain.Parcel,error){
	sess:=s.pdaRepo.ValidateSession(sessionToken)
	if sess==nil{return nil,fmt.Errorf("会话已过期")}
	op:=s.pdaRepo.GetOperatorByID(sess.OperatorID)
	p,err:=s.parcelRepo.GetByTrackingNo(ctx,1,trackingNo)
	if err!=nil||p==nil{return nil,fmt.Errorf("包裹不存在: %s",trackingNo)}
	if p.Status!=parcelDomain.StatusPreDeclared{return nil,fmt.Errorf("包裹状态为%s，无法入库",p.Status)}
	if !p.CanTransitionTo(parcelDomain.StatusReceived){return nil,fmt.Errorf("无法从%s转换到received",p.Status)}
	p.ActualWeight=weight;p.Length=length;p.Width=width;p.Height=height
	p.Status=parcelDomain.StatusReceived;p.UpdatedAt=time.Now()
	if err:=s.parcelRepo.Update(ctx,p);err!=nil{return nil,err}
	s.pdaRepo.LogScan(&domain.ScanLog{TenantID:1,WarehouseID:op.WarehouseID,OperatorID:op.ID,Action:"receive",TrackingNumber:trackingNo,Weight:weight,Success:true,Message:"入库成功"})
	return p,nil
}
func (s *PDAService) QuickReceive(ctx context.Context,sessionToken,trackingNo string)(*parcelDomain.Parcel,error){
	return s.ReceiveParcel(ctx,sessionToken,trackingNo,0.5,20,15,10)
}
func (s *PDAService) QueryParcel(ctx context.Context,sessionToken,trackingNo string)(*parcelDomain.Parcel,[]domain.ScanLog,error){
	sess:=s.pdaRepo.ValidateSession(sessionToken)
	if sess==nil{return nil,nil,fmt.Errorf("会话已过期")}
	p,err:=s.parcelRepo.GetByTrackingNo(ctx,1,trackingNo)
	if err!=nil||p==nil{return nil,nil,fmt.Errorf("包裹未找到")}
	logs:=s.pdaRepo.RecentScans(10)
	return p,logs,nil
}
func (s *PDAService) RecentScans(sessionToken string,limit int)[]domain.ScanLog{
	if s.pdaRepo.ValidateSession(sessionToken)==nil{return nil}
	return s.pdaRepo.RecentScans(limit)
}
func (s *PDAService) Menus()[]domain.PDAMenu{return domain.DefaultMenus()}
