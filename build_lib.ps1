<# llama.cpp - msvc/cmake #>
cd llama.cpp

mkdir -f build
cd build

cmake .. -DLLAMA_CLBLAST=1 -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

cp bin/Release/llama.dll ../../llama.dll

cd ../..

gendef ./llama.dll
dlltool -k -d ./llama.def -l ./libllama.a


# Not use
# <# llama.cpp - mingw #>
# mingw32-make.exe CC=gcc -C llama.cpp ggml.o
# mingw32-make.exe CC=gcc -C llama.cpp llama.o
# mingw32-make.exe CC=gcc -C llama.cpp common.o

# ar src libllama.a llama.cpp/llama.o llama.cpp/ggml.o llama.cpp/common.o
# Not use


<# binding #>
g++ -O3 -DNDEBUG -std=c++11 -fPIC -march=native -mtune=native -I./llama.cpp -I./llama.cpp/examples binding.cpp -o binding.o -c
ar src libbinding.a binding.o
