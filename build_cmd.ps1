if ($args[0] -eq "clblast") {
    # export CGO_LDFLAGS="-static -L. -lstdc++ -lllama_cl -lbinding_cl"
    $env:CGO_LDFLAGS="-static -L. -lstdc++ -lllama_cl -lbinding_cl"

    go build -trimpath -ldflags="-w -s -X main.deviceType=clblast" -o bin/run-myllama_cl.exe ./cmd
    cp -f llama_cl.dll bin/
    cp -f openclblast/lib/clblast.dll bin/
} else {
    # export CGO_LDFLAGS="-static -L. -lstdc++ -lllama -lbinding"
    $env:CGO_LDFLAGS="-static -L. -lstdc++ -lllama -lbinding"

    go build -trimpath -ldflags="-w -s" -o bin/run-myllama.exe ./cmd
    cp -f llama.dll bin/
}
