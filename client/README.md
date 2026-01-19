go run main.go

or

GOOS=js GOARCH=wasm go build -o game.wasm main.go
cp /usr/local/go/misc/wasm/wasm_exec.js .

then

python3 -m http.server 8080
http://localhost:8080/index.html
