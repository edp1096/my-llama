cd vendors/llama.cpp

make clean; rm -rf build/*

cd ../..

rm -f ./*.o


if [ "$1" = "all" ]; then
    rm -rf openclblast
fi

git restore vendors
