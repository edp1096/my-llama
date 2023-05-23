# Get submodule
cd llama.cpp
git checkout -q master
cd ..

cd llama.cpp
git clean -f .
git reset --hard
cd ..


# cpu
./clean.ps1
./build_lib.ps1
./build_cmd.ps1

cd bin
tar.exe -a -c -f my-llama_dll_cpu.zip llama.dll
tar.exe -a -c -f my-llama_cpu.zip run-myllama_cpu.exe llama.dll
cd ..


# clblast
./clean.ps1
./build_lib.ps1 clblast
./build_cmd.ps1 clblast

cd bin
tar.exe -a -c -f my-llama_dll_cl.zip llama_cl.dll clblast.dll
tar.exe -a -c -f my-llama_cl.zip run-myllama_cl.exe llama_cl.dll clblast.dll
cd ..


# cleaning
cd bin
rm -f *.dll *.exe
cd ..
