.PHONY: all
all:
	mkdir -p ./release && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./release/simpleotp_go main.go
