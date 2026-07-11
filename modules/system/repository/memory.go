package repository
import (
	"context";"sync";"sync/atomic";"time"
	domain "github.com/i56/modules/system/domain"
)
type MemSystemConfigRepo struct {
	mu sync.RWMutex
	logisticsAPIs map[int64]*domain.LogisticsAPIConfig
	brokers map[int64]*domain.CustomsBrokerAPIConfig
	printers map[int64]*domain.PrinterConfig
	channels map[int64]*domain.NotificationChannel
	settings map[int64]*domain.SystemSetting
	apiConfigs []APIConfigEntry
	nextID int64
}
func NewMemSystemConfigRepo() *MemSystemConfigRepo {
	r := &MemSystemConfigRepo{logisticsAPIs:map[int64]*domain.LogisticsAPIConfig{},brokers:map[int64]*domain.CustomsBrokerAPIConfig{},printers:map[int64]*domain.PrinterConfig{},channels:map[int64]*domain.NotificationChannel{},settings:map[int64]*domain.SystemSetting{}}
	r.seed();return r
}
func (r *MemSystemConfigRepo) next() int64 {return atomic.AddInt64(&r.nextID,1)}
func (r *MemSystemConfigRepo) seed() {
	r.logisticsAPIs[1]=&domain.LogisticsAPIConfig{ID:1,TenantID:1,CarrierID:1,Name:"顺丰速运API",BaseURL:"https://sfapi.sf-express.com/std/service",AppKey:"SF_KEY_001",WebhookURL:"https://wms.mikaplay.com/api/v1/webhook/sf",TimeoutMs:5000,MaxRetry:3,IsActive:true}
	r.nextID=10
	r.brokers[1]=&domain.CustomsBrokerAPIConfig{ID:1,TenantID:1,BrokerID:1,Name:"厦门清关",Code:"XM-CS",Prefix:"776XM",APIEndpoint:"https://customs.xm-port.gov.cn/api/v2",SupportedDocs:"invoice,packing_list,certificate_of_origin",IsActive:true}
	r.printers[1]=&domain.PrinterConfig{ID:1,TenantID:1,Name:"入库口打印機",PrinterType:"thermal",PaperWidth:100,PaperHeight:150,DPI:203,ConnectionType:"network",IPAddress:"192.168.1.100",Port:9100,IsDefault:true,IsActive:true}
	r.printers[2]=&domain.PrinterConfig{ID:2,TenantID:1,Name:"出库口打印機",PrinterType:"thermal",PaperWidth:100,PaperHeight:150,DPI:300,ConnectionType:"usb",IsActive:true}
	r.channels[1]=&domain.NotificationChannel{ID:1,TenantID:1,ChannelType:"email",Name:"系统邮件",Provider:"smtp",ConfigJSON:`{"smtp_host":"smtp.example.com","smtp_port":587,"from":"noreply@i56.com"}`,IsActive:true}
	r.channels[2]=&domain.NotificationChannel{ID:2,TenantID:1,ChannelType:"sms",Name:"短信通知",Provider:"aliyun_sms",ConfigJSON:`{"gateway_url":"https://sms-api.example.com/send"}`,IsActive:true}
	r.channels[3]=&domain.NotificationChannel{ID:3,TenantID:1,ChannelType:"webhook",Name:"业务回调",Provider:"webhook",ConfigJSON:`{"url":"https://client.example.com/api/callback","events":["parcel.received","order.shipped"]}`,IsActive:false}
	r.settings[1]=&domain.SystemSetting{ID:1,TenantID:1,Key:"warehouse.default_storage_days",Value:"365",Type:"int",Group:"warehouse",Label:"默认免费仓储天数"}
	r.settings[2]=&domain.SystemSetting{ID:2,TenantID:1,Key:"warehouse.storage_fee_per_day",Value:"1.00",Type:"string",Group:"warehouse",Label:"超期仓储费(元/天)"}
	r.settings[3]=&domain.SystemSetting{ID:3,TenantID:1,Key:"order.auto_cancel_hours",Value:"72",Type:"int",Group:"order",Label:"未付款自动取消(小时)"}
	r.settings[4]=&domain.SystemSetting{ID:4,TenantID:1,Key:"pda.auto_refresh_seconds",Value:"30",Type:"int",Group:"pda",Label:"PDA自动刷新间隔(秒)"}
	r.settings[5]=&domain.SystemSetting{ID:5,TenantID:1,Key:"api.rate_limit_per_minute",Value:"120",Type:"int",Group:"api",Label:"API每秒限流"}
	r.settings[6]=&domain.SystemSetting{ID:6,TenantID:1,Key:"session.idle_timeout_minutes",Value:"30",Type:"int",Group:"security",Label:"会话空闲超时(分)"}
}
func (r *MemSystemConfigRepo) ListLogisticsAPIs(tenantID int64)([]*domain.LogisticsAPIConfig,int64){r.mu.RLock();defer r.mu.RUnlock();var res []*domain.LogisticsAPIConfig;for _,c:=range r.logisticsAPIs{if c.TenantID==tenantID{res=append(res,c)}};return res,int64(len(res))}
func (r *MemSystemConfigRepo) ListBrokers(_ context.Context,tenantID int64) []*domain.CustomsBrokerAPIConfig {r.mu.RLock();defer r.mu.RUnlock();var res []*domain.CustomsBrokerAPIConfig;for _,b:=range r.brokers{if b.TenantID==tenantID{res=append(res,b)}};return res}
func (r *MemSystemConfigRepo) ListPrinters(tenantID int64) []*domain.PrinterConfig {r.mu.RLock();defer r.mu.RUnlock();var res []*domain.PrinterConfig;for _,p:=range r.printers{if p.TenantID==tenantID{res=append(res,p)}};return res}
func (r *MemSystemConfigRepo) ListChannels(_ context.Context,tenantID int64) []*domain.NotificationChannel {r.mu.RLock();defer r.mu.RUnlock();var res []*domain.NotificationChannel;for _,c:=range r.channels{if c.TenantID==tenantID{res=append(res,c)}};return res}
func (r *MemSystemConfigRepo) ListSettings(tenantID int64) []*domain.SystemSetting {r.mu.RLock();defer r.mu.RUnlock();var res []*domain.SystemSetting;for _,s:=range r.settings{if s.TenantID==tenantID{res=append(res,s)}};return res}

