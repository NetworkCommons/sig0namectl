PROJECT := sig0

PORT := 8822

${PROJECT}.wasm: ../../golang/wasm/*.go ../../golang/sig0/*.go
	GOOS=js GOARCH=wasm go build -o ${PROJECT}.wasm ../../golang/wasm/*.go

start: ${PROJECT}.wasm
	nohup python3 -m http.server -d . ${PORT} &
	$(info open http://localhost:${PORT}/ in your browser)

stop:
	pkill -f "python3 -m http.server -d . ${PORT}"
