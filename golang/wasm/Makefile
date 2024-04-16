PROJECT := demo

PORT := 8800

${PROJECT}.wasm: *.go ../sig0/*.go
	GOOS=js GOARCH=wasm go build  -o ${PROJECT}.wasm

start: ${PROJECT}.wasm
	nohup python3 -m http.server -d . ${PORT} &

stop:
	pkill -f "python3 -m http.server -d . ${PORT}"
