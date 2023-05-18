init:
	go mod download
	go mod vendor

run:
	go mod vendor;
	export APP_ID=GoCanal; go run main.go job --port=8080;

build-linux:
	go mod vendor;
	GOOS=linux GOARCH=amd64 go build -o "./build/go-canal" "./main.go";

build-windows:
	go mod vendor;
	GOOS=windows GOARCH=amd64 go build -o "./build/go-canal" "./main.go";

scp-dev:
	scp ./build/cruise-api root@测试api:/data/dianchi/

scp-pre:
	scp ./build/cruise-api root@正式api02:/data/dianchi_pre/

scp-prod:
	scp ./build/cruise-api root@店驰api01:/data/dianchi/

upx:
	upx --brute "./build/go-canal";
