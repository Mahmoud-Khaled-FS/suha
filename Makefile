build:
	@go build -o suha.exe main.go

run: build
	suha $(arg)