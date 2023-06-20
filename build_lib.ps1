$threadCount = 8

$llamaCppSharedLibName="llama"
$llamaCppSharedLibExt="dll"

$dllName="llama.dll"
$defName="llama.def"
$libLlamaName="libllama.a"
$libMyllamaName="libmyllama.a"

$cmakePrefixPath=""
$cmakeUseCLBLAST="OFF"
$cmakeUseCUDA="OFF"


cd vendors
import-module bitstransfer

<# Prepare clblast and opencl #>
if ($args[0] -eq "clblast") {
    if (-not (Test-Path -Path "opencl.zip")) {
        echo "Downloading OpenCL..."
        start-bitstransfer -destination opencl.zip -source "https://github.com/KhronosGroup/OpenCL-SDK/releases/download/v2023.04.17/OpenCL-SDK-v2023.04.17-Win-x64.zip"
    }

    if (-not (Test-Path -Path "clblast.zip")) {
        echo "Downloading CLBlast..."
        # "https://ci.appveyor.com/api/buildjobs/nikwayllaa7nia4c/artifacts/CLBlast-1.6.0-Windows-x64.zip"
        start-bitstransfer -destination clblast.7z -source "https://github.com/CNugteren/CLBlast/releases/download/1.6.0/CLBlast-1.6.0-windows-x64.7z"
    }

    mkdir -f openclblast >$null
    remove-item -r -force -ea 0 openclblast/*
    tar -xf opencl.zip -C openclblast
    # tar -xf clblast.zip -C openclblast
    if (-not (Test-Path -Path "7zr.exe")) {
        echo "Downloading 7zr..."
        start-bitstransfer -destination 7zr.exe -source "https://www.7-zip.org/a/7zr.exe"
    }
    .\7zr.exe x -y .\clblast.7z -o"." -r
    # mv -f CLBlast-1.6.0-windows-x64/* openclblast/
    move-item -force CLBlast-1.6.0-windows-x64/* openclblast/
    remove-item -force -ea 0 ./7zr.exe
    remove-item -r -force -ea 0 CLBlast-1.6.0-windows-x64

    copy-item -r -force openclblast/OpenCL-SDK-v2023.04.17-Win-x64/* openclblast
    copy-item -r -force openclblast_cmake/*.cmake openclblast/lib/cmake/CLBlast

    remove-item -r -force -ea 0 openclblast/OpenCL-SDK-v2023.04.17-Win-x64

    $dllName="llama_cl.dll"
    $defName="llama_cl.def"
    $libLlamaName="libllama_cl.a"
    $libMyllamaName="libmyllama_cl.a"

    $cmakePrefixPath="../vendors/openclblast"
    $cmakeUseCLBLAST="ON"
}

if ($args[0] -eq "cuda") {
    $dllName="llama_cu.dll"
    $defName="llama_cu.def"
    $libLlamaName="libllama_cu.a"
    $libMyllamaName="libmyllama_cu.a"

    $cmakeUseCUDA="ON"
}

cd ..


<# Compile llama.cpp - msvc/cmake #>
cd llama.cpp

# mkdir -f build >$null
# cd build

cmake -B ../build . `
    -DCMAKE_CXX_FLAGS="/EHsc /wd4819" -DCMAKE_CUDA_FLAGS="-Xcompiler /wd4819" `
    -DCMAKE_PREFIX_PATH="$cmakePrefixPath" `
    -DBUILD_SHARED_LIBS="ON" `
    -DLLAMA_CUBLAS="$cmakeUseCUDA" -DLLAMA_CLBLAST="$cmakeUseCLBLAST" `
    -DLLAMA_BUILD_EXAMPLES="OFF" -DLLAMA_BUILD_TESTS="OFF"

cd ../build

cmake --build . --config Release -j $threadCount

copy-item -force bin/Release/$llamaCppSharedLibName.$llamaCppSharedLibExt ../../../$dllName

cd ../../..

gendef ./$dllName
if ($args[0] -eq "clblast") {
    (Get-Content -Path "$defName") -replace "$llamaCppSharedLibName.$llamaCppSharedLibExt", "llama_cl.dll" | Set-Content -Path "$defName"
}
if ($args[0] -eq "cuda") {
    (Get-Content -Path "$defName") -replace "$llamaCppSharedLibName.$llamaCppSharedLibExt", "llama_cu.dll" | Set-Content -Path "$defName"
}
dlltool -k -d ./$defName -l ./$libLlamaName

<# Compile binding - myllama, myllama_llama_api #>
g++ `
    -O3 -std=c++11 -fPIC -march=native -mtune=native `
    -I./llama.cpp -I./llama.cpp/examples `
    myllama.cpp -o myllama.o -c
g++ `
    -O3 -std=c++11 -fPIC -march=native -mtune=native `
    -I./llama.cpp -I./llama.cpp/examples `
    myllama_llama_api.cpp -o myllama_llama_api.o -c

ar src $libMyllamaName myllama_llama_api.o myllama.o


# <# Restore overwritten llama.cpp_mod for clblast/cuda to original commit #>
# No more necessary
# if ($args[0] -eq "clblast" -or $args[0] -eq "cuda") {
#     git restore vendors
# }
