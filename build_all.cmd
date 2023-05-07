@echo off

cmd /c "clean.cmd"
cmd /c "build.cmd"

move /y bin\my-llama.exe bin\my-llama_cu.exe

cmd /c "clean.cmd"

cmd /c "mingw32-make.exe"

move /y bin\my-llama.exe bin\my-llama_cpu.exe

cmd /c "mingw32-make.exe clean"
