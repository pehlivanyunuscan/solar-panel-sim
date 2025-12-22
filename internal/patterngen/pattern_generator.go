package patterngen

import (
	"math/rand"
	"time"
)

// GenerateDailyPattern günlük güneş paneli güç deseni oluşturur
// startMinute: Desenin başlangıç dakikası (0-1439)
// endMinute: Desenin bitiş dakikası (0-1439)
// maxValue: Maksimum güç değeri (Watt)
// seed: Rastgele sayı üretici tohumu
func GenerateDailyPattern(startMinute, endMinute int, maxValue float64, seed int64) []float64 {
	r := rand.New(rand.NewSource(seed))

	// Günlük rastgele parametreler
	peakNoise := 1 + (r.NormFloat64() * 0.07) // %7 tepe oynaklığı
	morningSlope := 0.18 + r.Float64()*0.04   // sabah artış eğimi (0.18-0.22)
	peakHour := 12.0 + r.NormFloat64()*0.5    // tepe saati ~12:00

	pattern := make([]float64, endMinute-startMinute)

	for minute := startMinute; minute < endMinute; minute++ {
		hour := float64(minute) / 60.0
		var value float64

		switch {
		case hour < 6.0: // Gece - güç yok
			value = 0

		case hour < 8.0: // Sabah yavaş artış
			value = maxValue * morningSlope * (hour - 6.0) / 2.0 * peakNoise

		case hour < peakHour: // Sabah hızlı artış
			value = maxValue * (morningSlope + (1-morningSlope)*(hour-8.0)/(peakHour-8.0)) * peakNoise

		case hour < 15.0: // Öğlen - tepe güç
			base := maxValue * peakNoise
			noise := r.NormFloat64() * maxValue * 0.15 // %15 oynaklık
			value = base + noise

		case hour < 18.0: // Akşam - yavaş düşüş
			fraction := (hour - 15.0) / 3.0
			base := maxValue * peakNoise * (1.0 - fraction)
			noise := -1 * (r.Float64() * base * 0.08) // %8 negatif oynaklık
			value = base + noise

		default: // Gece - güç yok
			value = 0
		}

		if value < 0 {
			value = 0
		}
		pattern[minute-startMinute] = value
	}
	return pattern
}

// GetCurrentMinuteOfDay günün mevcut dakikasını döndürür (0-1439)
func GetCurrentMinuteOfDay() int {
	now := time.Now()
	return now.Hour()*60 + now.Minute()
}

// GetPatternValueForNow şu anki dakika için desen değerini döndürür
func GetPatternValueForNow(pattern []float64, startMinute, endMinute int) float64 {
	minuteOfDay := GetCurrentMinuteOfDay()
	if minuteOfDay < startMinute || minuteOfDay >= endMinute || len(pattern) == 0 {
		return 0
	}
	idx := minuteOfDay - startMinute
	if idx >= len(pattern) {
		return 0
	}
	return pattern[idx]
}
