# Get submodule
cd llama.cpp
git checkout -q master
cd ..


# cpu
./clean.ps1
./build_lib.ps1
./build_cmd.ps1

cd bin
mv run-myllama.exe run-myllama_cpu.exe
tar.exe -a -c -f my-llama_cpu.zip run-myllama_cpu.exe llama.dll
cd ..


# clblast
./clean.ps1
./build_lib.ps1 clblast
./build_cmd.ps1 clblast

cd bin
mv run-myllama.exe run-myllama_cl.exe
tar.exe -a -c -f my-llama_cl.zip run-myllama_cl.exe llama_cl.dll clblast.dll
cd ..


# # cleaning
# cd bin
# rm -f *.dll *.exe
# cd ..
