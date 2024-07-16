GOOS=linux # REQUIRED REPLACE
GOARCH=amd64 # REQUIRED REPLACE

export GOARCH=$GOARCH
export GOOS=$GOOS
file="ds389-exporter"
go build -ldflags "-s -w" -o $file ds389-exporter.go
