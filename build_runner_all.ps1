<# cpu #>
./clean.ps1
./build_lib.ps1

cd examples/runner
go build -tags cpu -trimpath -ldflags="-w -s" -o ../../bin/run-myllama_cpu.exe
copy-item -force ../../llama.dll ../../bin/
cd ../..

cd bin
tar.exe -a -c -f my-llama_dll_cpu.zip llama.dll
tar.exe -a -c -f my-llama_cpu.zip run-myllama_cpu.exe llama.dll
cd ..


<# clblast #>
./clean.ps1
./build_lib.ps1 clblast

cd examples/runner
go build -tags clblast -trimpath -ldflags="-w -s" -o ../../bin/run-myllama_cl.exe
copy-item -force ../../llama_cl.dll ../../bin/
copy-item -force ../../openclblast/lib/clblast.dll ../../bin/
cd ../..

cd bin
tar.exe -a -c -f my-llama_dll_cl.zip llama_cl.dll clblast.dll
tar.exe -a -c -f my-llama_cl.zip run-myllama_cl.exe llama_cl.dll clblast.dll
cd ..


<# cuda - cublas #>
./clean.ps1
./build_lib.ps1 cuda

cd examples/runner
go build -tags cuda -trimpath -ldflags="-w -s" -o ../../bin/run-myllama_cu.exe
copy-item -force ../../llama_cu.dll ../../bin/
cd ../..

cd bin
tar.exe -a -c -f my-llama_dll_cu.zip llama_cu.dll
tar.exe -a -c -f my-llama_cu.zip run-myllama_cu.exe llama_cu.dll
cd ..


<# cleaning #>
cd bin
remove-item -force -ea 0 *.dll
remove-item -force -ea 0 *.exe
cd ..
