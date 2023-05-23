![image description](doc/screenshot.gif)

Llama 7B runner on my windows machine

## This is a ..

* Go binding for interactive mode of `llama.cpp/examples/main`
* Websocket server
* Go embedded web ui


## Download pre-compiled binary
* [MS-Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.10/my-llama_cpu.exe)
* [MS-Windows clblast](https://github.com/edp1096/my-llama/releases/download/v0.1.10/my-llama_cl.zip)
    * Require one of them - NVIDIA CUDA SDK or AMD APP SDK or AMD ROCm or Intel OpenCL


## Usage
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
    * Memory >= 12GB
* GPU/CLBlast
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [OpenCL-SDK](https://github.com/KhronosGroup/OpenCL-SDK), [CLBlast](https://github.com/CNugteren/CLBlast)
        * <b>When build script running, download and build them automatically. No need to install manually</b>
        * If need change their version, just edit [build_cl.cmd](/build_cl.cmd).
        * And one of them - NVIDIA CUDA SDK or AMD APP SDK or AMD ROCm or Intel OpenCL
    * CPU Memory >= 12GB
    * Video Memory >= 4GB

### Build
* Scripts - `ExecutionPolicy` should be set to `RemoteSigned` and unblock `ps1` files
```powershell
# Check
ExecutionPolicy
# Set as RemoteSigned
Set-ExecutionPolicy -Scope CurrentUser RemoteSigned

# Unblock ps1 files
Unblock-File *.ps1
```

* CPU
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

git submodule update --init --recursive

mingw32-make.exe
```
* GPU/CLBlast
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

git submodule update --init --recursive

build_cl.ps1
```

### Use binding
See <a href="https://pkg.go.dev/github.com/edp1096/my-llama/cgollama"><img src="https://pkg.go.dev/badge/github.com/edp1096/my-llama/cgollama.svg" alt="Go Reference"></a> and [`main.go`](main.go)


## Todo
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
* https://github.com/ggerganov/llama.cpp
* https://github.com/go-skynet/go-llama.cpp
* https://github.com/cornelk/llama-go
