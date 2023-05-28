#!/bin/sh

# common
# llama.cpp compile error bdbda1b1 so always copy
cp -f llama.cpp_deallocate/* vendors/llama.cpp/

if [ "$1" = "cpu" ] || [ -z "$1" ]; then
    # cpu
    cd vendors/llama.cpp
    mkdir -p build; cd build

    cmake .. -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
    cmake --build . --config Release

    cp libllama.a ../../../libllama_lin64.a

    cd ../../..

    g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
    ar src libbinding_lin64.a libllama_lin64.a binding.o binding_llama_api.o

elif [ "$1" = "clblast" ]; then
    # clblast
    mkdir -p ./openclblast; cd openclblast

    rm -rf ./OpenCL-SDK
    git clone --recurse-submodules https://github.com/KhronosGroup/OpenCL-SDK.git
    mkdir -p OpenCL-SDK/build; cd OpenCL-SDK/build
    cmake .. -DBUILD_DOCS=OFF -DBUILD_TESTING=OFF -DBUILD_EXAMPLES=OFF -DOPENCL_SDK_BUILD_SAMPLES=OFF -DOPENCL_SDK_TEST_SAMPLES=OFF
    cmake --build . --config Release
    cmake --install . --prefix ../..

    cd ../..

    rm -rf ./CLBlast
    git clone https://github.com/CNugteren/CLBlast.git
    mkdir -p CLBlast/build; cd CLBlast/build
    cmake .. -DBUILD_SHARED_LIBS=OFF -DTUNERS=OFF
    cmake --build . --config Release
    cmake --install . --prefix ../..

    cd ../..

    cd ..

    cd vendors/llama.cpp
    mkdir -p build; cd build

    curl -L

    cmake .. -DLLAMA_CLBLAST=1 -DCMAKE_PREFIX_PATH="../../openclblast" -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
    cmake --build . --config Release

    cp libllama.a ../../../libllama_cl_lin64.a

    cd ../../..

    g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
    ar src libbinding_cl_lin64.a libllama_cl_lin64.a binding.o

elif [ "$1" = "cuda" ]; then
    # cuda
    cd vendors/llama.cpp
    mkdir -p build
    cd build

    cmake .. -DLLAMA_CUBLAS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
    cmake --build . --config Release

    cp libllama.a ../../../libllama_cu_lin64.a

    cd ../../..

    g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
    ar src libbinding_cu_lin64.a libllama_cu_lin64.a binding.o

else
    echo "Invalid argument"

fi
