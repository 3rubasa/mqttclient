.PHONY: default all

local: build_raspi deploy_local
default: build_raspy deploy_local

build_linux_amd64:
	go env -w GOOS=linux && go env -w GOARCH=amd64 && go build -o bin/mqttclient

build_raspi:
	go env -w GOOS=linux && go env -w GOARCH=arm64 && go build -o bin/mqttclient

deploy_local:
	plink -pw p dima@10.42.0.1 "sudo systemctl stop mqttclient.service" && pscp -pw p ./bin/mqttclient dima@10.42.0.1:/opt/mqttclient && plink -pw p dima@10.42.0.1 "sudo systemctl start mqttclient.service"