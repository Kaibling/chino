dep-air:
	go install github.com/cosmtrek/air@latest	
run: dep-air
	air
docker-start:
	docker-compose -f build/docker-compose.yml up -d

build-server:
	go build -v -ldflags="-X 'main.appVersion=v0.0.0' -X 'main.goVersion=$(go version)' -X 'main.buildTime=$(date -u '+%Y-%m-%d_%H-%M-%S')'" -o chino  -buildvcs=false 
