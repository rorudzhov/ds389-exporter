GOOS=linux # REQUIRED REPLACE
GOARCH=arm64 # REQUIRED REPLACE

export GOARCH=$GOARCH
export GOOS=$GOOS
file="ds389-exporter-$GOOS-$GOARCH"
go build -ldflags "-s -w" -o $file ds389-exporter.go

tar czvf $file.tar.gz $file
rm $file

