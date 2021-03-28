all: linux-amd64 linux-386 darwin-amd64 windows-amd64 windows-386

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/runai-container-runtime ./cmd

linux-386:
	GOOS=linux GOARCH=386 go build -o bin/linux/386/runai-container-runtime ./cmd

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/runai-container-runtime ./cmd

windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/x86_64/runai-container-runtime ./cmd

windows-386:
	GOOS=windows GOARCH=386 go build -o bin/windows/x86/runai-container-runtime ./cmd
