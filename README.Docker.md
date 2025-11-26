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

**İlk kullanımda:**
- Grafana'ya giriş yapınca **Solar Simulator Dashboard** otomatik yüklenmiş olacak
- Prometheus datasource otomatik yapılandırılmıştır, ekstra bir ayar gerekmez
- Dashboard'da 6 panel göreceksiniz: Power Output, Voltage, Current, Efficiency, Temperature, ve Gauge

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

### Hızlı Başlangıç

1. http://localhost:3000 adresine gidin
2. admin/admin ile giriş yapın (ilk girişte şifre değiştirmeniz istenebilir)
3. Dashboard otomatik olarak yüklenmiş olacak: **Solar Simulator Dashboard**

### Dashboard'ların Kalıcılığı

Bu projede Grafana dashboard'ları **iki şekilde** korunur:

#### 1. Named Volume (grafana_data)
- UI üzerinde yaptığınız tüm değişiklikler `/var/lib/grafana` içinde saklanır
- `docker-compose down` ve `docker-compose up` yaptığınızda veriler korunur
- ⚠️ **DİKKAT**: `docker-compose down -v` komutu volume'ları siler, verilerinizi kaybedersiniz!

#### 2. Provisioning (Versiyon Kontrolü)
- `grafana/dashboards/` klasöründeki JSON dosyaları otomatik yüklenir
- Repo ile birlikte gelir, değişiklikler git ile takip edilir
- Container her başlatıldığında bu dashboard'lar otomatik import edilir
- `grafana/provisioning/` altında datasource ve dashboard yapılandırmaları vardır

### Dashboard Yedekleme

Mevcut Grafana verilerini yedeklemek için:

```bash
# SQLite veritabanını yedekle (tüm dashboard'lar, kullanıcılar, ayarlar)
docker cp grafana:/var/lib/grafana/grafana.db ./grafana-backup-$(date +%Y%m%d).db

# Alternatif: Volume içeriğini tar olarak yedekle
docker run --rm -v prom-custom-metric_grafana_data:/data -v "$(pwd)":/backup alpine tar czf /backup/grafana-data-backup.tar.gz -C /data .
```

### Dashboard'ları Git'e Kaydetme

UI üzerinde oluşturduğunuz yeni dashboard'ları repoya eklemek için:

1. Grafana UI'de dashboard'unuzu açın
2. Sağ üstteki Share butonuna tıklayın
3. **Export** sekmesine geçin
4. **Save to file** ile JSON'u indirin
5. İndirdiğiniz dosyayı `grafana/dashboards/` klasörüne kopyalayın
6. Git ile commit edin:

```bash
cp ~/Downloads/my-dashboard.json grafana/dashboards/
git add grafana/dashboards/my-dashboard.json
git commit -m "Add new dashboard"
```

Container yeniden başlatıldığında bu dashboard otomatik yüklenecektir.

### Güvenli Volume Yönetimi

```bash
# GÜVENLİ: Volume'ları koruyarak durdur
docker-compose down

# TEHLİKELİ: Volume'ları da siler!
docker-compose down -v

# Volume'ları listeleme
docker volume ls | grep grafana

# Belirli volume'u silme (dikkatli kullanın!)
docker volume rm prom-custom-metric_grafana_data
```

### Datasource Yapılandırması

Prometheus datasource otomatik olarak yapılandırılır (`grafana/provisioning/datasources/datasource.yml`):
- Name: Prometheus
- URL: http://prometheus:9090
- Default: Yes

Farklı bir datasource eklemek için `grafana/provisioning/datasources/` altına yeni YAML dosyası ekleyin.

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
