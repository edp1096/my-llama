cd vendors/llama.cpp
mkdir -p build
cd build

# cmake .. -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake .. -DLLAMA_CUBLAS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

# cp libllama.a ../../../libllama_lin64.a
cp libllama.a ../../../libllama_cu_lin64.a

cd ../../..

g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
g++ -static -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding_llama_api.cpp -o binding_llama_api.o -c
# ar src libbinding_lin64.a libllama_lin64.a binding.o binding_llama_api.o
ar src libbinding_cu_lin64.a libllama_cu_lin64.a binding.o binding_llama_api.o
