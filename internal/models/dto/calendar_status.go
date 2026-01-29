package dto

type Calendar struct {
	Title      string   `json:"title"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Display    string   `json:"display,omitempty"`
	ClassNames []string `json:"classNames,omitempty"`
	Color      string   `json:"color,omitempty"`
	Background string   `json:"backgroundColor,omitempty"`
	TextColor  string   `json:"textColor,omitempty"`
}
