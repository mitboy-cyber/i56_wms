package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	domain "github.com/i56/modules/system/domain"
)

// ===================================================================
// MemAPIConfigRepo — in-memory repo for system API integration configs
// ===================================================================
type MemAPIConfigRepo struct {
	mu          sync.RWMutex
	courierAPIs map[int64]*domain.CourierAPIConfig
	customsAPIs map[int64]*domain.CustomsBrokerConfig
	notifChans  map[int64]*domain.NotificationChannel
	printTmpls  map[int64]*domain.PrintTemplate
	storageCfgs map[int64]*domain.StorageConfig
	nextID      int64
}

// NewMemAPIConfigRepo creates a new repo with realistic seed data.
func NewMemAPIConfigRepo() *MemAPIConfigRepo {
	r := &MemAPIConfigRepo{
		courierAPIs: make(map[int64]*domain.CourierAPIConfig),
		customsAPIs: make(map[int64]*domain.CustomsBrokerConfig),
		notifChans:  make(map[int64]*domain.NotificationChannel),
		printTmpls:  make(map[int64]*domain.PrintTemplate),
		storageCfgs: make(map[int64]*domain.StorageConfig),
	}
	r.seed()
	return r
}

func (r *MemAPIConfigRepo) next() int64 { return atomic.AddInt64(&r.nextID, 1) }

