![image description](doc/screenshot.gif)

Llama 7B runner on my windows machine

## This is a ..

* Go binding for interactive mode of `llama.cpp/examples/main`
* Websocket server
* Go embedded web ui


## Download pre-compiled binary
* [MS-Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.7/my-llama_cpu.exe)
* [MS-Windows cuda](https://github.com/edp1096/my-llama/releases/download/v0.1.7/my-llama_cu.zip) - require [CUDA Toolkit 12.1](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64) or [DLLs](https://github.com/ggerganov/llama.cpp/releases/download/master-e6a46b0/cudart-llama-bin-win-cu12.1.0-x64.zip)
* [MS-Windows clblast](https://github.com/edp1096/my-llama/releases/download/v0.1.7/my-llama_cl.zip)


## Usage
```powershell
# Just launch
./bin/my-llama.exe

# Launch with browser open
./bin/my-llama.exe -b
```
* When modified parameters in panel seem not working, try refresh the browser screen


## Build from source

### Requirements
* CPU
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * Memory >= 12GB
* GPU/CUDA
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [CUDA Toolkit 12.1](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64)
    * CPU Memory >= 12GB
    * Video Memory >= 4GB
* GPU/CLBlast
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [OpenCL-SDK](https://github.com/KhronosGroup/OpenCL-SDK), [CLBlast](https://github.com/KhronosGroup/OpenCL-SDK)
        * <b>When build script running, download and build them automatically. No need to install manually</b>
        * If need change their version, just edit [build_cl.cmd](/build_cl.cmd).
    * CPU Memory >= 12GB
    * Video Memory >= 4GB

### Build
* CPU
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

git submodule update --init --recursive

mingw32-make.exe
```
* GPU/CUDA
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

git submodule update --init --recursive

build_cu.cmd
```
* GPU/CLBlast
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

git submodule update --init --recursive

build_cl.cmd
```

### Use binding
See <a href="https://pkg.go.dev/github.com/edp1096/my-llama/cgollama"><img src="https://pkg.go.dev/badge/github.com/edp1096/my-llama/cgollama.svg" alt="Go Reference"></a> and [`main.go`](main.go)


## Todo
* [ ] antiprompt, response name
    * [x] Add response name input - `### Assistant:`
    * [ ] Parse `### Human:`, `### Assistant:`
* [ ] Add Papago, Kakao, DeepL translator
* `binding.cpp`
    * GGML Parameter settings - Set parameters from html to websocket server
        * Not touch. Probably I can't
            * ~~n_predict - new tokens to predict~~
            * ~~seed, n_keep, f16_kv, use_mmap, use_mlock~~
    * Not touch. Probably I can't
        * ~~Clean up functions~~
        * ~~Remove and integrate all unnecessary functions~~


## Source
* Code
    * https://github.com/ggerganov/llama.cpp
    * https://github.com/go-skynet/go-llama.cpp
    * https://github.com/cornelk/llama-go
* Prompt
    * https://arca.live/b/alpaca/73449389
    ```dos
    main -m ggml-vicuna-7b-4bit-rev1.bin --color -f ./prompts/vicuna.txt -i --n_parts 1 -t 6 --temp 0.15 --top_k 400 -c 2048 --repeat_last_n 2048 --repeat_penalty 1.0 -n 2048 -r "### Human:" -b 512
    ```
