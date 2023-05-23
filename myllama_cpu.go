//go:build cpu
// +build cpu

package myllama

/*
#cgo CXXFLAGS: -Illama.cpp -Illama.cpp/examples
#cgo LDFLAGS: -static -L. -lstdc++ -lllama -lbinding
#include "binding.h"
*/
import "C"
