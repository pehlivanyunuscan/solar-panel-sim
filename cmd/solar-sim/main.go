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
	"main.go/internal/logging"
	"main.go/internal/models"
	"main.go/internal/patterngen"
)

// Panel listesi
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
	startMinute  = 0    // 00:00
	endMinute    = 1440 // 24*60 = 1440
)

// updatePanelPatternIfNeeded güncel panel gücü desenini günceller
func updatePanelPatternIfNeeded() {
	now := time.Now()
	day := now.YearDay()
	for i := range panelList {
		if panelList[i].LastPatternDay != day {
			seed := int64(day) + int64(i)*1000 + panelSeed
			panelList[i].Pattern = patterngen.GenerateDailyPattern(startMinute, endMinute, panelList[i].MaxPower, seed)
			panelList[i].LastPatternDay = day
		}
	}
}

func initGauges() {
	sensorGauges = make(map[string]*metrics.Gauge)
	for _, sensor := range sensorLabels {
		key := fmt.Sprintf(`mppt_values{sensor="%s"}`, sensor)
		g := metrics.GetOrCreateGauge(key, nil)
		sensorGauges[key] = g
	}
	for _, role := range roleLabels {
		key := fmt.Sprintf(`mppt_values{sensor="role durumlari",role="%s"}`, role)
		g := metrics.GetOrCreateGauge(key, nil)
		sensorGauges[key] = g
	}
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
		return float64(rand.IntN(2))
	case "kapi bilgisi":
		return float64(rand.IntN(3))
	case "role durumlari":
		return float64(rand.IntN(2))
	default:
		return 0.0
	}
}

func main() {
	// PANEL_SEED environment variable
	panelSeedStr := os.Getenv("PANEL_SEED")
	if panelSeedStr != "" {
		if seed, err := strconv.ParseInt(panelSeedStr, 10, 64); err == nil {
			panelSeed = seed
		}
	}

	panelList = models.PanelList
	logging.SetLogLevel()

	app := fiber.New(fiber.Config{
		AppName: "Solar Panel Simulator",
	})

	once.Do(initGauges)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Prometheus metrics endpoint
	app.Get("/metrics", func(c *fiber.Ctx) error {
		updatePanelPatternIfNeeded()

		for _, sensor := range sensorLabels {
			var val float64
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

		for _, role := range roleLabels {
			val := randomValue("role durumlari")
			key := fmt.Sprintf(`mppt_values{sensor="role durumlari",role="%s"}`, role)
			if g, ok := sensorGauges[key]; ok {
				g.Set(val)
			}
		}

		c.Set("Content-Type", "text/plain; version=0.0.4")
		metrics.WritePrometheus(c.Context(), true)
		return nil
	})

	logging.LogApp(logging.INFO, "Uygulama başlatıldı, port: 8080")
	if err := app.Listen("0.0.0.0:8080"); err != nil {
		logging.LogApp(logging.ERROR, "Sunucu başlatılamadı: %v", err)
	}
}
