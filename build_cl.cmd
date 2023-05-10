@echo off

curl --progress-bar -Lo clblast.zip "https://github.com/CNugteren/CLBlast/releases/download/1.5.3/CLBlast-1.5.3-Windows-x64.zip"
curl --progress-bar -Lo opencl.zip "https://github.com/KhronosGroup/OpenCL-SDK/releases/download/v2023.04.17/OpenCL-SDK-v2023.04.17-Win-x64.zip"

md openclblast
tar -xf clblast.zip -C openclblast
tar -xf opencl.zip -C openclblast

xcopy openclblast\OpenCL-SDK-v2023.04.17-Win-x64 openclblast /E /Y
copy openclblast_cmake\*.cmake openclblast\lib\cmake\CLBlast /Y >nul 2>&1

rmdir openclblast\OpenCL-SDK-v2023.04.17-Win-x64 /s /q >nul 2>&1


cd llama.cpp

md build 2>nul
cd build

cmake .. -DCMAKE_PREFIX_PATH='../openclblast' -DLLAMA_CLBLAST=1 -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=1 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

copy bin\Release\llama.dll ..\..

cd ..\..

gendef.exe llama.dll
dlltool.exe -k -d llama.def -l libllama.a

mingw32-make.exe build_for_cuda

copy llama.dll .\bin /y
copy openclblast\lib\clblast.dll .\bin /y
