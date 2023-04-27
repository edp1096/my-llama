![image description](doc/screenshot.gif)

Llama 7B runner on my windows machine

## This is a ..

* Go binding for interactive mode of `llama.cpp/examples/main`
* Websocket server
* Go embedded web ui


## Download pre-compiled binary
* [Windows cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.2/my-llama.exe)
* [Windows cuda](https://github.com/edp1096/my-llama/releases/download/v0.1.2/my-llama_cu.zip) - require [CUDA Toolkit 12.1](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64)


## Usage
```powershell
# Just launch
./bin/my-llama.exe

# Launch with browser open
./bin/my-llama.exe -b
```


## Build from source

### Requirements
* CPU
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * Memory >= 12GB
* GPU
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [Ms Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [CUDA Toolkit 12.1](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64)
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
* GPU
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

git submodule update --init --recursive

build.cmd
```


## Todo
* `binding.cpp`
    * [ ] Send separated input instead of hard coded prompt when begin in html
    * [ ] GGML Parameter settings - Set parameters from html to websocket server
        * [ ] Slider
        * [x] topk, topp, temp, repeat_penalty
        * [ ] Others - ?
    * [ ] Clean up functions & crash predict function in cgollama.go
    * [ ] Remove and integrate all unnecessary functions
    * [v] Replace all `llama_` prefix to `bd_`
    * [v] Save light/dark mode to localstorage
* [ ] Save/Load state
* [ ] Add Papago, Kakao translator


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
