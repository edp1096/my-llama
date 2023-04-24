@echo off

cd llama.cpp

md build 2>nul
cd build

cmake .. -DLLAMA_CUBLAS=1 -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

copy bin\Release\llama.dll ..\..\
copy Release\llama.lib ..\..\

cd ..\..

C:\Windows\SysWOW64\WindowsPowerShell\v1.0\powershell.exe -noe -c "&{Import-Module 'C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\Tools\Microsoft.VisualStudio.DevShell.dll'; Enter-VsDevShell 4da9a52a} ; lib.exe /OUT:llama.obj llama.lib; exit"

mingw32-make.exe build_for_cuda

copy llama.dll .\bin /y