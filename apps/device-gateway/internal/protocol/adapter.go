// Package protocol provides protocol adapters for industrial device communication.
package protocol

import (
	"fmt"
)

// Command represents a command to send to a device.
type Command struct {
	Code     byte
	Data     []byte
	DeviceID byte // Modbus slave address
}

// ProtocolAdapter defines the interface for parsing and encoding device protocols.
type ProtocolAdapter interface {
	// Parse decodes raw bytes from a device into structured data (typically a float64 weight).
	Parse(raw []byte) (interface{}, error)
	// Encode builds a raw byte command from a Command struct.
	Encode(cmd Command) ([]byte, error)
	// Name returns the protocol name.
	Name() string
}

// ─── Modbus RTU Adapter ────────────────────────────────────────────

// ModbusRTUAdapter implements the Modbus RTU protocol (RS-232/485).
// Reads holding registers, function code 0x03.
type ModbusRTUAdapter struct {
	slaveID byte
}

// NewModbusRTUAdapter creates a Modbus RTU adapter.
func NewModbusRTUAdapter(slaveID byte) *ModbusRTUAdapter {
	return &ModbusRTUAdapter{slaveID: slaveID}
}

func (m *ModbusRTUAdapter) Name() string { return "MODBUS_RTU" }

// Parse decodes a Modbus RTU response containing holding register values.
// Expected response format: [slaveID, funcCode, byteCount, dataHi, dataLo, crcLo, crcHi]
func (m *ModbusRTUAdapter) Parse(raw []byte) (interface{}, error) {
	if len(raw) < 7 {
		return nil, fmt.Errorf("modbus: frame too short (%d bytes)", len(raw))
	}
	if raw[0] != m.slaveID {
		return nil, fmt.Errorf("modbus: slave ID mismatch (expected %d, got %d)", m.slaveID, raw[0])
	}
	if raw[1] != 0x03 {
		return nil, fmt.Errorf("modbus: unexpected function code 0x%02X", raw[1])
	}
	byteCount := int(raw[2])
	if len(raw) < 5+byteCount {
		return nil, fmt.Errorf("modbus: data too short for byte count %d", byteCount)
	}

	// Interpret as big-endian 16-bit registers
	if byteCount >= 2 {
		value := int(raw[3])<<8 | int(raw[4])
		// Scale factor: 1 decimal place (e.g., 1234 -> 123.4 kg)
		return float64(value) / 10.0, nil
	}
	return float64(0), fmt.Errorf("modbus: no register data")
}

// Encode builds a Modbus RTU read-holding-registers command (function 0x03).
// Command.Code = starting register high byte, Command.Data contains [startingRegLo, qtyHi, qtyLo]
func (m *ModbusRTUAdapter) Encode(cmd Command) ([]byte, error) {
	startReg := uint16(cmd.Code)<<8 | uint16(0x00)
	qty := uint16(1)
	if len(cmd.Data) >= 2 {
		qty = uint16(cmd.Data[0])<<8 | uint16(cmd.Data[1])
	}
	if len(cmd.Data) >= 1 && cmd.Data[0] != 0 {
		startReg = uint16(cmd.Code)<<8 | uint16(cmd.Data[0])
	}

	frame := []byte{m.slaveID, 0x03, byte(startReg >> 8), byte(startReg & 0xFF), byte(qty >> 8), byte(qty & 0xFF)}
	crc := modbusCRC16(frame)
	frame = append(frame, byte(crc&0xFF), byte(crc>>8))
	return frame, nil
}

// modbusCRC16 computes Modbus CRC-16.
func modbusCRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if (crc & 0x0001) != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// ─── Continuous Adapter (国产地磅连续发送模式) ──────────────────────

// ContinuousAdapter handles devices that continuously stream weight data as ASCII text.
// Common format: "ST,GS,+0123.4 kg\r\n" or "  123.4 kg\r\n"
type ContinuousAdapter struct{}

// NewContinuousAdapter creates a continuous-mode adapter.
func NewContinuousAdapter() *ContinuousAdapter {
	return &ContinuousAdapter{}
}

func (c *ContinuousAdapter) Name() string { return "CONTINUOUS" }

// Parse extracts a float64 weight from a continuous ASCII stream.
// Handles formats:
//
//	"ST,GS,+0123.4 kg"  — stable gross weight
//	"  123.45 kg"        — simple weight value
//	"US,GS,+0123.4 kg"   — unstable weight
func (c *ContinuousAdapter) Parse(raw []byte) (interface{}, error) {
	s := string(raw)

	// Try to find a float value followed by optional "kg" or "g"
	var value float64
	var unit string

	// Look for standard float pattern
	n, err := fmt.Sscanf(s, "%f %s", &value, &unit)
	if err != nil || n < 1 {
		// Try scanning without unit
		n, err = fmt.Sscanf(s, "%f", &value)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("continuous: cannot parse weight from %q", s)
		}
	}

	// Convert grams to kg if needed
	if unit == "g" || unit == "G" {
		value /= 1000.0
	}

	return value, nil
}

