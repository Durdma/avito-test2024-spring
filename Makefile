.PHONY:
.SILENT:

build:
	docker-compose up app

swag:
	swag init -g internal/app/app.go

test:
	go test -v internal/controller/httpv1/banners_test.go