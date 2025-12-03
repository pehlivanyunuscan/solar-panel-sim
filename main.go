package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gofiber/fiber/v2"
	"main.go/logging"
	"main.go/models"
	"main.go/patterngen"
)

// Panel listesi - data paketinden alıyoruz artık
var panelList []models.Panel

// Panel seed - her container için farklı pattern üretmek için
var panelSeed int64

var (
	sensorLabels = []string{
		"aku gerilimi",
		"panel gerilimi",
		"sarj akimi",
		"yuk akimi",
		"sicaklik",
		"soc",
		"yuk gucu",
		"panel gucu",
		"yuk durum",
		"aku tipi",
		"sarj durum",
		"kapi bilgisi",
		"sarj gucu",
		"panel akim",
	}

	roleLabels = []string{
		"role 1",
		"role 2",
		"role 3",
		"role 4",
		"role 5",
		"role 6",
		"role 7",
	}

	once         sync.Once
	sensorGauges map[string]*metrics.Gauge
	// Pattern için global değişkenler
	startMinute = 0    // 00:00
	endMinute   = 1440 // 24*60 = 1440
	// maxPanelGucu     = 1000.0 // örnek için
	// panelGucuPattern []float64
)

// updatePanelPatternIfNeeded güncel panel gücü desenini günceller
// Eğer gün değiştiyse yeni bir desen oluşturur.
func updatePanelPatternIfNeeded() {
	now := time.Now()
	day := now.YearDay()
	for i := range panelList {
		if panelList[i].LastPatternDay != day {
			// Her panel için farklı seed - panelSeed kullanarak her container farklı pattern üretir
			seed := int64(day) + int64(i)*1000 + panelSeed
			panelList[i].Pattern = patterngen.GenerateDailyPattern(startMinute, endMinute, panelList[i].MaxPower, seed)
			panelList[i].LastPatternDay = day
		}
	}
}

func initGauges() {
	sensorGauges = make(map[string]*metrics.Gauge)
	// Normal sensörler
	for _, sensor := range sensorLabels {
		key := fmt.Sprintf(`mppt_values{sensor="%s"}`, sensor)
		g := metrics.GetOrCreateGauge(key, nil)
		sensorGauges[key] = g
	}
	// Role durumları, sensör olarak "role durumlari" ve ek olarak "role" etiketi
	for _, role := range roleLabels {
		key := fmt.Sprintf(`mppt_values{sensor="role durumlari",role="%s"}`, role)
		g := metrics.GetOrCreateGauge(key, nil)
		sensorGauges[key] = g
	}

	// time_active metric removed — do not export time_active to Prometheus
}

func randomValue(sensor string) float64 {
	switch sensor {
	case "aku gerilimi", "panel gerilimi":
		return float64(rand.IntN(1500) + 1200)
	case "sarj akimi", "yuk akimi", "panel akim":
		return float64(rand.IntN(3000) + 500)
	case "sicaklik":
		return float64(rand.IntN(100))
	case "soc":
		return float64(rand.IntN(101))
	case "yuk gucu", "panel gucu", "sarj gucu":
		return float64(rand.IntN(4000) + 1000)
	case "yuk durum", "aku tipi", "sarj durum":
		return float64(rand.IntN(2)) // 0 veya 1
	case "kapi bilgisi":
		return float64(rand.IntN(3)) // 0, 1 veya 2
	case "role durumlari":
		return float64(rand.IntN(2)) // 0 veya 1
	default:
		return 0.0
	}
}

func main() {

	// PANEL_SEED environment variable'ını oku, yoksa default 0 kullan
	panelSeedStr := os.Getenv("PANEL_SEED")
	if panelSeedStr != "" {
		if seed, err := strconv.ParseInt(panelSeedStr, 10, 64); err == nil {
			panelSeed = seed
		}
	}

	panelList = models.PanelList // models paketinden panel listesini al
	// Logging seviyesini başlat
	logging.SetLogLevel()

	app := fiber.New()
	once.Do(initGauges)

	app.Get("/metrics", func(c *fiber.Ctx) error {
		updatePanelPatternIfNeeded() // Panel gücü desenini güncelle
		// Sensorlar için değer güncelle
		for _, sensor := range sensorLabels {
			var val float64
			// Sadece "panel gucu" için tek bir pattern değeri kullan,
			// diğer sensörler kendi random değerlerini üretmeye devam etsin.
			if sensor == "panel gucu" {
				if len(panelList) > 0 {
					val = patterngen.GetPatternValueForNow(panelList[0].Pattern, startMinute, endMinute)
				} else {
					val = randomValue(sensor)
				}
			} else {
				val = randomValue(sensor)
			}
			key := fmt.Sprintf(`mppt_values{sensor="%s"}`, sensor)
			if g, ok := sensorGauges[key]; ok {
				g.Set(val)
			}
		}
		// Role'ler için değer güncelle
		for _, role := range roleLabels {
			val := randomValue("role durumlari")
			key := fmt.Sprintf(`mppt_values{sensor="role durumlari",role="%s"}`, role)
			if g, ok := sensorGauges[key]; ok {
				g.Set(val)
			}
		}

		// time_active metric disabled — do not set or export it

		// logging.LogApp(logging.INFO, "/metrics endpointine istek geldi. IP: %s", c.IP())
		/* logging.LogAudit(
			"anonymous",
			"/metrics",
			c.Method(),
			http.StatusOK,
			c.IP(),
			nil,
			"Panel metrikleri listelendi",
		) */

		c.Set("Content-Type", "text/plain; version=0.0.4")
		metrics.WritePrometheus(c.Context(), true)
		return nil
	})

	logging.LogApp(logging.INFO, "Uygulama başlatıldı")
	if err := app.Listen("0.0.0.0:8080"); err != nil {
		logging.LogApp(logging.ERROR, "Sunucu başlatılamadı: %v", err)
	}
}
