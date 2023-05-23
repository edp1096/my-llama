//go:build clblast
// +build clblast

package myllama

/*
#cgo CXXFLAGS: -Illama.cpp -Illama.cpp/examples
#cgo LDFLAGS: -static -L. -lstdc++ -lllama_cl -lbinding_cl
#include "binding.h"
*/
import "C"
