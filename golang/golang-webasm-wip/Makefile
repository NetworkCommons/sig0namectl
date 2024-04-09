PROJECT := golang-wasm-WIP

dnsfoo:
	make -C dnsfoo build
	make -C webtest start
${PROJECT}:
	GOOS=js GOARCH=wasm go build  -o webtest/${PROJECT}.wasm
	make -C webtest start
down:
	make -C webtest stop
