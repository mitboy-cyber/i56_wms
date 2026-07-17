package repository
import (
	"crypto/rand"; "encoding/hex"; "sync"; "sync/atomic"; "time"
	"github.com/i56/modules/pda/domain"
)
type MemPDARepo struct {
	mu        sync.RWMutex
	operators map[int64]*domain.Operator
	sessions  map[string]*domain.Session
	scanLogs  []domain.ScanLog
	nextOpID  int64
}
func NewMemPDARepo() *MemPDARepo {
	r:=&MemPDARepo{operators:make(map[int64]*domain.Operator),sessions:make(map[string]*domain.Session),nextOpID:1}
	r.seedOperators()
	return r
}
func (r *MemPDARepo) seedOperators() {
	for _,op:=range []struct{name,code,pin string; whID int64}{{"张操作员","OP001","1234",1},{"李操作员","OP002","1234",1},{"王操作员","OP003","1234",1}}{
		id:=atomic.AddInt64(&r.nextOpID,1)-1
		r.operators[id]=&domain.Operator{ID:id,TenantID:1,WarehouseID:op.whID,Name:op.name,Code:op.code,Pin:op.pin,IsActive:true}
	}
}
func (r *MemPDARepo) GetOperatorByID(id int64)*domain.Operator{
	r.mu.RLock();defer r.mu.RUnlock()
	op,ok:=r.operators[id];if ok&&op.IsActive{return op};return nil
}
func (r *MemPDARepo) GetOperatorByCode(code string)*domain.Operator{
	r.mu.RLock();defer r.mu.RUnlock()
	for _,op:=range r.operators{if op.Code==code&&op.IsActive{return op}};return nil
}
func (r *MemPDARepo) VerifyPin(opID int64,pin string)bool{
	r.mu.RLock();defer r.mu.RUnlock()
	op,ok:=r.operators[opID];return ok&&op.Pin==pin
}
func (r *MemPDARepo) CreateSession(opID int64,deviceID string)*domain.Session{
	r.mu.Lock();defer r.mu.Unlock()
	tok:=make([]byte,16);rand.Read(tok)
	s:=&domain.Session{OperatorID:opID,Token:hex.EncodeToString(tok),DeviceID:deviceID,IsActive:true,LoginAt:time.Now(),LastSeen:time.Now()}
	r.sessions[s.Token]=s;return s
}
func (r *MemPDARepo) ValidateSession(token string)*domain.Session{
	r.mu.RLock();defer r.mu.RUnlock()
	s,ok:=r.sessions[token]
	if !ok||!s.IsActive{return nil}
	s.LastSeen=time.Now();return s
}
func (r *MemPDARepo) Logout(token string){
	r.mu.Lock();defer r.mu.Unlock()
	if s,ok:=r.sessions[token];ok{s.IsActive=false}
}
func (r *MemPDARepo) LogScan(log *domain.ScanLog){
	r.mu.Lock();defer r.mu.Unlock()
	log.ScannedAt=time.Now();r.scanLogs=append(r.scanLogs,*log)
}
func (r *MemPDARepo) RecentScans(limit int)[]domain.ScanLog{
	r.mu.RLock();defer r.mu.RUnlock()
	if limit>len(r.scanLogs){limit=len(r.scanLogs)}
	return r.scanLogs[len(r.scanLogs)-limit:]
}
