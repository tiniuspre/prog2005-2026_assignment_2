package models

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Country struct {
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Capital     string      `json:"capital,omitempty"`
	Population  int64       `json:"population,omitempty"`
	Area        float64     `json:"area,omitempty"`
	Region      string      `json:"region,omitempty"`
	Languages   []string    `json:"languages,omitempty"`
	Currencies  []string    `json:"currencies,omitempty"`
	Borders     []string    `json:"borders,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
}
