# My-llama

`My-llama` to run llama 7B on my windows machine.

## Requirements
* [Go](https://golang.org/dl)
* [MinGW](https://github.com/brechtsanders/winlibs_mingw)

## Build
```powershell
git submodule update --init --recursive

# Then modify /llama.cpp/Makefile. - See /Makefile
## Append 'CC = gcc'(not g++) at top of Makefile in /llama.cpp
## Comment or remove CCV, CXXV definition and print

mingw32-make.exe
```


## Usage
```powershell
# Download ggml weights
## https://huggingface.co/Drararara/llama-7b-ggml/tree/main
## https://huggingface.co/Pi3141/alpaca-lora-7B-ggml/tree/main
## https://huggingface.co/Sosaka/Vicuna-7B-4bit-ggml/tree/main
## https://huggingface.co/eachadea/ggml-vicuna-7b-4bit/tree/main

./bin/my-llama.exe -m <ggml_model_file>
```


## Todo
* Cleanup code
* Web server


## Source
* https://github.com/ggerganov/llama.cpp
* https://github.com/go-skynet/go-llama.cpp
* https://github.com/cornelk/llama-go