// ===================================================================
// Seed data — 3 real-looking configs per module
// ===================================================================
func (r *MemAPIConfigRepo) seed() {
	now := time.Now()

	// --- Courier API Configs (快递公司 API) ---
	r.courierAPIs[1] = &domain.CourierAPIConfig{
		ID: 1, TenantID: 1, CourierID: 1,
		Name:            "顺丰速运API",
		APIEndpoint:     "https://open.sf-express.com/std/service",
		APIKey:          "SF_APP_KEY_001",
		APISecret:       "SF_SECRET_001",
		TrackingPattern: `^SF\d{12}$`,
		AuthType:        "hmac",
		ExtraHeaders:    `{"Accept-Language":"zh-CN"}`,
		IsActive:        true,
		CreatedAt:       now, UpdatedAt: now,
	}
	r.courierAPIs[2] = &domain.CourierAPIConfig{
		ID: 2, TenantID: 1, CourierID: 2,
		Name:            "圆通快递API",
		APIEndpoint:     "https://open.yto.net.cn/api/v1",
		APIKey:          "YTO_APP_KEY_002",
		APISecret:       "YTO_SECRET_002",
		TrackingPattern: `^YT\d{13}$`,
		AuthType:        "api_key",
		ExtraHeaders:    `{}`,
		IsActive:        true,
		CreatedAt:       now, UpdatedAt: now,
	}
	r.courierAPIs[3] = &domain.CourierAPIConfig{
		ID: 3, TenantID: 1, CourierID: 3,
		Name:            "中通快递API",
		APIEndpoint:     "https://open.zto.com/api/order",
		APIKey:          "ZTO_APP_KEY_003",
		APISecret:       "ZTO_SECRET_003",
		TrackingPattern: `^ZTO\d{10,14}$`,
		AuthType:        "api_key",
		ExtraHeaders:    `{}`,
		IsActive:        false,
		CreatedAt:       now, UpdatedAt: now,
	}

	// --- Customs Broker Configs (报关行 API) ---
	r.customsAPIs[1] = &domain.CustomsBrokerConfig{
		ID: 1, TenantID: 1, BrokerID: 1,
		Name:               "厦门电子口岸清关",
		DeclarationAPIURL:  "https://customs.xm-port.gov.cn/api/v2/declaration",
		APIKey:             "XM_CUSTOMS_KEY_001",
		APISecret:          "XM_CUSTOMS_SECRET_001",
		CustomsPointID:     "CN_XM_3701",
		NumberPrefix:       "776XM",
		SupportedDocuments: `["invoice","packing_list","certificate_of_origin","bill_of_lading"]`,
		IsActive:           true,
		CreatedAt:          now, UpdatedAt: now,
	}
	r.customsAPIs[2] = &domain.CustomsBrokerConfig{
		ID: 2, TenantID: 1, BrokerID: 2,
		Name:               "深圳海关申报系统",
		DeclarationAPIURL:  "https://customs.sz-port.gov.cn/api/v1/declaration",
		APIKey:             "SZ_CUSTOMS_KEY_002",
		APISecret:          "SZ_CUSTOMS_SECRET_002",
		CustomsPointID:     "CN_SZ_5302",
		NumberPrefix:       "SZ53",
		SupportedDocuments: `["invoice","packing_list","certificate_of_origin"]`,
		IsActive:           true,
		CreatedAt:          now, UpdatedAt: now,
	}
	r.customsAPIs[3] = &domain.CustomsBrokerConfig{
		ID: 3, TenantID: 1, BrokerID: 3,
		Name:               "上海保税区清关",
		DeclarationAPIURL:  "https://customs.sh-port.gov.cn/api/v2/declaration",
		APIKey:             "SH_CUSTOMS_KEY_003",
		APISecret:          "SH_CUSTOMS_SECRET_003",
		CustomsPointID:     "CN_SH_2201",
		NumberPrefix:       "SH22",
		SupportedDocuments: `["invoice","packing_list","certificate_of_origin","fumigation_cert"]`,
		IsActive:           false,
		CreatedAt:          now, UpdatedAt: now,
	}

	// --- Notification Channels (通知渠道) ---
	r.notifChans[1] = &domain.NotificationChannel{
		ID: 1, TenantID: 1,
		Name:        "阿里云短信服务",
		ChannelType: "sms",
		Provider:    "aliyun_sms",
		ConfigJSON:  `{"access_key_id":"LTAI5tXXXX","access_key_secret":"XXXX","sign_name":"i56物流","template_code":"SMS_123456789","region":"cn-hangzhou"}`,
		IsActive:    true,
		CreatedAt:   now, UpdatedAt: now,
	}
	r.notifChans[2] = &domain.NotificationChannel{
		ID: 2, TenantID: 1,
		Name:        "SendGrid邮件服务",
		ChannelType: "email",
		Provider:    "sendgrid",
		ConfigJSON:  `{"api_key":"SG.xxxxx","from_email":"noreply@i56.com","from_name":"I56物流系统","templates":{"order_shipped":"d-xxx","parcel_received":"d-yyy"}}`,
		IsActive:    true,
		CreatedAt:   now, UpdatedAt: now,
	}
	r.notifChans[3] = &domain.NotificationChannel{
		ID: 3, TenantID: 1,
		Name:        "Line Notify通知",
		ChannelType: "line",
		Provider:    "line_notify",
		ConfigJSON:  `{"access_token":"LN_xxxxx","notify_endpoint":"https://notify-api.line.me/api/notify"}`,
		IsActive:    true,
		CreatedAt:   now, UpdatedAt: now,
	}

	// --- Print Templates (打印模板) ---
	r.printTmpls[1] = &domain.PrintTemplate{
		ID: 1, TenantID: 1,
		Name:            "顺丰标准面单",
		Type:            "label",
		PaperSize:       "100x150mm",
		TemplateContent: `^XA^FO50,50^A0N,30^FD{i56_order_no}^FS^FO50,100^A0N,24^FD{recipient_name}^FS^FO50,140^A0N,24^FD{recipient_phone}^FS^FO50,180^A0N,24^FD{recipient_address}^FS^XZ`,
		PrinterType:     "thermal",
		IsActive:        true,
		CreatedAt:       now, UpdatedAt: now,
	}
	r.printTmpls[2] = &domain.PrintTemplate{
		ID: 2, TenantID: 1,
		Name:            "标准商业发票",
		Type:            "invoice",
		PaperSize:       "A4",
		TemplateContent: `<html><body><h1>Commercial Invoice</h1><table><tr><td>Invoice No:</td><td>{invoice_no}</td></tr><tr><td>Shipper:</td><td>{shipper}</td></tr><tr><td>Consignee:</td><td>{consignee}</td></tr><tr><td>Items:</td><td>{items_table}</td></tr><tr><td>Total:</td><td>{total_amount}</td></tr></table></body></html>`,
		PrinterType:     "laser",
		IsActive:        true,
		CreatedAt:       now, UpdatedAt: now,
	}
	r.printTmpls[3] = &domain.PrintTemplate{
		ID: 3, TenantID: 1,
		Name:            "装箱单模板",
		Type:            "packing_list",
		PaperSize:       "A4",
		TemplateContent: `<html><body><h1>Packing List</h1><table><tr><th>箱号</th><th>品名</th><th>数量</th><th>重量(kg)</th><th>尺寸(cm)</th></tr>{items_rows}</table></body></html>`,
		PrinterType:     "laser",
		IsActive:        true,
		CreatedAt:       now, UpdatedAt: now,
	}

	// --- Storage Configs (对象存储) ---
	r.storageCfgs[1] = &domain.StorageConfig{
		ID: 1, TenantID: 1,
		Name:      "厦门仓MinIO存储",
		Provider:  "minio",
		Bucket:    "i56-xiamen-prod",
		Endpoint:  "https://minio.xm.i56.internal:9000",
		AccessKey: "MINIO_ACCESS_KEY_XM_001",
		SecretKey: "MINIO_SECRET_KEY_XM_001",
		Region:    "cn-xiamen",
		IsActive:  true,
		CreatedAt: now, UpdatedAt: now,
	}
	r.storageCfgs[2] = &domain.StorageConfig{
		ID: 2, TenantID: 1,
		Name:      "阿里云OSS归档存储",
		Provider:  "oss",
		Bucket:    "i56-archive",
		Endpoint:  "https://oss-cn-hangzhou.aliyuncs.com",
		AccessKey: "OSS_ACCESS_KEY_002",
		SecretKey: "OSS_SECRET_KEY_002",
		Region:    "cn-hangzhou",
		IsActive:  true,
		CreatedAt: now, UpdatedAt: now,
	}
	r.storageCfgs[3] = &domain.StorageConfig{
		ID: 3, TenantID: 1,
		Name:      "AWS S3备份存储",
		Provider:  "s3",
		Bucket:    "i56-backup-ap-southeast-1",
		Endpoint:  "https://s3.ap-southeast-1.amazonaws.com",
		AccessKey: "AKIAIOSFODNN7EXAMPLE",
		SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		Region:    "ap-southeast-1",
		IsActive:  false,
		CreatedAt: now, UpdatedAt: now,
	}

	r.nextID = 100
}

