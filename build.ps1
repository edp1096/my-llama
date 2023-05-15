cd llama.cpp

md build -ea 0
cd build

cmake .. -DLLAMA_CUBLAS=1 -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

copy bin\Release\llama.dll ..\..

cd ..\..

gendef.exe llama.dll
dlltool.exe -k -d llama.def -l libllama.a

go build -ldflags "-w -s" -o bin/
