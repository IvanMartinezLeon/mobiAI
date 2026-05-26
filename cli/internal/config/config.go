package config

import "encoding/json"

type Settings struct {
	QuietStartup    *bool  `json:"quietStartup,omitempty"`
	Theme           string `json:"theme,omitempty"`
	DefaultProvider string `json:"defaultProvider,omitempty"`
	DefaultModel    string `json:"defaultModel,omitempty"`
}

func Parse(data []byte) (*Settings, error) {
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Settings) Marshal() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func Merge(base, overlay *Settings) *Settings {
	if base == nil {
		base = &Settings{}
	}
	if overlay == nil {
		return base
	}
	if overlay.QuietStartup != nil {
		base.QuietStartup = overlay.QuietStartup
	}
	if overlay.Theme != "" {
		base.Theme = overlay.Theme
	}
	if overlay.DefaultProvider != "" {
		base.DefaultProvider = overlay.DefaultProvider
	}
	if overlay.DefaultModel != "" {
		base.DefaultModel = overlay.DefaultModel
	}
	return base
}
