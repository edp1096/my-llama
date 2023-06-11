$llamaCppSharedLibName="llama"
$llamaCppSharedLibExt="dll"

$dllName="llama.dll"
$defName="llama.def"
$libLlamaName="libllama.a"
$libBindingName="libbinding.a"

$cmakePrefixPath=""
$cmakeUseCLBLAST="OFF"
$cmakeUseCUDA="OFF"


<# All #>
# llama.cpp compile error bdbda1b1 so always copy
# cp -f llama.cpp_mod/* vendors/llama.cpp/

<# Prepare clblast and opencl #>
if ($args[0] -eq "clblast") {
    # llama.cpp compile error bdbda1b1 so always copy
    # cp -f llama.cpp_mod/* vendors/llama.cpp/

    if (-not (Test-Path -Path "opencl.zip")) {
        echo "Downloading OpenCL..."
        curl --progress-bar -Lo opencl.zip "https://github.com/KhronosGroup/OpenCL-SDK/releases/download/v2023.04.17/OpenCL-SDK-v2023.04.17-Win-x64.zip"
    }

    if (-not (Test-Path -Path "clblast.zip")) {
        echo "Downloading CLBlast..."
        # curl --progress-bar -Lo clblast.zip "https://github.com/CNugteren/CLBlast/releases/download/1.5.3/CLBlast-1.5.3-Windows-x64.zip"
        # curl --progress-bar -Lo clblast.zip "https://ci.appveyor.com/api/buildjobs/nikwayllaa7nia4c/artifacts/CLBlast-1.6.0-Windows-x64.zip"
        curl --progress-bar -Lo clblast.7z "https://github.com/CNugteren/CLBlast/releases/download/1.6.0/CLBlast-1.6.0-windows-x64.7z"
    }

    mkdir -f openclblast >$null
    rm -rf openclblast/*
    tar -xf opencl.zip -C openclblast
    # tar -xf clblast.zip -C openclblast
    if (-not (Test-Path -Path "7zr.exe")) {
        echo "Downloading 7zr..."
        curl --progress-bar -Lo 7zr.exe "https://www.7-zip.org/a/7zr.exe"
    }
    .\7zr.exe x -y .\clblast.7z -o"." -r
    mv -f CLBlast-1.6.0-windows-x64/* openclblast/
    rm -f ./7zr.exe
    rm -rf CLBlast-1.6.0-windows-x64

    cp -rf openclblast/OpenCL-SDK-v2023.04.17-Win-x64/* openclblast
    cp -rf openclblast_cmake/*.cmake openclblast/lib/cmake/CLBlast

    rm -rf openclblast/OpenCL-SDK-v2023.04.17-Win-x64

    $dllName="llama_cl.dll"
    $defName="llama_cl.def"
    $libLlamaName="libllama_cl.a"
    $libBindingName="libbinding_cl.a"

    $cmakePrefixPath="../../openclblast"
    $cmakeUseCLBLAST="ON"
}

if ($args[0] -eq "cuda") {
    # llama.cpp compile error bdbda1b1 so always copy
    # cp -f llama.cpp_mod/* vendors/llama.cpp/

    $dllName="llama_cu.dll"
    $defName="llama_cu.def"
    $libLlamaName="libllama_cu.a"
    $libBindingName="libbinding_cu.a"

    $cmakeUseCUDA="ON"
}


<# Compile vendors/llama.cpp - msvc/cmake #>
cd vendors/llama.cpp

mkdir -f build >$null
cd build

cmake .. -DCMAKE_PREFIX_PATH="$cmakePrefixPath" -DLLAMA_CUBLAS="$cmakeUseCUDA" -DLLAMA_CLBLAST="$cmakeUseCLBLAST" -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
# cmake --build . --config Release
cmake --build . --config Debug

# cp bin/Release/$llamaCppSharedLibName.$llamaCppSharedLibExt ../../../$dllName
cp bin/Debug/$llamaCppSharedLibName.$llamaCppSharedLibExt ../../../$dllName

cd ../../..

gendef ./$dllName
if ($args[0] -eq "clblast") {
    (Get-Content -Path "$defName") -replace "$llamaCppSharedLibName.$llamaCppSharedLibExt", "llama_cl.dll" | Set-Content -Path "$defName"
}
if ($args[0] -eq "cuda") {
    (Get-Content -Path "$defName") -replace "$llamaCppSharedLibName.$llamaCppSharedLibExt", "llama_cu.dll" | Set-Content -Path "$defName"
}
dlltool -k -d ./$defName -l ./$libLlamaName

<# Compile binding #>
g++ -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples binding.cpp -o binding.o -c
g++ -O3 -std=c++11 -fPIC -march=native -mtune=native -I./vendors/llama.cpp -I./vendors/llama.cpp/examples myllama_llama_api.cpp -o myllama_llama_api.o -c
ar src $libBindingName myllama_llama_api.o binding.o


# # <# Restore overwritten vendors/llama.cpp_mod for clblast/cuda to original commit #>
# if ($args[0] -eq "clblast" -or $args[0] -eq "cuda") {
#     git restore vendors
# }
