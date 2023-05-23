if ($args[0] -eq "clblast") {
    go build -trimpath -ldflags="-w -s -X main.deviceType=clblast" -o bin/ ./cmd
    cp -f llama.dll bin/
    cp -f openclblast/lib/clblast.dll bin/
} else {
    go build -trimpath -ldflags="-w -s" -o bin/ ./cmd
    cp -f llama.dll bin/
}