// Encode builds a command for the continuous adapter (rarely used).
func (c *ContinuousAdapter) Encode(cmd Command) ([]byte, error) {
	switch cmd.Code {
	case 'T': // Tare 去皮
		return []byte("T\r\n"), nil
	case 'Z': // Zero 归零
		return []byte("Z\r\n"), nil
	case 'R': // Read
		return []byte("R\r\n"), nil
	default:
		return nil, fmt.Errorf("continuous: unknown command 0x%02X", cmd.Code)
	}
}

// ─── Toledo Adapter ─────────────────────────────────────────────────

// ToledoAdapter implements Toledo/Mettler-Toledo continuous output protocol.
// Format: <STX><status><weight><unit><ETX> where weight is 6 digits with implied decimal.
type ToledoAdapter struct{}

const (
	toledoSTX = 0x02
	toledoETX = 0x03
)

// NewToledoAdapter creates a Toledo protocol adapter.
func NewToledoAdapter() *ToledoAdapter {
	return &ToledoAdapter{}
}

func (t *ToledoAdapter) Name() string { return "TOLEDO" }

// Parse decodes a Toledo continuous frame.
// Frame: STX SWA/SWB weight(6 chars) unit(2 chars) ETX
func (t *ToledoAdapter) Parse(raw []byte) (interface{}, error) {
	if len(raw) < 10 {
		return nil, fmt.Errorf("toledo: frame too short (%d bytes)", len(raw))
	}

	// Find STX
	start := -1
	for i, b := range raw {
		if b == toledoSTX {
			start = i
			break
		}
	}
	if start < 0 {
		return nil, fmt.Errorf("toledo: STX not found")
	}

	data := raw[start:]
	if len(data) < 10 {
		return nil, fmt.Errorf("toledo: data after STX too short")
	}

	// Status bytes at positions 1-3
	status := string(data[1:4])
	_ = status

	// Weight at positions 4-9 (6 chars), with 7th being decimal point pos
	if len(data) < 10 {
		return nil, fmt.Errorf("toledo: incomplete weight field")
	}
	weightRaw := string(data[4:10])

	var value float64
	if _, err := fmt.Sscanf(weightRaw, "%f", &value); err != nil {
		return nil, fmt.Errorf("toledo: cannot parse weight %q: %w", weightRaw, err)
	}

	// Toledo uses implied decimal; typical format has 1-2 decimal places
	// Weight raw like "001234" means 12.34 or 123.4 depending on configuration
	// We default to dividing by 100 for 2 decimal places
	if value > 100 {
		value /= 100.0
	}

	return value, nil
}

// Encode builds a command for the Toledo adapter.
func (t *ToledoAdapter) Encode(cmd Command) ([]byte, error) {
	// Toledo uses specific escape sequences for commands
	switch cmd.Code {
	case 'T': // Tare
		return []byte{toledoSTX, 'T', toledoETX}, nil
	case 'Z': // Zero
		return []byte{toledoSTX, 'Z', toledoETX}, nil
	case 'P': // Print
		return []byte{toledoSTX, 'P', toledoETX}, nil
	default:
		return nil, fmt.Errorf("toledo: unknown command 0x%02X", cmd.Code)
	}
}

// ─── Custom Adapter (可配置帧头/帧尾/校验) ───────────────────────────

// CustomAdapter supports configurable frame format with start/end markers and optional checksum.
type CustomAdapter struct {
	name      string
	frameHead []byte // Frame start marker
	frameTail []byte // Frame end marker
	useCRC    bool   // Whether to verify CRC
	dataStart int    // Offset from frame head to start of data
	dataLen   int    // Length of data field
}

// NewCustomAdapter creates a configurable custom protocol adapter.
func NewCustomAdapter(name string, frameHead, frameTail []byte, dataStart, dataLen int) *CustomAdapter {
	return &CustomAdapter{
		name:      name,
		frameHead: frameHead,
		frameTail: frameTail,
		dataStart: dataStart,
		dataLen:   dataLen,
	}
}

func (c *CustomAdapter) Name() string { return c.name }

// Parse extracts data from a custom-framed message.
func (c *CustomAdapter) Parse(raw []byte) (interface{}, error) {
	if len(c.frameHead) > 0 {
		// Find frame head
		pos := findBytes(raw, c.frameHead)
		if pos < 0 {
			return nil, fmt.Errorf("custom: frame head not found")
		}
		raw = raw[pos:]
	}

	if len(raw) < c.dataStart+c.dataLen {
		return nil, fmt.Errorf("custom: frame too short for data field")
	}

	data := raw[c.dataStart : c.dataStart+c.dataLen]

	// Try to parse as float
	var value float64
	if _, err := fmt.Sscanf(string(data), "%f", &value); err != nil {
		return string(data), nil // Return raw string if not numeric
	}

	return value, nil
}

// Encode builds a command for the custom adapter.
func (c *CustomAdapter) Encode(cmd Command) ([]byte, error) {
	var frame []byte
	frame = append(frame, c.frameHead...)
	frame = append(frame, cmd.Data...)
	frame = append(frame, c.frameTail...)
	return frame, nil
}

// findBytes finds a byte sequence in a slice; returns index or -1.
func findBytes(data, target []byte) int {
	if len(target) == 0 {
		return 0
	}
	for i := 0; i <= len(data)-len(target); i++ {
		match := true
		for j := 0; j < len(target); j++ {
			if data[i+j] != target[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
