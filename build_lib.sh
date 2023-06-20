#!/bin/sh

if [ "$1" = "cpu" ] || [ -z "$1" ]; then
    # cpu
    cd llama.cpp

    cmake -B ../build . -DLLAMA_BUILD_EXAMPLES="OFF" -DLLAMA_BUILD_TESTS="OFF"
    cd ../build
    cmake --build . --config Release

    cp libllama.a ../libllama_lin64.a

    cd ..

    g++ \
    -static -O3 -std=c++11 -fPIC -march=native -mtune=native \
    -I./llama.cpp -I./llama.cpp/examples \
    myllama.cpp -o myllama.o -c
    g++ \
    -static -O3 -std=c++11 -fPIC -march=native -mtune=native \
    -I./llama.cpp -I./llama.cpp/examples \
    myllama_llama_api.cpp -o myllama_llama_api.o -c

    ar src libmyllama_lin64.a libllama_lin64.a myllama.o myllama_llama_api.o

elif [ "$1" = "clblast" ]; then
    # clblast
    export LD_LIBRARY_PATH=$(pwd):$LD_LIBRARY_PATH

    cd vendors
    mkdir -p ./openclblast; cd openclblast

    rm -rf ./OpenCL-SDK
    rm -rf ./CLBlast

    git clone --recurse-submodules https://github.com/KhronosGroup/OpenCL-SDK.git
    git clone https://github.com/CNugteren/CLBlast.git

    mkdir -p OpenCL-SDK/build; cd OpenCL-SDK/build
    cmake .. -DBUILD_DOCS=OFF -DBUILD_TESTING=OFF -DBUILD_EXAMPLES=OFF -DOPENCL_SDK_BUILD_SAMPLES=OFF -DOPENCL_SDK_TEST_SAMPLES=OFF
    cmake --build . --config Release
    cmake --install . --prefix ../..

    cd ../..

    mkdir -p CLBlast/build; cd CLBlast/build
    cmake .. -DBUILD_SHARED_LIBS=OFF -DTUNERS=OFF
    cmake --build . --config Release
    cmake --install . --prefix ../..

    cd ../..

    ar src ../../libOpenCL_lin64.a ./lib/libOpenCL.so 
    cp -f ./lib/libclblast.a ../../libclblast_lin64.a

    cd ../..


    cd llama.cpp

    cmake -B ../build . -DLLAMA_CLBLAST="ON" -DCMAKE_PREFIX_PATH="../vendors/openclblast" -DLLAMA_BUILD_EXAMPLES="OFF" -DLLAMA_BUILD_TESTS="OFF"
    cd ../build
    cmake --build . --config Release

    cp libllama.a ../libllama_cl_lin64.a

    cd ..

    g++ \
    -static -O3 -std=c++11 -fPIC -march=native -mtune=native \
    -I./llama.cpp -I./llama.cpp/examples \
    myllama.cpp -o myllama.o -c
    g++ \
    -static -O3 -std=c++11 -fPIC -march=native -mtune=native \
    -I./llama.cpp -I./llama.cpp/examples \
    myllama_llama_api.cpp -o myllama_llama_api.o -c

    ar src libmyllama_cl_lin64.a libllama_cl_lin64.a myllama.o myllama_llama_api.o

elif [ "$1" = "cuda" ]; then
    # cuda
    cd llama.cpp

    cmake -B ../build . -DLLAMA_CUBLAS="ON" -DLLAMA_BUILD_EXAMPLES="OFF" -DLLAMA_BUILD_TESTS="OFF"
    cd ../build
    cmake --build . --config Release

    cp libllama.a ../libllama_cu_lin64.a

    cd ..

    g++ \
    -static -O3 -std=c++11 -fPIC -march=native -mtune=native \
    -I./llama.cpp -I./llama.cpp/examples \
    myllama.cpp -o myllama.o -c
    g++ \
    -static -O3 -std=c++11 -fPIC -march=native -mtune=native \
    -I./llama.cpp -I./llama.cpp/examples \
    myllama_llama_api.cpp -o myllama_llama_api.o -c
    
    ar src libmyllama_cu_lin64.a libllama_cu_lin64.a myllama.o myllama_llama_api.o

else
    echo "Invalid argument"

fi
