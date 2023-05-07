![image description](doc/screenshot.gif)

내 컴퓨터에서 돌리는 로컬 라마7B 실행기

## 뭐냐면요..

* `llama.cpp/examples/main`에 있는 코드 interactive 모드를 흉내낸 Go언어 바인딩
* 간단한 웹소켓 서버
* 간단한 embed 웹ui


## 실행파일 다운로드
* [MS윈도우 cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.4/my-llama_cpu.exe)
* [MS윈도우 cuda](https://github.com/edp1096/my-llama/releases/download/v0.1.4/my-llama_cu.zip) - [CUDA Toolkit 12.1](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64)를 설치해야됩니다


## 실행 방법
```powershell
# 실행만
./bin/my-llama.exe

# 실행하면서 웹브라우저 같이 띄우기
./bin/my-llama.exe -b
```
* 파라미터가 말을 안듣는 것 같으면 브라우저 새로고침 하세요

## 소스 빌드하기

### 요구사항
* CPU
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * Memory >= 12GB
* GPU
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [CUDA Toolkit 12.1](https://developer.nvidia.com/cuda-downloads?target_os=Windows&target_arch=x86_64)
    * CPU Memory >= 12GB
    * Video Memory >= 4GB

### 컴파일
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

build_cu.cmd
```


## 출처/참고
* 코드
    * https://github.com/ggerganov/llama.cpp
    * https://github.com/go-skynet/go-llama.cpp
    * https://github.com/cornelk/llama-go
* 프롬프트
    * https://arca.live/b/alpaca/73449389
    ```dos
    main -m ggml-vicuna-7b-4bit-rev1.bin --color -f ./prompts/vicuna.txt -i --n_parts 1 -t 6 --temp 0.15 --top_k 400 -c 2048 --repeat_last_n 2048 --repeat_penalty 1.0 -n 2048 -r "### Human:" -b 512
    ```
