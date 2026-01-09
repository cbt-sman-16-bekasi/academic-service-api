package school_response

// SchoolBriefInfo - Info singkat sekolah untuk halaman public (login, dll)
type SchoolBriefInfo struct {
	SchoolCode string `json:"school_code"`
	SchoolName string `json:"school_name"`
	Logo       string `json:"logo"`
	Banner     string `json:"banner"`
	Address    string `json:"address"`
}

type ConfigurationResponse struct {
	ClientID    string           `json:"client_id"` // Public identifier untuk API key validation
	ThemePreset string           `json:"theme_preset"`
	ThemeCustom string           `json:"theme_custom"`
	SchoolInfo  *SchoolBriefInfo `json:"school_info"` // Info sekolah untuk halaman login
}
