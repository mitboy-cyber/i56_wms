// Package common provides shared helpers for admin route modules.
package common

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ===================================================================
// Parsing helpers
// ===================================================================

// ParseID parses an int64 from a form value string.
func ParseID(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

// ParseFloat parses a float64 from a form value string.
func ParseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// StatusLabelText returns a plain text status label.
func StatusLabelText(active bool) string {
	if active {
		return "启用"
	}
	return "停用"
}

// ===================================================================
// Response helpers
// ===================================================================

// HtmlOK sets the Content-Type header for HTML responses.
func HtmlOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// Redirect sends an HTMX HX-Redirect response.
func Redirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(200)
}

// ===================================================================
// Modal form helpers (BDL 1.0 modal overlay)
// ===================================================================

// ModalStart opens a modal overlay div.
func ModalStart(title string) string {
	return `<div class="modal-overlay" onclick="event.target===this&&this.remove()"><div class="modal-content"><div class="modal-header"><span class="modal-title">` + title + `</span><button class="modal-close" onclick="this.closest('.modal-overlay').remove()">&times;</button></div><div class="modal-body">`
}

// ModalEnd closes a modal overlay div.
func ModalEnd() string { return `</div></div></div>` }

// FormField renders a form field with label and input.
func FormField(label, name, value, placeholder string) string {
	return fmt.Sprintf(`<div class="form-group"><label class="form-label">%s</label><input name="%s" value="%s" class="form-input" placeholder="%s"></div>`, label, name, value, placeholder)
}

// FormSelect renders a form field with a select dropdown.
func FormSelect(label, name, value string, opts ...[2]string) string {
	h := fmt.Sprintf(`<div class="form-group"><label class="form-label">%s</label><select name="%s" class="form-input">`, label, name)
	for _, o := range opts {
		sel := ""
		if o[0] == value {
			sel = " selected"
		}
		h += fmt.Sprintf(`<option value="%s"%s>%s</option>`, o[0], sel, o[1])
	}
	return h + `</select></div>`
}

// FormSave opens a form element with HTMX post action.
func FormSave(action string) string {
	return fmt.Sprintf(`<form hx-post="%s" hx-swap="none">`, action)
}

// FormFooter closes a modal form with cancel and submit buttons.
func FormFooter() string {
	return `<div class="modal-footer"><button type="button" class="btn" onclick="this.closest('.modal-overlay').remove()">取消</button><button type="submit" class="btn btn-primary">保存</button></div></form>`
}

// ===================================================================
// Rendering context
// ===================================================================

// ExecTpl is a function type for executing templates.
type ExecTpl func(tmpl map[string]*template.Template, key string, w http.ResponseWriter, name string, data any)

// GenericListFunc renders a page through the BDL generic_list template.
type GenericListFunc func(w http.ResponseWriter, page, title string, total int, cols []string, rows [][]string, addURL ...string)

// RenderCtx bundles template rendering helpers.
type RenderCtx struct {
	Tmpl  map[string]*template.Template
	Exec  ExecTpl
}

// DefaultExecTpl is the standard template execution helper.
func DefaultExecTpl(tmpl map[string]*template.Template, key string, w http.ResponseWriter, name string, data any) {
	tmpl[key].ExecuteTemplate(w, name, data)
}

// formatCell fixes common Sprintf format errors in cell values.
func formatCell(v string) string {
	v = strings.ReplaceAll(v, "%!(EXTRA int64=", "")
	v = strings.ReplaceAll(v, "%!(EXTRA float64=", "")
	v = strings.ReplaceAll(v, ") ", ")")
	for strings.Contains(v, "))") {
		v = strings.ReplaceAll(v, "))", ")")
	}
	v = strings.TrimRight(v, ")")
	v = strings.ReplaceAll(v, "sea_express", "海快")
	v = strings.ReplaceAll(v, "sea", "海运")
	v = strings.ReplaceAll(v, "air", "空运")
	return v
}

// NewGenericList creates a genericList closure using the render context.
func (rc *RenderCtx) NewGenericList() GenericListFunc {
	return func(w http.ResponseWriter, page, title string, total int, cols []string, rows [][]string, addURL ...string) {
		fmtRows := make([][]string, len(rows))
		for i, row := range rows {
			fmtRow := make([]string, len(row))
			for j, cell := range row {
				fmtRow[j] = formatCell(cell)
			}
			fmtRows[i] = fmtRow
		}
		data := map[string]any{
			"Page":       page,
			"Title":      title,
			"Total":      total,
			"Columns":    cols,
			"Rows":       fmtRows,
			"HasActions": true,
		}
		if len(addURL) > 0 && addURL[0] != "" {
			data["AddURL"] = addURL[0]
		}
		rc.Exec(rc.Tmpl, "generic_list", w, "generic_list.html", data)
	}
}
