# I56 Device Gateway — 设备网关集成指南

> 版本: 1.0.0 | 日期: 2026-07-12 | I56 WMS 硬件集成服务

---

## 1. Architecture Overview (架构概览)

```
┌─────────────────────────────────────────────────────────────┐
│                      I56 WMS Server                         │
│                     Port :8080 (HTTP)                       │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │         Device Gateway API Endpoints                  │   │
│  │  GET  /api/device/inbound-task?barcode=xxx           │   │
│  │  POST /api/device/weight-record                      │   │
│  │  POST /api/device/inbound-confirm                    │   │
│  │  POST /api/device/heartbeat                          │   │
│  └─────────────────────────────────────────────────────┘   │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP (JSON)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              I56 Device Gateway Service                      │
│                   Port :9100 (HTTP)                          │
│                  Port :9101 (WebSocket)                      │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │Dispatcher│  │ Session  │  │ Config   │                  │
│  │  (任务分派)│  │ Manager  │  │ Loader   │                  │
│  └────┬─────┘  └──────────┘  └──────────┘                  │
│       │                                                     │
│  ┌────┴──────────────────────────────┐                     │
│  │      Protocol Adapters            │                     │
│  │  MODBUS_RTU | CONTINUOUS          │                     │
│  │  TOLEDO     | CUSTOM              │                     │
│  └────┬──────────────────────────────┘                     │
└───────┼──────────────────────────────────────────────────────┘
        │ RS-232 / RS-485 Serial
        ▼
┌─────────────────────────────────────────────────────────────┐
│                    Hardware Devices                         │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │  地磅     │  │ 入库机    │  │ 扫码枪    │                  │
│  │  Scale   │  │Conveyor  │  │ Scanner  │                  │
│  │RS-232/485│  │ RS-232   │  │ USB/232  │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow (数据流)

```
扫描条码 → 查询WMS入库任务 → 称重 → 记录重量 → 分拨 → 确认入库
   ↓           ↓              ↓        ↓        ↓        ↓
 Scanner   WMS API         Scale    WMS API  Conveyor  WMS API
```

---

## 2. Supported Protocols (支持的协议)

### 2.1 MODBUS_RTU (标准工业协议)

- **物理层**: RS-485 两线制 (A/B) 或 RS-232
- **波特率**: 9600 / 19200 / 38400 / 115200
- **数据位**: 8
- **停止位**: 1
- **校验**: 无 / 偶校验 / 奇校验
- **功能码**: 0x03 (读取保持寄存器)
- **从站地址**: 可配置 (默认 0x01)

**读取命令帧** (Modbus RTU):
```
[地址] [功能码] [起始寄存器H] [起始寄存器L] [数量H] [数量L] [CRC_L] [CRC_H]
0x01   0x03      0x00           0x00         0x00    0x01    0x84    0x0A
```

**响应帧**:
```
[地址] [功能码] [字节数] [数据H] [数据L] [CRC_L] [CRC_H]
0x01   0x03      0x02     0x04    0xD2    0x39    0x85
→ 0x04D2 = 1234 → 123.4 kg (除以10)
```

### 2.2 CONTINUOUS (连续发送模式 — 国产地磅)

大部分国产地磅使用此模式。设备持续以 ASCII 文本形式发送重量数据。

**格式 1 — 标准格式**:
```
ST,GS,+0123.4 kg\r\n
ST = 稳定(Stable), US = 不稳定(Unstable)
GS = 毛重(Gross), NT = 净重(Net)
```

**格式 2 — 简化格式**:
```
  123.45 kg\r\n
```

**命令**:
| 命令 | ASCII | 功能 |
|------|-------|------|
| 去皮 | `T\r\n` | 设置当前重量为皮重 |
| 归零 | `Z\r\n` | 重量清零 |
| 读取 | `R\r\n` | 请求一次读数 |

### 2.3 TOLEDO (托利多专用协议)

Mettler-Toledo 地磅使用的连续输出格式。

**帧格式**:
```
<STX> <状态A> <状态B> <状态C> <重量6位> <单位2位> <ETX>
0x02  S      W      A      001234   kg     0x03
```

**状态字节**:
- **Status A**: 空格=有效, `?`=无效, `S`=过载, `-`=欠载
- **Status B**: ` `=毛重, `N`=净重
- **Status C**: ` `=正常, `M`=运动/不稳定, `*`=零点

**命令**:
| 命令 | 帧 | 功能 |
|------|-----|------|
| 去皮 | `<STX>T<ETX>` | 去皮 |
| 归零 | `<STX>Z<ETX>` | 归零 |
| 打印 | `<STX>P<ETX>` | 打印 |

### 2.4 CUSTOM (自定义协议)

支持配置帧头、帧尾和校验方式。

**配置参数**:
```yaml
protocol: "CUSTOM"
frame_head: [0x02]    # STX
frame_tail: [0x03]    # ETX
data_start: 1         # 数据起始偏移
data_len: 8           # 数据长度
```

---

## 3. Hardware Wiring Guide (硬件接线指南)

### 3.1 RS-232 串口 (DB9 公头)

```
DB9 Pinout (DTE — 计算机端):
┌─────────────────────┐
│ \ 1  2  3  4  5 /   │
│  \ 6  7  8  9  /    │
└──────────────────────┘

