package models

// Panel tanımı - il, ilçe, mahalle eklendi
type Panel struct {
	/*
		City           string  `json:"city"`
		District       string  `json:"district"`
		Neighborhood   string  `json:"neighborhood"` */
	Brand          string  `json:"brand"`
	MaxPower       float64 `json:"max_power"`
	Pattern        []float64
	LastPatternDay int
}

// Panel listesi
var PanelList = []Panel{
	{MaxPower: 3000, LastPatternDay: -1},
}
