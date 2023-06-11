#!/bin/sh

# common
# llama.cpp compile error bdbda1b1 so always copy
# cp -f llama.cpp_mod/* vendors/llama.cpp/

# if [ "$1" = "cpu" ] || [ -z "$1" ]; then
#     # cpu
#     cd vendors/llama.cpp
#     mkdir -p build; cd build

#     cmake .. -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
#     cmake --build . --config Release

#     cp libllama.a ../../../libllama_lin64.a

#     cd ../../..

#     g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
#     g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples myllama_llama_api.cpp -o myllama_llama_api.o -c
#     ar src libbinding_lin64.a libllama_lin64.a binding.o myllama_llama_api.o

# elif [ "$1" = "clblast" ]; then
#     # clblast
#     export LD_LIBRARY_PATH=$(pwd):$LD_LIBRARY_PATH

#     mkdir -p ./openclblast; cd openclblast

#     rm -rf ./OpenCL-SDK
#     git clone --recurse-submodules https://github.com/KhronosGroup/OpenCL-SDK.git
#     mkdir -p OpenCL-SDK/build; cd OpenCL-SDK/build
#     cmake .. -DBUILD_DOCS=OFF -DBUILD_TESTING=OFF -DBUILD_EXAMPLES=OFF -DOPENCL_SDK_BUILD_SAMPLES=OFF -DOPENCL_SDK_TEST_SAMPLES=OFF
#     cmake --build . --config Release
#     cmake --install . --prefix ../..

#     cd ../..
#     # cp -f ./lib/libOpenCL_lin64.so ../libOpenCL.so
#     ar src ../libOpenCL_lin64.a ./lib/libOpenCL.so 

#     rm -rf ./CLBlast
#     git clone https://github.com/CNugteren/CLBlast.git
#     mkdir -p CLBlast/build; cd CLBlast/build
#     cmake .. -DBUILD_SHARED_LIBS=OFF -DTUNERS=OFF
#     cmake --build . --config Release
#     cmake --install . --prefix ../..

#     cd ../..
#     cp -f ./lib/libclblast.a ../libclblast_lin64.a

#     cd ..

#     cd vendors/llama.cpp
#     mkdir -p build; cd build

#     cmake .. -DLLAMA_CLBLAST=1 -DCMAKE_PREFIX_PATH="../../openclblast" -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
#     cmake --build . --config Release

#     cp libllama.a ../../../libllama_cl_lin64.a

#     cd ../../..

#     g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
#     g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples myllama_llama_api.cpp -o myllama_llama_api.o -c
#     ar src libbinding_cl_lin64.a libllama_cl_lin64.a binding.o myllama_llama_api.o

#     git restore vendors

# elif [ "$1" = "cuda" ]; then
if [ "$1" = "cuda" ]; then
    # cuda
    cd vendors/llama.cpp
    mkdir -p build; cd build

    cmake .. -DLLAMA_CUBLAS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
    cmake --build . --config Release

    cp libllama.a ../../../libllama_cu_lin64.a

    cd ../../..

    g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
    g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples myllama_llama_api.cpp -o myllama_llama_api.o -c
    ar src libbinding_cu_lin64.a libllama_cu_lin64.a binding.o myllama_llama_api.o

    git restore vendors

else
    echo "Invalid argument"

fi
