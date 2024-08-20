.PHONY: build
build:
	go build
release:
	# clean
	go clean
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gchangelog-windows-amd64-${VERSION}.exe
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o gchangelog-linux-amd64-${VERSION}
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o gchangelog-darwin-amd64-${VERSION}
help:
	echo "run \"make release VERSION=$VERSION\" to build all binary"