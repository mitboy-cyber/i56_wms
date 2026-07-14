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

// Redirect sends an HTTP 303 See Other redirect.
// Plain HTTP redirect — most reliable cross-browser, no JS/HTMX dependency.
func Redirect(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(303)
}

// ===================================================================
// Modal form helpers (BDL 1.0 modal overlay)
// ===================================================================

// ModalStart opens a modal overlay div.
func ModalStart(title string) string {
	return `<div class="i56-modal-overlay" onclick="if(event.target===this) closeI56Modal()"><div class="i56-modal-content"><div class="i56-modal-header"><span class="i56-modal-title">` + title + `</span><button class="i56-modal-close" onclick="closeI56Modal()">&times;</button></div><div class="i56-modal-body">`
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

// FormSave opens a plain HTML form element.
// Pure HTML POST + 303 redirect — no HTMX AJAX, most reliable.
func FormSave(action string) string {
	return fmt.Sprintf(`<form action="%s" method="POST">`, action)
}

// FormFooter closes a modal form with cancel and submit buttons.
func FormFooter() string {
	return `<div class="i56-modal-footer"><button type="button" class="i56-btn" onclick="closeI56Modal()">取消</button><button type="submit" class="i56-btn i56-btn-primary">保存</button></div></form>`
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
	// Clean up Sprintf artifacts like ") " that become ")"
	for strings.Contains(v, "))") {
		v = strings.ReplaceAll(v, "))", ")")
	}
	// Only strip the trailing ")" if it's a Sprintf artifact (i.e., there's no matching opening paren)
	// Don't strip closing parens that are part of legitimate parenthetical text (e.g., "厦门→台湾(海快)")
	v = strings.ReplaceAll(v, "sea_express", "海快")
	v = strings.ReplaceAll(v, "sea", "海运")
	v = strings.ReplaceAll(v, "air", "空运")
	// Cargo type Chinese labels
	v = strings.ReplaceAll(v, "general", "普货")
	v = strings.ReplaceAll(v, "sensitive", "特货")
	v = strings.ReplaceAll(v, "dangerous", "危险品")
	// Device type Chinese labels
	v = strings.ReplaceAll(v, "scale", "地磅")
	v = strings.ReplaceAll(v, "conveyor", "入库机")
	v = strings.ReplaceAll(v, "scanner", "扫码枪")
	return v
}

// formatCell fixes common Sprintf format errors in cell values.
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

// FormImageUpload renders a file input for image upload.
// Files are auto-uploaded via JS fetch to /admin/upload/parcel-image.
// URLs are stored in a hidden field "uploaded_urls" for form submission.
func FormImageUpload(label string) string {
	return fmt.Sprintf(`<div class="form-group">
		<label class="form-label">%s</label>
		<input type="file" accept="image/*" multiple
			style="display:block;font-size:12px;margin-bottom:4px"
			onchange="i56UploadImages(this)" />
		<div class="i56-preview-row" style="display:flex;gap:4px;flex-wrap:wrap;margin-top:4px"></div>
		<input type="hidden" name="uploaded_urls" value="" />
	</div>
	<script>
	if(!window._i56UploadInit){
		window._i56UploadInit=true;
		async function i56UploadImages(el){
			var files=el.files;if(!files.length)return;
			var fd=new FormData();
			for(var f of files)fd.append("images",f);
			var r=await fetch("/admin/upload/parcel-image",{method:"POST",body:fd});
			var d=await r.json();
			if(d.ok&&d.urls){
				var h=el.parentElement.querySelector("[name=uploaded_urls]");
				if(h.value)h.value+=",";
				h.value+=d.urls.join(",");
				var pv=el.parentElement.querySelector(".i56-preview-row");
				for(var u of d.urls){
					var img=document.createElement("img");
					img.src=u;img.style="width:64px;height:64px;object-fit:cover;border-radius:4px;border:1px solid #ddd";
					pv.appendChild(img);
				}
			}
		}
	}
	</script>`, label)
}
