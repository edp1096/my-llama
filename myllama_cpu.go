//go:build cpu
// +build cpu

package myllama

/*
#cgo CXXFLAGS: -Ivendors/llama.cpp -Ivendors/llama.cpp/examples
#cgo LDFLAGS: -static -L. -lstdc++ -lllama -lbinding
#include "binding.h"
*/
import "C"
