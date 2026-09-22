CGO_ENABLED=0 go build -a -ldflags="-s -w" -gcflags=all="-l -B -C" -o ./bin/client pkg/client/main.go
CGO_ENABLED=0 go build -a -ldflags="-s -w" -gcflags=all="-l -B -C" -o ./bin/server pkg/server/main.go
ls -lah ./bin