<# llama.cpp #>
# cd llama.cpp

# mkdir -f build
# cd build

# # cmake .. -DLLAMA_STATIC=1 -DBUILD_SHARED_LIBS=0 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
# cmake .. -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
# cmake --build . --config Release

# cp bin/Release/llama.dll ../../llama.dll

# cd ../..

mingw32-make.exe CC=gcc -C llama.cpp ggml.o
mingw32-make.exe CC=gcc -C llama.cpp llama.o
mingw32-make.exe CC=gcc -C llama.cpp common.o

ar src libllama.a llama.cpp/llama.o llama.cpp/ggml.o llama.cpp/common.o

# gcc -I. -Iexamples ggml.c -o ggml.o -c
# g++ -I. -Iexamples llama.cpp -o llama.o -c
# g++ -I. -Iexamples examples/common.cpp -o common.o -c

# cp ggml.o ../ggml.o
# cp llama.o ../llama.o
# cp common.o ../common.o

# cd ..
# ar rcs libllama.a llama.o ggml.o common.o

# gendef ./llama.dll
# dlltool -k -d ./llama.def -l ./libllama.a

<# binding #>
g++ -O3 -DNDEBUG -std=c++11 -fPIC -march=native -mtune=native -I./llama.cpp -I./llama.cpp/examples binding/binding.cpp -o binding.o -c
# ar rcs libbinding.a binding.o libllama.a
ar src libbinding.a binding.o

# cp libbinding.a ../libbinding.a
# cd ..

<# compile golang sample cmd #>
# go build -trimpath -ldflags="-w -s"
