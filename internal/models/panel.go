package models

// Panel, güneş paneli simülasyonu için veri modelidir
type Panel struct {
	MaxPower       float64   // Maksimum güç kapasitesi (Watt)
	Pattern        []float64 // Günlük güç deseni (dakikalık değerler)
	LastPatternDay int       // Son desen güncellemesinin yapıldığı gün
}

// PanelList varsayılan panel listesi
var PanelList = []Panel{
	{MaxPower: 3000, LastPatternDay: -1},
}
