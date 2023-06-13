![image description](doc/screenshot.gif)

Llama 7B runner on my windows machine

## This is a ..

* Go binding for interactive mode of `llama.cpp/examples/main`
* `cmd`
    * Websocket server
    * Go embedded web ui


## Download pre-compiled binary, dll
* [DLL](https://github.com/edp1096/my-llama/releases)
* [MS-Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.18/my-llama_cpu.zip)
* [MS-Windows clblast](https://github.com/edp1096/my-llama/releases/download/v0.1.18/my-llama_cl.zip)
* [MS-Windows cuda](https://github.com/edp1096/my-llama/releases/download/v0.1.18/my-llama_cu.zip) - Require [CUDA toolkit 12](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64) or [this](https://github.com/ggerganov/llama.cpp/releases/download/master-66874d4/cudart-llama-bin-win-cu12.1.0-x64.zip)


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
* GPU/CUDA
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

build_lib.ps1 cuda
```

* Clean
```powershell
clean.ps1

# or

clean.ps1 all
```

### Then build runner in cmd folder or in example folder
* CPU
```powershell
cd cmd
go build [-tags cpu]
```
* GPU/CLBlast
```powershell
cd cmd
go build -tags clblast
```
* GPU/CUDA
```powershell
cd cmd
go build -tags cuda
```

### Linux
* See [build_lib.sh](/build_lib.sh) and [clean.sh](/clean.sh)
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