// ===================================================================
// Courier API Configs CRUD
// ===================================================================

// ListCouriers returns all courier API configs for a tenant.
func (r *MemAPIConfigRepo) ListCouriers(tenantID int64) []*domain.CourierAPIConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.CourierAPIConfig
	for _, c := range r.courierAPIs {
		if c.TenantID == tenantID {
			res = append(res, c)
		}
	}
	return res
}

// SaveCourier creates or updates a courier API config.
func (r *MemAPIConfigRepo) SaveCourier(_ context.Context, c *domain.CourierAPIConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		c.ID = r.next()
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
	} else {
		c.UpdatedAt = time.Now()
	}
	if c.TenantID == 0 {
		c.TenantID = 1
	}
	r.courierAPIs[c.ID] = c
}

// DeleteCourier removes a courier API config.
func (r *MemAPIConfigRepo) DeleteCourier(_ context.Context, tenantID, id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.courierAPIs[id]; ok && c.TenantID == tenantID {
		delete(r.courierAPIs, id)
	}
}

// ===================================================================
// Customs Broker Configs CRUD
// ===================================================================

// ListCustomsBrokers returns all customs broker configs for a tenant.
func (r *MemAPIConfigRepo) ListCustomsBrokers(tenantID int64) []*domain.CustomsBrokerConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.CustomsBrokerConfig
	for _, c := range r.customsAPIs {
		if c.TenantID == tenantID {
			res = append(res, c)
		}
	}
	return res
}

