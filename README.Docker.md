# Docker ile Çalıştırma

Bu proje Docker ve Docker Compose kullanarak kolayca çalıştırılabilir.

## Gereksinimler

- Docker
- Docker Compose

## Kurulum ve Çalıştırma

### 1. Projeyi klonlayın veya indirin

```bash
git clone <repo-url>
cd prom-custom-metric
```

### 2. Ortam değişkenlerini ayarlayın (Opsiyonel)

```bash
cp .env.example .env
# .env dosyasını düzenleyerek özel ayarlarınızı yapın
```

### 3. Docker Compose ile çalıştırın

```bash
# Tüm servisleri başlat (solar-sim, prometheus, grafana)
docker-compose up -d

# Sadece uygulamayı başlat
docker-compose up -d solar-sim

# Logları görüntüle
docker-compose logs -f

# Belirli bir servisin logları
docker-compose logs -f solar-sim
```

### 4. Servislere erişim

- **Solar Simulator Metrics**: http://localhost:8080/metrics
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (kullanıcı: admin, şifre: admin)

## Docker Komutları

```bash
# Servisleri durdur
docker-compose down

# Servisleri durdur ve volume'leri sil
docker-compose down -v

# Servisleri yeniden başlat
docker-compose restart

# Image'ı yeniden build et
docker-compose build

# Build edip başlat
docker-compose up -d --build

# Servislerin durumunu kontrol et
docker-compose ps
```

## Sadece Uygulamayı Docker ile Çalıştırma

Eğer sadece solar simulator'ı çalıştırmak istiyorsanız:

```bash
# Image'ı build et
docker build -t solar-sim .

# Çalıştır
docker run -d -p 8080:8080 --name solar-sim solar-sim

# Logları görüntüle
docker logs -f solar-sim

# Durdur ve sil
docker stop solar-sim
docker rm solar-sim
```

## Prometheus Konfigürasyonu

Prometheus konfigürasyonunu değiştirmek için `prometheus.yml` dosyasını düzenleyin ve servisi yeniden başlatın:

```bash
docker-compose restart prometheus
```

## Grafana Konfigürasyonu

1. http://localhost:3000 adresine gidin
2. admin/admin ile giriş yapın
3. Data Source ekleyin:
   - Type: Prometheus
   - URL: http://prometheus:9090
4. Dashboard oluşturun ve `mppt_values` metriklerini kullanın

## Sorun Giderme

### Port zaten kullanımda hatası

```bash
# Kullanılan portları kontrol et
sudo netstat -tulpn | grep :8080
sudo netstat -tulpn | grep :9090
sudo netstat -tulpn | grep :3000

# docker-compose.yml dosyasında portları değiştir
```

### Logları detaylı görmek

```bash
# Tüm servisler
docker-compose logs -f

# Sadece solar-sim
docker-compose logs -f solar-sim

# Son 100 satır
docker-compose logs --tail=100 solar-sim
```

### Container'a bağlan

```bash
docker exec -it solar-sim sh
```

### Volume'leri temizle

```bash
docker-compose down -v
docker volume prune
```
