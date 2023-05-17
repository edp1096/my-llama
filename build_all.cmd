@echo off

@REM @REM Old ggml

@REM @REM cuda
@REM cd llama.cpp
@REM git checkout -q b608b55a
@REM cd ..

@REM cmd /c "clean.cmd"
@REM cmd /c "build_cu.cmd USE_OLD_GGML"
@REM move /y bin\my-llama.exe bin\my-llama_cu.exe
@REM cd bin
@REM tar.exe -a -c -f my-llama_cu_old_ggml.zip my-llama_cu.exe llama.dll
@REM cd ..


@REM @REM clblast
@REM cmd /c "clean.cmd"
@REM cmd /c "build_cl.cmd USE_OLD_GGML"
@REM move /y bin\my-llama.exe bin\my-llama_cl.exe
@REM cd bin
@REM tar.exe -a -c -f my-llama_cl_old_ggml.zip my-llama_cl.exe llama.dll clblast.dll
@REM cd ..


@REM @REM cpu
@REM cmd /c "clean.cmd"
@REM cmd /c "mingw32-make.exe USE_OLD_GGML=1"
@REM move /y bin\my-llama.exe bin\my-llama_cpu_old_ggml.exe
@REM cmd /c "mingw32-make.exe clean"



@REM New ggml

@REM cuda
cd llama.cpp
git checkout -q master
cd ..

cmd /c "clean.cmd"
cmd /c "build_cu.cmd"
move /y bin\my-llama.exe bin\my-llama_cu.exe
cd bin
tar.exe -a -c -f my-llama_cu.zip my-llama_cu.exe llama.dll my-runner.exe
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
