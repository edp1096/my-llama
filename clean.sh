cd vendors/llama.cpp

make clean; rm -rf build/*

cd ../..

rm -f ./*.o
# rm -f ./libbinding_lin64.a
# rm -f ./libllama_lin64.a

if [ "$1" = "all" ]; then
    rm -rf openclblast
fi

git restore vendors
