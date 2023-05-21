# Waiting for 0cc4m's PR accepting. until then, must copy this to llama.cpp folder.
cp -f llama.cpp_0cc4m/* llama.cpp/


if (-not (Test-Path -Path "clblast.zip")) {
    echo "Downloading CLBlast..."
    curl --progress-bar -Lo clblast.zip "https://github.com/CNugteren/CLBlast/releases/download/1.5.3/CLBlast-1.5.3-Windows-x64.zip"
}

if (-not (Test-Path -Path "openclblast_cmake.zip")) {
    echo "Downloading OpenCLBlast CMake..."
    curl --progress-bar -Lo opencl.zip "https://github.com/KhronosGroup/OpenCL-SDK/releases/download/v2023.04.17/OpenCL-SDK-v2023.04.17-Win-x64.zip"
}

mkdir -Force -Path "openclblast"
tar -xf clblast.zip -C openclblast
tar -xf opencl.zip -C openclblast

cp -rf openclblast/OpenCL-SDK-v2023.04.17-Win-x64/* openclblast
cp -rf openclblast_cmake/*.cmake openclblast/lib/cmake/CLBlast

rm -rf openclblast/OpenCL-SDK-v2023.04.17-Win-x64


cd llama.cpp

mkdir -f build
cd build

cmake .. -DCMAKE_PREFIX_PATH='../openclblast' -DLLAMA_CLBLAST=1 -DBUILD_SHARED_LIBS=1 -DLLAMA_BUILD_EXAMPLES=0 -DLLAMA_BUILD_TESTS=0
cmake --build . --config Release

cp bin/Release/llama.dll ../..

cd ../..

gendef.exe llama.dll
dlltool.exe -k -d llama.def -l libllama.a

mingw32-make.exe build_for_cuda


cp -f llama.dll bin/
cp -f openclblast/lib/clblast.dll bin/


# Waiting for 0cc4m's PR accepting. until then, must copy this to llama.cpp folder.
cd llama.cpp
git clean -f .
git reset --hard
cd ..