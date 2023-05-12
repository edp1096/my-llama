@echo off

cd llama.cpp

md build 2>nul
cd build

cmake .. -DLLAMA_CUBLAS=1 -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=1 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

copy bin\Release\llama.dll ..\..

@REM Static not work
@REM copy Release\llama.lib ..\..
@REM copy examples\common.dir\Release\common.lib ..\..
@REM copy ggml.dir\Release\ggml.lib ..\..

cd ..\..


@REM Static not work
@REM C:\Windows\SysWOW64\WindowsPowerShell\v1.0\powershell.exe -noe -c "&{Import-Module 'C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\Tools\Microsoft.VisualStudio.DevShell.dll'; Enter-VsDevShell 4da9a52a} ; lib.exe /OUT:common.obj common.lib; exit"
@REM C:\Windows\SysWOW64\WindowsPowerShell\v1.0\powershell.exe -noe -c "&{Import-Module 'C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\Tools\Microsoft.VisualStudio.DevShell.dll'; Enter-VsDevShell 4da9a52a} ; lib.exe /OUT:ggml.obj ggml.lib; exit"
@REM C:\Windows\SysWOW64\WindowsPowerShell\v1.0\powershell.exe -noe -c "&{Import-Module 'C:\Program Files\Microsoft Visual Studio\2022\Community\Common7\Tools\Microsoft.VisualStudio.DevShell.dll'; Enter-VsDevShell 4da9a52a} ; lib.exe /OUT:llama.obj llama.lib; exit"

gendef.exe llama.dll
dlltool.exe -k -d llama.def -l libllama.a

if "%1" == "USE_OLD_GGML" (
    mingw32-make.exe USE_OLD_GGML=1 build_for_cuda 
) else (
    mingw32-make.exe build_for_cuda
)

copy llama.dll .\bin /y