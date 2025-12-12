all: go-build \
	docker-build \
	docker-save \
	docker-clean

go-build:
	CGO_ENABLED=0 go build -v \
	-buildvcs=false \
	-installsuffix 'static' \
	-ldflags="-X 'main.version=$$(git rev-parse --short HEAD)' -X 'main.build=$$(date --iso-8601=seconds)'" \
	-o ./dist/server \
	./cmd/server

docker-build:
	docker build -f ./cmd/server/Dockerfile \
	-t ChangerzaryX1602/sysware-test \
	--pull \
	.

docker-save:
	docker save ChangerzaryX1602/sysware-test | gzip > dist/ChangerzaryX1602-sysware-test.tar.gz

docker-clean:
	docker image prune -f

go-mod:
	go get -v -u && go mod tidy

go-test:
	go test -v ./...

go-run:
	go run ./cmd/server/

docker-start:
	colima start

docker-init-db:
	docker run -p 15432:5432 --name sysware-test-postgres -e POSTGRES_PASSWORD=SyswareTestPassword -d postgres:16-alpine

swagger-gen:
	swag init -g cmd/server/main.go

privkey-generate:
	openssl ecparam -name prime256v1 -genkey -noout -out internal/assets/dev/jwt/privkey.pem
pubkey-generate:
	openssl ec -in internal/assets/dev/jwt/privkey.pem -pubout -out internal/assets/dev/jwt/pubkey.pem