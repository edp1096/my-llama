#!/bin/sh

./clean.sh
./build_lib.sh

./clean.sh
./build_lib.sh cuda

./clean.sh
./build_lib.sh clblast
