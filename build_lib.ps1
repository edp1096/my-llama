$useCLBlast = 0

<# clblast #>
if ($args[0] -eq "clblast") {
    cp -f llama.cpp_deallocate/* llama.cpp/

    if (-not (Test-Path -Path "clblast.zip")) {
        echo "Downloading CLBlast..."
        curl --progress-bar -Lo clblast.zip "https://github.com/CNugteren/CLBlast/releases/download/1.5.3/CLBlast-1.5.3-Windows-x64.zip"
        # curl --progress-bar -Lo clblast.7z "https://github.com/CNugteren/CLBlast/releases/download/1.6.0/CLBlast-1.6.0-windows-x64.7z"
    }

    if (-not (Test-Path -Path "openclblast_cmake.zip")) {
        echo "Downloading OpenCLBlast CMake..."
        curl --progress-bar -Lo opencl.zip "https://github.com/KhronosGroup/OpenCL-SDK/releases/download/v2023.04.17/OpenCL-SDK-v2023.04.17-Win-x64.zip"
    }

    mkdir -f openclblast
    rm -rf openclblast/*
    tar -xf clblast.zip -C openclblast
    # if (-not (Test-Path -Path "7zr.exe")) {
    #     echo "Downloading 7zr..."
    #     curl --progress-bar -Lo 7zr.exe "https://www.7-zip.org/a/7zr.exe"
    # }
    # .\7zr.exe x -y .\clblast.7z -o"." -r
    # mv -f clblast_release_dir/* openclblast/
    # rm -rf clblast_release_dir
    tar -xf opencl.zip -C openclblast

    cp -rf openclblast/OpenCL-SDK-v2023.04.17-Win-x64/* openclblast
    cp -rf openclblast_cmake/*.cmake openclblast/lib/cmake/CLBlast

    rm -rf openclblast/OpenCL-SDK-v2023.04.17-Win-x64

    $useCLBlast = 1
}


<# llama.cpp - msvc/cmake #>
cd llama.cpp

mkdir -f build
cd build

cmake .. -DLLAMA_CLBLAST=$useCLBlast -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

cp bin/Release/llama.dll ../../llama.dll

cd ../..

gendef ./llama.dll
dlltool -k -d ./llama.def -l ./libllama.a


<# binding #>
g++ -O3 -DNDEBUG -std=c++11 -fPIC -march=native -mtune=native -I./llama.cpp -I./llama.cpp/examples binding.cpp -o binding.o -c
ar src libbinding.a binding.o


<# clblast - restore modified from llama.cpp_deallocate #>
if ($args[0] -eq "clblast") {
    cd llama.cpp
    git clean -f .
    git reset --hard
    cd ..
}
