![image description](doc/screenshot.gif)

Llama 7B runner on my windows machine

## This is a ..

* Go binding for interactive mode of `llama.cpp/examples/main`
* `cmd`
    * Websocket server
    * Go embedded web ui


## Download pre-compiled binary, dll
* [DLL](https://github.com/edp1096/my-llama/releases)
* [MS-Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.15/my-llama_cpu.zip)
* [MS-Windows clblast](https://github.com/edp1096/my-llama/releases/download/v0.1.15/my-llama_cl.zip)


## Usage

### Use this as go module
See <a href="https://pkg.go.dev/github.com/edp1096/my-llama"><img src="https://pkg.go.dev/badge/github.com/edp1096/my-llama.svg" alt="Go Reference"></a> or [my-llama-app](https://github.com/edp1096/my-llama-app) repo.

### runner in cmd
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
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [OpenCL-SDK](https://github.com/KhronosGroup/OpenCL-SDK), [CLBlast](https://github.com/CNugteren/CLBlast)
        * <b>When build script running, download and build them automatically. No need to install manually</b>
        * If need change their version, just edit [build_lib.ps1](/build_lib.ps1).
        * And one of them
            * NVIDIA CUDA SDK
            * AMD APP SDK
            * AMD ROCm
            * Intel OpenCL
    * CPU Memory >= 12GB
    * Video Memory >= 4GB

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
* CPU
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

build_lib.ps1
```
* GPU/CLBlast
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

build_lib.ps1 clblast
```

### Then build runner in cmd folder
* CPU
```powershell
build_cmd.ps1

# or

cd cmd
go build -tags cpu
```
* GPU/CLBlast
```powershell
build_cmd.ps1 clblast

# or

cd cmd
go build -tags clblast
```


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
