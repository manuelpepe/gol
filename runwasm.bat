set GOOS=js
set GOARCH=wasm

go build -o dist/yourgame.wasm main.go

cd dist\
python -m http.server 8000
cd ..
