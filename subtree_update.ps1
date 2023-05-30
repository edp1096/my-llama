# Update submodule
git submodule update --recursive --remote

# Update subtree
git subtree pull --prefix=vendors/llama.cpp https://github.com/ggerganov/llama.cpp.git master --squash