Pin  Signal  Direction  接法
1    DCD     Input      —
2    RXD     Input      接设备 TXD
3    TXD     Output     接设备 RXD
4    DTR     Output     —
5    GND     —          接设备 GND
6    DSR     Input      —
7    RTS     Output     —
8    CTS     Input      —
9    RI      Input      —
```

**最少接线**: Pin 2 (RXD), Pin 3 (TXD), Pin 5 (GND)

### 3.2 RS-485 两线制

```
主机 (USB-485转换器)          设备 (地磅/入库机)
  A+ ───────────────────────── A+
  B- ───────────────────────── B-
  GND ──────────────────────── GND

注意事项:
• 终端电阻: 120Ω 并联在总线两端
• 最长距离: 1200m @ 9600 bps
• 推荐线缆: 双绞屏蔽线 (STP)
• 偏置电阻: 上拉 A+ → VCC (680Ω), 下拉 B- → GND (680Ω)
```

### 3.3 USB-串口转换器

推荐型号:
- **FTDI FT232RL**: 最稳定, 支持所有波特率
- **CH340G**: 性价比高, 9600-115200
- **CP2102**: 工业级, 宽温

Linux 设备节点: `/dev/ttyUSB0`, `/dev/ttyUSB1`, ...

---

## 4. Configuration (配置说明)

### 4.1 配置文件 `configs/device-gateway.yaml`

```yaml
server:
  port: 9100          # 设备网关 HTTP 端口
  ws_port: 9101       # WebSocket 实时推送端口

wms:
  api_url: "http://localhost:8080"   # WMS 服务器地址
  api_key: ""                         # API 密钥 (可选)

devices:
  scales:             # 地磅列表
    - id: "SCALE-001"
      port: "/dev/ttyUSB0"
      baud: 9600
      protocol: "CONTINUOUS"   # MODBUS_RTU | CONTINUOUS | TOLEDO | CUSTOM
      warehouse: "厦门仓"

  conveyors:          # 入库机列表
    - id: "CONV-001"
      port: "/dev/ttyUSB1"
      baud: 115200
      protocol: "MODBUS_RTU"
      warehouse: "厦门仓"

  scanners:           # 扫码枪列表
    - id: "SCAN-001"
      port: "/dev/ttyUSB2"
      baud: 9600
```

### 4.2 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DEVICE_GATEWAY_CONFIG` | 配置文件路径 | `configs/device-gateway.yaml` |

---

## 5. API Reference (API 参考)

### 5.1 Device Gateway → WMS API

#### GET /api/device/inbound-task
查询入库任务。

**参数**: `?barcode=SF1234567890`

**响应**:
```json
{
  "id": 1,
  "waybill_no": "SF1234567890",
  "tracking_number": "SF1234567890",
  "product_name": "手机壳",
  "declared_weight": 0.35,
  "location_code": "A-01-03",
  "status": 0
}
```

#### POST /api/device/weight-record
记录称重数据。

**请求体**:
```json
{
  "waybill_no": "SF1234567890",
  "gross_weight": 0.362,
  "net_weight": 0.362,
  "declared_weight": 0.35,
  "scale_id": "SCALE-001",
  "status": 0
}
```

**响应**:
```json
{
  "id": 0,
  "waybill_no": "SF1234567890",
  "gross_weight": 0.362,
  "net_weight": 0.362,
  "declared_weight": 0.35,
  "weight_diff": 0.012,
  "scale_id": "SCALE-001",
  "status": 1
}
```

**重量验证逻辑**: `|实际—申报| / 申报 > 5%` → 标记为异常 (status=2)

#### POST /api/device/inbound-confirm
确认入库完成。

**请求体**:
```json
{
  "waybill_no": "SF1234567890",
  "location_code": "A-01-03"
}
```

#### POST /api/device/heartbeat
设备心跳。

**请求体**:
```json
{
  "device_id": "SCALE-001"
}
```

### 5.2 Device Gateway 本地 API (Port 9100)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查, 返回设备数量 |
| GET | `/api/sessions` | 查看所有设备连接状态 |
| GET | `/api/scale/{id}/read` | 手动读取地磅重量 |
| POST | `/api/scale/{id}/tare` | 地磅去皮 |
| POST | `/api/scale/{id}/zero` | 地磅归零 |

---

## 6. Database Tables (数据库表)

