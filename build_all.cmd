@echo off

cmd /c "clean.cmd"
cmd /c "build_cu.cmd"
move /y bin\my-llama.exe bin\my-llama_cu.exe

cmd /c "clean.cmd"

cmd /c "mingw32-make.exe"
move /y bin\my-llama.exe bin\my-llama_cpu.exe

cmd /c "mingw32-make.exe clean"

tar.exe -a -c -f bin\my-llama_cu.zip bin\my-llama_cu.exe bin\llama.dll
