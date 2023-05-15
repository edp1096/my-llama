@echo off

@REM New ggml

@REM cuda
cmd /c "clean.cmd"
cmd /c "build_cu.cmd"
move /y bin\my-llama.exe bin\my-llama_cu.exe
cd bin
tar.exe -a -c -f my-llama_cu.zip my-llama_cu.exe llama.dll
cd ..


@REM clblast
cmd /c "clean.cmd"
cmd /c "build_cl.cmd"
move /y bin\my-llama.exe bin\my-llama_cl.exe
cd bin
tar.exe -a -c -f my-llama_cl.zip my-llama_cl.exe llama.dll clblast.dll
cd ..


@REM cpu
cmd /c "clean.cmd"
cmd /c "mingw32-make.exe"
move /y bin\my-llama.exe bin\my-llama_cpu.exe
cmd /c "mingw32-make.exe clean"


@REM Delete unnecessary files

del bin\my-llama_cu.exe /s /q
del bin\my-llama_cl.exe /s /q
del bin\llama.dll /s /q
del bin\clblast.dll /s /q
