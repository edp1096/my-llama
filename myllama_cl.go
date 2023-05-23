//go:build clblast
// +build clblast

package myllama

/*
#cgo CXXFLAGS: -Ivendor/llama.cpp -Ivendor/llama.cpp/examples
#cgo LDFLAGS: -static -L. -lstdc++ -lllama_cl -lbinding_cl
#include "binding.h"
*/
import "C"
