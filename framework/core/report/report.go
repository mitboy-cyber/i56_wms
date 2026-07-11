// Package report provides a composable report generation engine.
// Supports CSV, Excel (XLSX), PDF output formats with pluggable data sources.
package report

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Format defines the output format for reports.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatXLSX Format = "xlsx"
	FormatPDF  Format = "pdf"
	FormatJSON Format = "json"
)

// Report defines a report definition with columns, data source, and output.
type Report struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Columns     []Column    `json:"columns"`
	Format      Format      `json:"format"`
	CreatedAt   time.Time   `json:"created_at"`
}

// Column defines a single column in a report.
type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Width int    `json:"width,omitempty"` // for XLSX/PDF
}

// DataSource provides rows for report generation.
type DataSource interface {
	// Fetch retrieves all rows for the report.
	Fetch(ctx context.Context, params map[string]any) ([]map[string]any, error)
	// Count returns total rows (for pagination).
	Count(ctx context.Context, params map[string]any) (int64, error)
}

// Writer is the interface for output-format writers.
type Writer interface {
	Write(ctx context.Context, report *Report, rows []map[string]any, w io.Writer) error
}

// Engine orchestrates report generation.
type Engine struct {
	writers map[Format]Writer
}

// NewEngine creates a report engine with built-in writers.
func NewEngine() *Engine {
	e := &Engine{
		writers: make(map[Format]Writer),
	}
	// Register built-in writers
	e.RegisterWriter(FormatCSV, &CSVWriter{})
	e.RegisterWriter(FormatJSON, &JSONWriter{})
	return e
}

// RegisterWriter adds a custom format writer.
func (e *Engine) RegisterWriter(format Format, w Writer) {
	e.writers[format] = w
}

// Generate fetches data from a DataSource and writes the report to an io.Writer.
func (e *Engine) Generate(ctx context.Context, report *Report, source DataSource, params map[string]any, w io.Writer) error {
	writer, ok := e.writers[report.Format]
	if !ok {
		return fmt.Errorf("report: unsupported format %q", report.Format)
	}

	rows, err := source.Fetch(ctx, params)
	if err != nil {
		return fmt.Errorf("report: fetch: %w", err)
	}

	return writer.Write(ctx, report, rows, w)
}

// GenerateToWriter fetches, optionally counts pagination info, and writes.
func (e *Engine) GenerateToWriter(ctx context.Context, report *Report, source DataSource, params map[string]any, w io.Writer) (*Result, error) {
	writer, ok := e.writers[report.Format]
	if !ok {
		return nil, fmt.Errorf("report: unsupported format %q", report.Format)
	}

	rows, err := source.Fetch(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("report: fetch: %w", err)
	}

	total, err := source.Count(ctx, params)
	if err != nil {
		total = int64(len(rows))
	}

	if err := writer.Write(ctx, report, rows, w); err != nil {
		return nil, fmt.Errorf("report: write: %w", err)
	}

	return &Result{
		ReportID:   report.ID,
		RowCount:   len(rows),
		TotalCount: total,
		Format:     report.Format,
	}, nil
}

// Result holds metadata about a generated report.
type Result struct {
	ReportID   string `json:"report_id"`
	RowCount   int    `json:"row_count"`
	TotalCount int64  `json:"total_count"`
	Format     Format `json:"format"`
}

// CSVWriter writes report rows as CSV.
type CSVWriter struct{}

func (cw *CSVWriter) Write(ctx context.Context, report *Report, rows []map[string]any, w io.Writer) error {
	// Write header
	for i, col := range report.Columns {
		if i > 0 {
			io.WriteString(w, ",")
		}
		io.WriteString(w, escapeCSV(col.Label))
	}
	io.WriteString(w, "\n")

	// Write data rows
	for _, row := range rows {
		for i, col := range report.Columns {
			if i > 0 {
				io.WriteString(w, ",")
			}
			val := fmt.Sprintf("%v", row[col.Key])
			io.WriteString(w, escapeCSV(val))
		}
		io.WriteString(w, "\n")
	}
	return nil
}

func escapeCSV(s string) string {
	needsQuoting := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' {
			needsQuoting = true
			break
		}
	}
	if !needsQuoting {
		return s
	}
	// Escape double-quotes by doubling them
	escaped := ""
	for _, c := range s {
		if c == '"' {
			escaped += "\"\""
		} else {
			escaped += string(c)
		}
	}
	return "\"" + escaped + "\""
}

// JSONWriter writes report rows as a JSON array.
type JSONWriter struct{}

func (jw *JSONWriter) Write(ctx context.Context, report *Report, rows []map[string]any, w io.Writer) error {
	io.WriteString(w, "[\n")
	for i, row := range rows {
		if i > 0 {
			io.WriteString(w, ",\n")
		}
		io.WriteString(w, "  {")
		first := true
		for _, col := range report.Columns {
			if !first {
				io.WriteString(w, ", ")
			}
			first = false
			io.WriteString(w, fmt.Sprintf("%q: %q", col.Key, fmt.Sprintf("%v", row[col.Key])))
		}
		io.WriteString(w, "}")
	}
	if len(rows) > 0 {
		io.WriteString(w, "\n")
	}
	io.WriteString(w, "]\n")
	return nil
}

// NewReport creates a new Report with the given parameters.
func NewReport(id, title string, columns []Column, format Format) *Report {
	return &Report{
		ID:        id,
		Title:     title,
		Columns:   columns,
		Format:    format,
		CreatedAt: time.Now(),
	}
}

// InMemSource is an in-memory DataSource for testing and small datasets.
type InMemSource struct {
	rows []map[string]any
}

// NewInMemSource creates an in-memory data source.
func NewInMemSource(rows []map[string]any) *InMemSource {
	return &InMemSource{rows: rows}
}

func (s *InMemSource) Fetch(ctx context.Context, params map[string]any) ([]map[string]any, error) {
	return s.rows, nil
}

func (s *InMemSource) Count(ctx context.Context, params map[string]any) (int64, error) {
	return int64(len(s.rows)), nil
}