// SaveCustomsBroker creates or updates a customs broker config.
func (r *MemAPIConfigRepo) SaveCustomsBroker(_ context.Context, c *domain.CustomsBrokerConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		c.ID = r.next()
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
	} else {
		c.UpdatedAt = time.Now()
	}
	if c.TenantID == 0 {
		c.TenantID = 1
	}
	r.customsAPIs[c.ID] = c
}

// DeleteCustomsBroker removes a customs broker config.
func (r *MemAPIConfigRepo) DeleteCustomsBroker(_ context.Context, tenantID, id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.customsAPIs[id]; ok && c.TenantID == tenantID {
		delete(r.customsAPIs, id)
	}
}

// ===================================================================
// Notification Channels CRUD (richer version)
// ===================================================================

// ListNotificationChannels returns all notification channels for a tenant.
func (r *MemAPIConfigRepo) ListNotificationChannels(tenantID int64) []*domain.NotificationChannel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.NotificationChannel
	for _, c := range r.notifChans {
		if c.TenantID == tenantID {
			res = append(res, c)
		}
	}
	return res
}

// SaveNotificationChannel creates or updates a notification channel.
func (r *MemAPIConfigRepo) SaveNotificationChannel(_ context.Context, c *domain.NotificationChannel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		c.ID = r.next()
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
	} else {
		c.UpdatedAt = time.Now()
	}
	if c.TenantID == 0 {
		c.TenantID = 1
	}
	r.notifChans[c.ID] = c
}

// DeleteNotificationChannel removes a notification channel.
func (r *MemAPIConfigRepo) DeleteNotificationChannel(_ context.Context, tenantID, id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.notifChans[id]; ok && c.TenantID == tenantID {
		delete(r.notifChans, id)
	}
}

// ===================================================================
// Print Templates CRUD
// ===================================================================

// ListPrintTemplates returns all print templates for a tenant.
func (r *MemAPIConfigRepo) ListPrintTemplates(tenantID int64) []*domain.PrintTemplate {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.PrintTemplate
	for _, p := range r.printTmpls {
		if p.TenantID == tenantID {
			res = append(res, p)
		}
	}
	return res
}

// SavePrintTemplate creates or updates a print template.
func (r *MemAPIConfigRepo) SavePrintTemplate(_ context.Context, c *domain.PrintTemplate) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		c.ID = r.next()
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
	} else {
		c.UpdatedAt = time.Now()
	}
	if c.TenantID == 0 {
		c.TenantID = 1
	}
	r.printTmpls[c.ID] = c
}

// DeletePrintTemplate removes a print template.
func (r *MemAPIConfigRepo) DeletePrintTemplate(_ context.Context, tenantID, id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.printTmpls[id]; ok && c.TenantID == tenantID {
		delete(r.printTmpls, id)
	}
}

// ===================================================================
// Storage Configs CRUD
// ===================================================================

// ListStorageConfigs returns all storage configs for a tenant.
func (r *MemAPIConfigRepo) ListStorageConfigs(tenantID int64) []*domain.StorageConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.StorageConfig
	for _, s := range r.storageCfgs {
		if s.TenantID == tenantID {
			res = append(res, s)
		}
	}
	return res
}

// SaveStorageConfig creates or updates a storage config.
func (r *MemAPIConfigRepo) SaveStorageConfig(_ context.Context, c *domain.StorageConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		c.ID = r.next()
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
	} else {
		c.UpdatedAt = time.Now()
	}
	if c.TenantID == 0 {
		c.TenantID = 1
	}
	r.storageCfgs[c.ID] = c
}

// DeleteStorageConfig removes a storage config.
func (r *MemAPIConfigRepo) DeleteStorageConfig(_ context.Context, tenantID, id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.storageCfgs[id]; ok && c.TenantID == tenantID {
		delete(r.storageCfgs, id)
	}
}
