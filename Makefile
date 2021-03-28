all: linux-amd64 linux-386 darwin-amd64 windows-amd64 windows-386

linux-amd64: runtime-linux-amd64 hook-linux-amd64

runtime-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/runtime/linux/amd64/runai-container-runtime ./cmd/runtime

hook-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/hook/linux/amd64/runai-container-runtime-hook ./cmd/hook


linux-386: runtime-linux-386 hook-linux-386

runtime-linux-386:
	GOOS=linux GOARCH=386 go build -o bin/runtime/linux/386/runai-container-runtime ./cmd/runtime

hook-linux-386:
	GOOS=linux GOARCH=386 go build -o bin/hook/linux/386/runai-container-runtime-hook ./cmd/hook


darwin-amd64: runtime-darwin-amd64 hook-darwin-amd64

runtime-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/runtime/darwin/amd64/runai-container-runtime ./cmd/runtime

hook-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/hook/darwin/amd64/runai-container-runtime-hook ./cmd/hook


windows-amd64: runtime-windows-amd64 hook-windows-amd64

runtime-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/runtime/windows/x86_64/runai-container-runtime ./cmd/runtime

hook-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/hook/windows/x86_64/runai-container-runtime-hook ./cmd/hook


windows-386: runtime-windows-386 hook-windows-386

runtime-windows-386:
	GOOS=windows GOARCH=386 go build -o bin/runtime/windows/x86/runai-container-runtime ./cmd/runtime

hook-windows-386:
	GOOS=windows GOARCH=386 go build -o bin/hook/windows/x86/runai-container-runtime-hook ./cmd/hook
