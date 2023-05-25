if ($args[0] -eq "clblast") {
    go build -tags clblast -trimpath -ldflags="-w -s" -o bin/run-myllama_cl.exe ./cmd
    cp -f llama_cl.dll bin/
    cp -f openclblast/lib/clblast.dll bin/
} else {
    go build -tags cpu -trimpath -ldflags="-w -s" -o bin/run-myllama_cpu.exe ./cmd
    cp -f llama.dll bin/
}
