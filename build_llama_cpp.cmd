@echo off

cd llama.cpp
mkdir build
cd build

cmake .. -G "MinGW Makefiles" -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake --build .

copy bin\libllama.dll ..\..
cd ..\..
