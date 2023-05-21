![image description](doc/screenshot.gif)

내 컴퓨터에서 돌리는 로컬 라마7B 실행기

## 뭐냐면요..

* `llama.cpp/examples/main`에 있는 코드 interactive 모드를 흉내낸 Go언어 바인딩
* 간단한 웹소켓 서버
* 간단한 Go언어 embed 웹ui


## 실행파일 다운로드
* [MS윈도우 cpu](https://github.com/edp1096/my-llama/releases/download/v0.1.10/my-llama_cpu.exe)
* [MS윈도우 clblast] - (https://github.com/edp1096/my-llama/releases/download/v0.1.10/my-llama_cl.zip)
    * NVIDIA CUDA SDK, AMD APP SDK, AMD ROCm, Intel OpenCL 중 하나가 필요합니다.


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
* GPU/CLBlast
    * [Go](https://golang.org/dl)
    * [MinGW>=12.2.0](https://github.com/brechtsanders/winlibs_mingw/releases/tag/12.2.0-16.0.0-10.0.0-ucrt-r5)
    * [Git](https://github.com/git-for-windows/git/releases)
    * [MS Visual Studio 2022 Community](https://visualstudio.microsoft.com/vs)
    * [Cmake >= 3.26](https://cmake.org/download)
    * [OpenCL-SDK](https://github.com/KhronosGroup/OpenCL-SDK), [CLBlast](https://github.com/KhronosGroup/OpenCL-SDK)
        * <b>빌드스크립트에 다운로드, 빌드 명령 포함되어있으므로, 수동으로 다운받을 필요 없습니다.</b>
        * 버전숫자를 바꾸려면 [build_cl.cmd](/build_cl.cmd)파일을 수정해주세요.
        * 추가로 NVIDIA CUDA SDK, AMD APP SDK, AMD ROCm, Intel OpenCL 중 하나가 필요합니다.
    * CPU Memory >= 12GB
    * Video Memory >= 4GB

### 컴파일
* 파워쉘 스크립트 - `ps1` 실행 안되면 아래와 같이 실행권한 설정해주세요.
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
* GPU/CLBLast
```powershell
git clone https://github.com/edp1096/my-llama.git

cd my-llama

git submodule update --init --recursive

build_cl.cmd
```


## 출처/참고
* https://github.com/ggerganov/llama.cpp
* https://github.com/go-skynet/go-llama.cpp
* https://github.com/cornelk/llama-go
