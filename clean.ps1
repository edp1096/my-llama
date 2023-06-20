rm -f *.o
rm -f *.s
rm -f *.def
rm -f *.dll

rm -f vendors/*.o
rm -f vendors/*.a
rm -f vendors/*.def
rm -f vendors/*.dll

rm -f llama.cpp/*.o
rm -f llama.cpp/*.a

rm -rf build
rm -f llama.cpp/build-info.h

rm -f 7zr.exe

if ($args[0] -eq "all") {
    rm -f opencl.zip
    rm -f clblast.7z
    rm -f clblast.zip
    rm -rf openclblast

    rm -f output.log
}

./subtree_restore.ps1
