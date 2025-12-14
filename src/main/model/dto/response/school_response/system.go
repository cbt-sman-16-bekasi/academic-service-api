package school_response

type ConfigurationResponse struct {
	SchoolCode  string `json:"school_code"`
	Key         string `json:"key"`
	ThemePreset string `json:"theme_preset"`
	ThemeCustom string `json:"theme_custom"`
}