type APIConfigEntry struct {
	Name        string
	Provider    string
	Endpoint    string
	APIKey      string
	APISecret   string
	WebhookURL  string
	Description string
	Status      string
}

func (r *MemSystemConfigRepo) ListAPIConfig() []APIConfigEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.apiConfigs
}

func (r *MemSystemConfigRepo) SaveAPIConfig(name, provider, endpoint, apiKey, apiSecret, webhookURL, description string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.apiConfigs = append(r.apiConfigs, APIConfigEntry{
		Name: name, Provider: provider, Endpoint: endpoint,
		APIKey: apiKey, APISecret: apiSecret, WebhookURL: webhookURL,
		Description: description, Status: "已配置",
	})
}

// SaveNotificationChannel adds a new notification channel.
func (r *MemSystemConfigRepo) SaveNotificationChannel(channelType, name, config string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.next()
	r.channels[id] = &domain.NotificationChannel{ID: id, TenantID: 1, ChannelType: channelType, Name: name, Provider: "manual", ConfigJSON: config, IsActive: true}
}

// DeleteChannel removes a notification channel.
func (r *MemSystemConfigRepo) DeleteChannel(_ context.Context, tenantID, id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.channels[id]; ok && c.TenantID == tenantID {
		delete(r.channels, id)
	}
}

// SaveSetting creates or updates a system setting.
func (r *MemSystemConfigRepo) SaveSetting(tenantID int64, key, value, typ, group, label string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check if key already exists and update it
	for _, s := range r.settings {
		if s.TenantID == tenantID && s.Key == key {
			s.Value = value
			s.Type = typ
			s.Group = group
			s.Label = label
			s.UpdatedAt = time.Now()
			return
		}
	}
	// Create new
	id := r.next()
	r.settings[id] = &domain.SystemSetting{
		ID: id, TenantID: tenantID,
		Key: key, Value: value, Type: typ,
		Group: group, Label: label,
	}
}
