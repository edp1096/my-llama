![image description](doc/screenshot.gif)

Llama 7B runner on my windows machine

## This is a ..

* Go binding for interactive mode of `llama.cpp/examples/main`
* Websocket server
* Go embedded web ui


## Download pre-compiled binary
* ggjt v2 (GGML new)
    * [MS-Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.8/my-llama_cpu.exe)
    * [MS-Windows cuda](https://github.com/edp1096/my-llama/releases/download/v0.1.8/my-llama_cu.zip) - require [CUDA Toolkit 12](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64) or [DLLs](https://github.com/ggerganov/llama.cpp/releases/download/master-e6a46b0/cudart-llama-bin-win-cu12.1.0-x64.zip) and VRAM >= 7GB
    * [MS-Windows clblast](https://github.com/edp1096/my-llama/releases/download/v0.1.8/my-llama_cl.zip)
* ggjt v1 (GGML old)
    * [MS-Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.8/my-llama_cpu_old_ggml.exe)
    * [MS-Windows cuda](https://github.com/edp1096/my-llama/releases/download/v0.1.8/my-llama_cu_old_ggml.zip) - require [CUDA Toolkit 12](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64) or [DLLs](https://github.com/ggerganov/llama.cpp/releases/download/master-e6a46b0/cudart-llama-bin-win-cu12.1.0-x64.zip)
    * [MS-Windows clblast](https://github.com/edp1096/my-llama/releases/download/v0.1.8/my-llama_cl_old_ggml.zip)


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
    * [CUDA Toolkit 12](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64)
    * CPU Memory >= 12GB
    * Video Memory >= 7GB
        * Beacuse of GPU token generation, hard coded `n_gpu_layer` as 32 for my 3060ti
        * If you have different GPU, you may need to change it in [cgollama/binding.cpp](/cgollama/binding.cpp)
* GPU/CLBlast
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [OpenCL-SDK](https://github.com/KhronosGroup/OpenCL-SDK), [CLBlast](https://github.com/CNugteren/CLBlast)
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
