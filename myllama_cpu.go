//go:build cpu
// +build cpu

package myllama

/*
#cgo CXXFLAGS: -Ivendor/llama.cpp -Ivendor/llama.cpp/examples
#cgo LDFLAGS: -static -L. -lstdc++ -lllama -lbinding
#include "binding.h"
*/
import "C"
