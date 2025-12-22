# Docker ile Çalıştırma

## Hızlı Başlangıç

```bash
make docker-up

# veya
docker compose up -d
```

## Servisler

| Servis | Port | Açıklama |
|--------|------|----------|
| solar-sim-1 | 8081 | Simülatör instance 1 |
| solar-sim-2 | 8082 | Simülatör instance 2 |
| prometheus | 9090 | Metrik toplama |
| grafana | 3000 | Görselleştirme |

## Komutlar

```bash
make docker-up       # Başlat
make docker-down     # Durdur
make docker-logs     # Logları izle
make docker-build    # Yeniden derle
make docker-ps       # Container durumları
```

## Loglar

Container'lar tüm logları **stdout**'a yazar. Docker bu logları otomatik toplar.

```bash
# Tüm loglar
make docker-logs

# Belirli servis
docker compose logs -f solar-sim-1
```

## Volume'lar

| Volume | Kullanım |
|--------|----------|
| `prometheus_data` | Prometheus verileri |
| `grafana_data` | Grafana ayarları |

```bash
# Volume'ları koruyarak durdur
make docker-down

# Volume'ları da sil
docker compose down -v
```

## Ortam Değişkenleri

| Değişken | Varsayılan |
|----------|------------|
| `LOG_LEVEL` | `INFO` |
| `GRAFANA_USER` | `admin` |
| `GRAFANA_PASSWORD` | `admin` |
