package report

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestCSVWriter_Write(t *testing.T) {
	report := NewReport("r1", "Orders", []Column{
		{Key: "id", Label: "ID"},
		{Key: "name", Label: "Name"},
	}, FormatCSV)

	rows := []map[string]any{
		{"id": "1", "name": "Order A"},
		{"id": "2", "name": "Order B"},
	}

	var buf bytes.Buffer
	writer := &CSVWriter{}
	err := writer.Write(context.Background(), report, rows, &buf)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ID,Name") {
		t.Errorf("expected header 'ID,Name', got %q", output)
	}
	if !strings.Contains(output, "Order A") {
		t.Errorf("expected 'Order A' in output, got %q", output)
	}
}

func TestCSVWriter_Escape(t *testing.T) {
	report := NewReport("r1", "Test", []Column{
		{Key: "desc", Label: "Description"},
	}, FormatCSV)

	rows := []map[string]any{
		{"desc": `Contains, comma and "quotes"`},
	}

	var buf bytes.Buffer
	(&CSVWriter{}).Write(context.Background(), report, rows, &buf)

	output := buf.String()
	if !strings.Contains(output, `"Contains, comma and ""quotes"""`) {
		t.Errorf("unexpected CSV escaping: %q", output)
	}
}

func TestJSONWriter_Write(t *testing.T) {
	report := NewReport("r1", "Orders", []Column{
		{Key: "id", Label: "ID"},
		{Key: "status", Label: "Status"},
	}, FormatJSON)

	rows := []map[string]any{
		{"id": "1", "status": "shipped"},
		{"id": "2", "status": "pending"},
	}

	var buf bytes.Buffer
	writer := &JSONWriter{}
	err := writer.Write(context.Background(), report, rows, &buf)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"id":`) {
		t.Errorf("expected JSON output, got %q", output)
	}
	if !strings.Contains(output, "shipped") {
		t.Errorf("expected 'shipped' in output, got %q", output)
	}
}

func TestEngine_Generate(t *testing.T) {
	engine := NewEngine()

	report := NewReport("r1", "Orders", []Column{
		{Key: "id", Label: "ID"},
		{Key: "amount", Label: "Amount"},
	}, FormatCSV)

	source := NewInMemSource([]map[string]any{
		{"id": "1", "amount": "100"},
		{"id": "2", "amount": "200"},
	})

	var buf bytes.Buffer
	err := engine.Generate(context.Background(), report, source, nil, &buf)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ID,Amount") {
		t.Errorf("expected CSV header in output: %q", output)
	}
}

func TestEngine_GenerateToWriter(t *testing.T) {
	engine := NewEngine()

	report := NewReport("r2", "Clients", []Column{
		{Key: "name", Label: "Name"},
	}, FormatJSON)

	source := NewInMemSource([]map[string]any{
		{"name": "ACME"},
		{"name": "XYZ"},
	})

	var buf bytes.Buffer
	result, err := engine.GenerateToWriter(context.Background(), report, source, nil, &buf)
	if err != nil {
		t.Fatalf("GenerateToWriter: %v", err)
	}
	if result.RowCount != 2 {
		t.Errorf("expected 2 rows, got %d", result.RowCount)
	}
	if result.TotalCount != 2 {
		t.Errorf("expected total 2, got %d", result.TotalCount)
	}
	if result.Format != FormatJSON {
		t.Errorf("expected JSON format, got %q", result.Format)
	}
}

func TestEngine_UnsupportedFormat(t *testing.T) {
	engine := NewEngine()
	report := NewReport("r1", "Test", nil, "unsupported")
	source := NewInMemSource(nil)

	err := engine.Generate(context.Background(), report, source, nil, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestInMemSource(t *testing.T) {
	data := []map[string]any{
		{"a": 1},
		{"b": 2},
	}
	source := NewInMemSource(data)

	rows, err := source.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}

	count, err := source.Count(context.Background(), nil)
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestNewReport(t *testing.T) {
	r := NewReport("r1", "Test Report", []Column{{Key: "id", Label: "ID"}}, FormatCSV)
	if r.ID != "r1" {
		t.Errorf("expected 'r1', got %q", r.ID)
	}
	if r.Title != "Test Report" {
		t.Errorf("expected 'Test Report', got %q", r.Title)
	}
	if r.Format != FormatCSV {
		t.Errorf("expected CSV, got %q", r.Format)
	}
	if r.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestRegisterWriter(t *testing.T) {
	engine := NewEngine()
	engine.RegisterWriter("custom", &CSVWriter{})

	report := NewReport("r1", "Test", []Column{{Key: "id", Label: "ID"}}, "custom")
	source := NewInMemSource([]map[string]any{{"id": "1"}})

	err := engine.Generate(context.Background(), report, source, nil, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Generate with custom writer: %v", err)
	}
}
