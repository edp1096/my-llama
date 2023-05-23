<# llama.cpp #>
cd llama.cpp

mkdir -f build
cd build

# cmake .. -DLLAMA_STATIC=1 -DBUILD_SHARED_LIBS=0 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake .. -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=1 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

cp bin/Release/llama.dll ../../llama.dll

cd ..

# gcc -I. -Iexamples ggml.c -o ggml.o -c
# g++ -I. -Iexamples llama.cpp -o llama.o -c
# g++ -I. -Iexamples examples/common.cpp -o common.o -c

# cp ggml.o ../ggml.o
# cp llama.o ../llama.o
# cp common.o ../common.o

cd ..
# ar rcs libllama.a llama.o ggml.o common.o

gendef ./llama.dll
dlltool -k -d ./llama.def -l ./libllama.a

<# binding #>
cd binding
g++ -I../llama.cpp -I../llama.cpp/examples -lllama binding.cpp -o binding.o -c -L..
# ar rcs libbinding.a binding.o libllama.a
ar rcs libbinding.a binding.o

cp libbinding.a ../libbinding.a

cd ..
