package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

type ExportHandler struct{}

func NewExportHandler() *ExportHandler { return &ExportHandler{} }

// WriteCSV writes a slice of structs as CSV response.
func (h *ExportHandler) WriteCSV(w http.ResponseWriter, filename string, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice { return fmt.Errorf("data must be a slice") }
	if v.Len() == 0 { return fmt.Errorf("no data") }

	first := v.Index(0)
	if first.Kind() == reflect.Ptr { first = first.Elem() }
	t := first.Type()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s.csv", filename, time.Now().Format("20060102")))
	w.Header().Set("X-Content-Type-Options", "nosniff")

	writer := csv.NewWriter(w)
	writer.UseCRLF = true

	// Write BOM for Excel UTF-8 compatibility
	w.Write([]byte{0xEF, 0xBB, 0xBF})

	// Header
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" { tag = f.Name }
		headers = append(headers, tag)
	}
	writer.Write(headers)

	// Rows
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr { elem = elem.Elem() }
		var row []string
		for j := 0; j < elem.NumField(); j++ {
			row = append(row, fmt.Sprintf("%v", elem.Field(j).Interface()))
		}
		writer.Write(row)
	}
	writer.Flush()
	return writer.Error()
}
