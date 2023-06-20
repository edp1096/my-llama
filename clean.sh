rm -rf build
cd llama.cpp; make clean; cd ..

rm -f ./*.o

if [ "$1" = "all" ]; then
    rm -rf vendors/openclblast
fi

git checkout -- llama.cpp
git clean llama.cpp -df
