.PHONY: build run test clean docker-build docker-up docker-down docker-logs help

# Varsayılan hedef
.DEFAULT_GOAL := help

# Değişkenler
BINARY_NAME=solar-sim
DOCKER_COMPOSE=docker compose

## build: Ana uygulamayı derle
build:
	go build -o $(BINARY_NAME) ./cmd/solar-sim

## run: Uygulamayı çalıştır
run:
	go run ./cmd/solar-sim

## test: Pattern test aracını çalıştır
test:
	go run ./cmd/checkpattern

## clean: Derleme çıktılarını temizle
clean:
	rm -f $(BINARY_NAME) checkpattern

## docker-build: Docker image'ı derle
docker-build:
	$(DOCKER_COMPOSE) build

## docker-up: Tüm servisleri başlat
docker-up:
	$(DOCKER_COMPOSE) up -d

## docker-down: Tüm servisleri durdur
docker-down:
	$(DOCKER_COMPOSE) down

## docker-logs: Container loglarını göster
docker-logs:
	$(DOCKER_COMPOSE) logs -f

## docker-restart: Servisleri yeniden başlat
docker-restart:
	$(DOCKER_COMPOSE) restart

## docker-ps: Container durumlarını göster
docker-ps:
	$(DOCKER_COMPOSE) ps

## deps: Bağımlılıkları indir
deps:
	go mod download
	go mod tidy

## help: Bu yardım mesajını göster
help:
	@echo "Kullanılabilir komutlar:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
