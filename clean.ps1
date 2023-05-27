rm -f *.o
rm -f *.s
rm -f *.def
rm -f *.dll

rm -f vendors/*.o
rm -f vendors/*.a
rm -f vendors/*.def
rm -f vendors/*.dll

rm -f vendors/llama.cpp/*.o
rm -f vendors/llama.cpp/*.a

rm -rf vendors/llama.cpp/build
rm -f vendors/llama.cpp/build-info.h

rm -f output.log >$null
rm -f 7zr.exe

if ($args[0] -eq "all") {
    rm -f opencl.zip
    rm -f clblast.7z
    rm -f clblast.zip
    rm -rf openclblast
}

./subtree_restore.ps1