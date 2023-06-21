Llama 7B runner on my windows machine

## This is a ..

* Go binding for interactive mode of `llama.cpp/examples/main`
* `examples/runner`
    * Websocket server
    * Go embedded web ui


## Download pre-compiled binary, dll
* [DLL](https://github.com/edp1096/my-llama/releases)
* [MS-Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.20/my-llama_cpu.zip)
* [MS-Windows clblast](https://github.com/edp1096/my-llama/releases/download/v0.1.20/my-llama_cl.zip)
* [MS-Windows cuda](https://github.com/edp1096/my-llama/releases/download/v0.1.20/my-llama_cu.zip) - Require [CUDA toolkit 12](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64) or [this](https://github.com/ggerganov/llama.cpp/releases/download/master-66874d4/cudart-llama-bin-win-cu12.1.0-x64.zip)


## Usage

### Use this as go module
See <a href="https://pkg.go.dev/github.com/edp1096/my-llama"><img src="https://pkg.go.dev/badge/github.com/edp1096/my-llama.svg" alt="Go Reference"></a> or [examples](/examples).

* [Example](/examples/minimal/main.go)
```go
package main // import "minimal"

import (
	"fmt"

	llama "github.com/edp1096/my-llama"
)

func main() {
	modelName := "vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin"
	numPredict := 16

	l, err := llama.New()
	if err != nil {
		panic(err)
	}

	l.LlamaApiInitBackend()
	l.InitGptParams()

	l.SetNumThreads(4)
	l.SetUseMlock(true)
	l.SetNumPredict(numPredict)
	l.SetNumGpuLayers(32)
	l.SetSeed(42)

	l.InitContextParamsFromGptParams()

	err = l.LoadModel(modelName)
	if err != nil {
		panic(err)
	}

	l.AllocateTokens()

	numPast := 0
	prompt := "The quick brown fox"

	promptTokens, promptNumTokens := l.LlamaApiTokenize(prompt, true)
	fmt.Println("promptTokens:", promptTokens)

	if promptNumTokens < 1 {
		fmt.Println("numToken < 1")
		panic("numToken < 1")
	}

	isOK := l.LlamaApiEval(promptTokens, promptNumTokens, numPast)
	numPast += promptNumTokens

	fmt.Println("n_prompt_token, n_past, isOK:", promptNumTokens, numPast, isOK)
	fmt.Println("numPredict:", numPredict)

	for i := 0; i < numPredict; i++ {
		l.LlamaApiGetLogits()
		numVocab := l.LlamaApiNumVocab()

		l.PrepareCandidates(numVocab)
		nextToken := l.LlamaApiSampleToken()
		nextTokenStr := l.LlamaApiTokenToStr(nextToken)

		fmt.Print(nextTokenStr)
		l.LlamaApiEval([]int32{nextToken}, 1, numPast)

		numPast++
	}

	fmt.Println()

	l.LlamaApiFree()
}

/*
# CPU
$ go build [-tags cpu]

# GPU/CLBlast
$ go build -tags clblast

# GPU/CUDA
$ go build -tags cuda

# Before run, copy shared libraries(DLL) to folder where executable file exists

--------------------------------------------------

$ ./minimal
System Info: AVX = 1 | AVX2 = 1 | AVX512 = 0 | AVX512_VBMI = 0 | AVX512_VNNI = 0 | FMA = 1 | NEON = 0 | ARM_FMA = 0 | F16C = 1 | FP16_VA = 0 | WASM_SIMD = 0 | BLAS = 0 | SSE3 = 1 | VSX = 0 |
Model: vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin
llama.cpp: loading model from vicuna-7B-1.1-ggml_q4_0-ggjt_v3.bin
llama_model_load_internal: format     = ggjt v3 (latest)
llama_model_load_internal: n_vocab    = 32000
llama_model_load_internal: n_ctx      = 512
llama_model_load_internal: n_embd     = 4096
llama_model_load_internal: n_mult     = 256
llama_model_load_internal: n_head     = 32
llama_model_load_internal: n_layer    = 32
llama_model_load_internal: n_rot      = 128
llama_model_load_internal: ftype      = 2 (mostly Q4_0)
llama_model_load_internal: n_ff       = 11008
llama_model_load_internal: n_parts    = 1
llama_model_load_internal: model size = 7B
llama_model_load_internal: ggml ctx size =    0.07 MB
llama_model_load_internal: mem required  = 5407.71 MB (+ 1026.00 MB per state)
...................................................................................................
llama_init_from_file: kv self size  =  256.00 MB
promptTokens: [1 1576 4996 17354 1701 29916]
n_prompt_token, n_past, isOK: 6 6 true
numPredict: 16
 jumps over the lazy dog.
...
 */
```

### [Runner](/examples/runner/main.go) in [Release page](https://github.com/edp1096/my-llama/releases)
```powershell
# Just launch
my-llama.exe

# Launch with browser open
my-llama.exe -b
```
* When modified parameters in panel seem not working, try refresh the browser screen


## Build from source

### Requirements
* CPU
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * Memory >= 12GB
* GPU/CLBlast
    * Above CPU requirements and below
    * [OpenCL-SDK](https://github.com/KhronosGroup/OpenCL-SDK), [CLBlast](https://github.com/CNugteren/CLBlast)
        * <b>Build script download and build them automatically. No need to install manually</b>
        * If need change their version, just edit [build_lib.ps1](/build_lib.ps1).
        * And one of them
            * NVIDIA CUDA SDK
            * AMD APP SDK
            * AMD ROCm
            * Intel OpenCL
    * CPU Memory >= 12GB
    * Video Memory >= 6GB
* GPU/CUDA
    * Above CPU requirements and below
    * [CUDA Toolkit 12](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64)
    * CPU Memory >= 12GB
    * Video Memory >= 6GB

### Powershell scripts
* Before execute `ps1` script files, `ExecutionPolicy` should be set to `RemoteSigned` and unblock `ps1` files
```powershell
# Check
ExecutionPolicy
# Set as RemoteSigned
Set-ExecutionPolicy -Scope CurrentUser RemoteSigned

# Unblock ps1 files
Unblock-File *.ps1
```

### Clone repository then build library
* Clone
```powershell
git clone https://github.com/edp1096/my-llama.git
```
* Build
```powershell
# CPU
./build_lib.ps1

# GPU/CLBlast
./build_lib.ps1 clblast

# GPU/CUDA
./build_lib.ps1 cuda
```

* Clean temporary files
```powershell
./clean.ps1

# or

./clean.ps1 all
```

### Then build in examples folder
* Build runner in `examples/runner` folder
```powershell
cd examples/runner

# CPU
go build [-tags cpu]

# GPU/CLBlast
go build -tags clblast

# GPU/CUDA
go build -tags cuda
```

### Linux
* Build library
```sh
# CPU
./build_lib.sh

# GPU/CLBlast
./build_lib.sh clblast

# GPU/CUDA
./build_lib.sh cuda
```
* Clean temporary files
```sh
./clean.sh
```
* Tested with nVidia driver 530, CUDA toolkit 12.1
    * Ubuntu 20.04, RTX 1080ti
    * Ubuntu 20.04, RTX 3090
    * WSL Ubuntu 22.04, RTX 3060ti
* WSL
    * Because [not support opencl](https://github.com/microsoft/WSL/issues/6951), clblast not work
    * Set environment value `export GGML_CUDA_NO_PINNED=1` if CUDA not work


## Source
* https://github.com/ggerganov/llama.cpp
* https://github.com/go-skynet/go-llama.cpp
* https://github.com/cornelk/llama-go
