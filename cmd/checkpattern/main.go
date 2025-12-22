package main

import (
	"fmt"

	"main.go/internal/models"
	"main.go/internal/patterngen"
)

func main() {
	p := models.PanelList[0]
	// Sabit seed ile tekrarlanabilir sonuç
	p.Pattern = patterngen.GenerateDailyPattern(0, 1440, p.MaxPower, 42)

	// Kontrol edilecek dakikalar (günün dakikası)
	times := []int{
		0,    // 00:00
		359,  // 05:59
		360,  // 06:00
		420,  // 07:00
		480,  // 08:00
		600,  // 10:00
		720,  // 12:00
		780,  // 13:00
		900,  // 15:00
		1079, // 17:59
		1080, // 18:00
		1200, // 20:00
		1439, // 23:59
	}

	fmt.Println("Örnek desen değerleri:")
	fmt.Println("========================")
	for _, m := range times {
		if m < 0 || m >= len(p.Pattern) {
			fmt.Printf("%02d:%02d -> aralık dışı\n", m/60, m%60)
			continue
		}
		fmt.Printf("%02d:%02d -> %.2f W\n", m/60, m%60, p.Pattern[m])
	}
}
