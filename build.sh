VERSION=1.0
GOOS_LIST=(linux)
GOARCH_LIST=(amd64 arm64 arm)

for GOOS in "${GOOS_LIST[@]}";
do
  for GOARCH in "${GOARCH_LIST[@]}";
  do
    export GOARCH=$GOARCH
    export GOOS=$GOOS
    file="ds389-exporter-$VERSION-$GOOS-$GOARCH"
    go build -ldflags "-s -w" -o $file ds389-exporter.go

    tar czvf $file.tar.gz $file
    rm $file
  done
done
