# Get submodule
cd llama.cpp
git checkout -q master
cd ..


# clblast
clean.ps1
build_cl.ps1
mv -f bin/my-llama.exe bin/my-llama_cl.exe
cd bin
tar.exe -a -c -f my-llama_cl.zip my-llama_cl.exe llama.dll clblast.dll
cd ..


# cpu
clean.ps1
mingw32-make.exe
mv -f bin/my-llama.exe bin/my-llama_cpu.exe
mingw32-make.exe clean


# Delete unnecessary files
Remove-Item -Path "bin\my-llama_cl.exe" -Recurse -Force
Remove-Item -Path "bin\llama.dll" -Recurse -Force
Remove-Item -Path "bin\clblast.dll" -Recurse -Force
