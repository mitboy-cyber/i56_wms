package domain
import "time"

type PDAVersion struct {
	ID           int64     `json:"id"`
	TenantID     int64     `json:"tenant_id"`
	VersionCode  int       `json:"version_code"`
	VersionName  string    `json:"version_name"`
	ReleaseNotes string    `json:"release_notes"`
	DownloadURL  string    `json:"download_url"`
	Platform     string    `json:"platform"` // android | ios | pwa
	ForceUpdate  bool      `json:"force_update"`
	IsActive     bool      `json:"is_active"`
	PublishedAt  time.Time `json:"published_at"`
}