### weight_records — 称重记录

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| tenant_id | INT | 租户ID |
| waybill_no | VARCHAR(50) | 运单号 |
| gross_weight | DECIMAL(10,3) | 毛重 |
| tare_weight | DECIMAL(10,3) | 皮重 |
| net_weight | DECIMAL(10,3) | 净重 |
| declared_weight | DECIMAL(10,3) | 申报重量 |
| weight_diff | DECIMAL(10,3) | 差异 |
| scale_id | VARCHAR(20) | 地磅编号 |
| weigh_time | TIMESTAMPTZ | 称重时间 |
| status | SMALLINT | 0:待确认 1:已确认 2:异常 |

### inbound_tasks — 入库任务

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| waybill_no | VARCHAR(50) | 运单号 |
| tracking_number | VARCHAR(100) | 快递单号 |
| location_code | VARCHAR(20) | 目标库位 |
| conveyor_id | VARCHAR(20) | 入库机ID |
| status | SMALLINT | 0:待执行 1:执行中 2:已完成 3:异常 |

### device_registry — 设备注册表

| 字段 | 类型 | 说明 |
|------|------|------|
| device_type | VARCHAR(30) | scale / conveyor / scanner |
| device_id | VARCHAR(50) UNIQUE | 设备编号 |
| protocol | VARCHAR(30) | MODBUS_RTU / CONTINUOUS / TOLEDO |
| port_config | TEXT | JSON 端口配置 |
| last_heartbeat | TIMESTAMPTZ | 最后心跳时间 |

---

## 7. Troubleshooting Guide (故障排查)

### 7.1 串口通信问题

**症状: `read error: bad file descriptor`**

```bash
# 1. 检查设备节点是否存在
ls -la /dev/ttyUSB*

# 2. 检查权限 (当前用户需在 dialout 组)
groups $USER
sudo usermod -a -G dialout $USER

# 3. 测试串口通信
stty -F /dev/ttyUSB0 9600 cs8 -cstopb -parenb
cat /dev/ttyUSB0  # 查看原始数据

# 4. 串口调试工具
sudo apt install minicom
minicom -D /dev/ttyUSB0 -b 9600
```

### 7.2 称重数据异常

**症状: 重量始终为 0 或超大值**

```
排查步骤:
1. 检查协议是否匹配 (CONTINUOUS vs MODBUS_RTU vs TOLEDO)
2. 用 minicom 查看原始数据格式
3. 检查地磅是否已归零
4. 检查 Modbus 从站地址是否匹配 (默认 0x01)
5. 检查 CRC 校验 — Modbus 响应帧最后2字节
```

### 7.3 扫码枪无响应

**症状: 扫描条码无反应**

```
排查步骤:
1. 确认扫码枪为 "串口模式" (非 USB HID 键盘模式)
   通常扫描配置条码 "RS-232 模式"
2. 确认波特率匹配 (常见 9600)
3. 确认换行符: 扫码枪通常发送 CR+LF (\r\n)
4. 测试: cat /dev/ttyUSB2 后扫描, 应看到条码数字
```

### 7.4 WMS 连接失败

**症状: `WMS query failed`**

```
排查步骤:
1. 确认 WMS 服务器运行: curl http://localhost:8080/api/v1/health
2. 检查 device-gateway.yaml 中 api_url 配置
3. 检查防火墙: WMS 端口 8080 是否开放
4. 查看 WMS 日志: 确认 /api/device/* 路由已注册
```

---

## 8. Build & Run (编译运行)

```bash
# 进入设备网关目录
cd apps/device-gateway

# 下载依赖
go mod tidy

# 编译
go build -o bin/device-gateway ./cmd/gateway/

# 运行 (需要配置文件)
./bin/device-gateway

# 或指定配置文件
DEVICE_GATEWAY_CONFIG=/path/to/config.yaml ./bin/device-gateway

# 健康检查
curl http://localhost:9100/api/health
# → {"devices":3,"scales":1,"scanners":1,"service":"device-gateway","status":"ok","version":"1.0.0"}
```

---

## 9. Benchmark Data (性能基准)

### Scan-to-Weight Latency (扫描到称重延迟)

| 环节 | 典型延迟 | 说明 |
|------|---------|------|
| 条码扫描 | < 50ms | 扫码枪输出 |
| WMS 任务查询 | 5-20ms | 局域网 HTTP 调用 |
| 皮带传输 | 1-3s | 物理传送时间 |
| 称重稳定 | 0.3-1.5s | 地磅稳定窗口 |
| 重量记录 | 5-15ms | 返回 WMS |
| 分拨指令 | 10-30ms | 发送到入库机 |
| **总计** | **2-5s** | 从扫描到分拨完成 |

### Throughput (吞吐量)

- 单台入库机: 600-900 件/小时
- 双入库机并行: 1200-1600 件/小时
- 设备网关 CPU: < 5% (单核)
- 内存: < 50MB
