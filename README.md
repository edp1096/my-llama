![image description](doc/screenshot.png)

Llama 7B runner on my windows machine


## Download pre-compiled binary
* [Windows](https://github.com/edp1096/my-llama/releases/download/v0.0.2/my-llama.exe)


## Build from source

### Requirements
* [Go](https://golang.org/dl)
* [MinGW](https://github.com/brechtsanders/winlibs_mingw)

### Build
```powershell
git clone https://github.com/edp1096/my-llama.git

git submodule update --init --recursive

mingw32-make.exe
```
* About submodule [llama.cpp](https://github.com/ggerganov/llama.cpp), since `<regex>` header is removed and came huge changes from the [commits beyond aaf3b23](https://github.com/ggerganov/llama.cpp/commit/f963b63afa0e057cfb9eba4d88407c6a0850a0d8), you should append `<time.h>` to `llama.cpp/llama.cpp` manually and also should modify many things in `cgollama/*`. Otherwise, do keep the commit hash as [aaf3b23](https://github.com/ggerganov/llama.cpp/commit/aaf3b23debc1fe1a06733c8c6468fb84233cc44f).


## Usage
```powershell
# Download ggml weights
## https://huggingface.co/Drararara/llama-7b-ggml/tree/main
## https://huggingface.co/Pi3141/alpaca-lora-7B-ggml/tree/main
## https://huggingface.co/Sosaka/Vicuna-7B-4bit-ggml/tree/main
## https://huggingface.co/eachadea/ggml-vicuna-7b-4bit/tree/main

./bin/my-llama.exe [-m <ggml_model_file>] [-t <cpu_count>] [-n <token_count>]
```


## TODO
* Prompt
* Add Papago, Kakao API


## Source
* https://github.com/ggerganov/llama.cpp
* https://github.com/go-skynet/go-llama.cpp
* https://github.com/cornelk/llama-go
