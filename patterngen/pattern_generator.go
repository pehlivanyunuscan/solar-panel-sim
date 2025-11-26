package patterngen

import (
	"math/rand"
	"time"
)

// Günlük bir desen oluşturur
// startMinute: Desenin başlangıç dakikası (0-1439 arası)
// endMinute: Desenin bitiş dakikası (0-1439 arası, startMinute'dan büyük olmalı)
func GenerateDailyPattern(startMinute, endMinute int, maxValue float64, seed int64) []float64 {
	// Günlük rastgele parametreler
	r := rand.New(rand.NewSource(seed))
	peakNoise := 1 + (r.NormFloat64() * 0.07) // %7 civarı tepe oynaklığı
	morningSlope := 0.18 + r.Float64()*0.04   // sabah artış eğimi (0.18-0.22)
	peakHour := 12.0 + r.NormFloat64()*0.5    // tepe saati 12 civarı oynar

	pattern := make([]float64, endMinute-startMinute)
	for minute := startMinute; minute < endMinute; minute++ {
		hour := float64(minute) / 60.0
		var value float64

		switch {
		case hour < 6.0:
			value = 0
		case hour < 8.0: // Sabah yavaş artış (random slope)
			value = maxValue * morningSlope * (hour - 6.0) / 2.0 * peakNoise
		case hour < peakHour: // Sabah hızlı artış (random slope)
			value = maxValue * (morningSlope + (1-morningSlope)*(hour-8.0)/(peakHour-8.0)) * peakNoise
		case hour < 15.0: // Öğlen (peak, noise ile)
			base := maxValue * peakNoise
			// use seeded rand.Rand instance for deterministic noise per-seed
			noise := r.NormFloat64() * maxValue * 0.15 // %15 oynaklık
			value = base + noise
		case hour < 18.0: // Akşam yavaş azalış (random slope)
			// 15:00'da mevcut değerden başla, 18:00'da 0'a in
			fraction := (hour - 15.0) / 3.0 // 15:00'da 0, 18:00'da 1
			// Düşüş eğrisini belirle (başlangıç değeri: peak, bitiş: 0)
			base := maxValue * peakNoise * (1.0 - fraction)
			// Noise ekle ama sadece negatif (azaltıcı) olacak şekilde
			// use seeded rand.Rand instance for deterministic noise per-seed
			noise := -1 * (r.Float64() * base * 0.08) // %8'e kadar negatif oynaklık
			value = base + noise
			if value < 0 {
				value = 0
			}
		default:
			value = 0
		}

		if value < 0 {
			value = 0
		}
		pattern[minute-startMinute] = value
	}
	return pattern
}

// Şu anki günün dakikasını döner (0-1439 arası)
// 0 = 00:00, 1439 = 23:59
// Bu fonksiyon, günün saatini ve dakikasını kullanarak toplam dakikayı hesaplar.
func GetCurrentMinuteOfDay() int {
	now := time.Now()
	return now.Hour()*60 + now.Minute()
}

func GetPatternValueForNow(pattern []float64, startMinute, endMinute int) float64 {
	minuteOfDay := GetCurrentMinuteOfDay()
	if minuteOfDay < startMinute || minuteOfDay >= endMinute {
		return 0
	}
	return pattern[minuteOfDay-startMinute]
}
