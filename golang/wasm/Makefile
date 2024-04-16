PROJECT := wasm-demo

PORT := 8000

${PROJECT}.wasm: 
	GOOS=js GOARCH=wasm go build  -o ${PROJECT}.wasm

start: ${PROJECT}.wasm
	nohup python3 -m http.server -d . ${PORT} &

stop:
	pkill -f "python3 -m http.server -d . ${PORT}"
