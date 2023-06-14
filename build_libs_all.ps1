git restore vendors

<# cpu #>
./clean.ps1
./build_lib.ps1


<# clblast #>
./clean.ps1
./build_lib.ps1 clblast


<# cuda - cublas #>
./clean.ps1
./build_lib.ps1 cuda


<# cleaning #>
cd bin
rm -f *.dll *.exe
cd ..
