//go:build linux && clblast
// +build linux,clblast

package myllama

// #cgo LDFLAGS: -L. -lstdc++ -lllama_cl_lin64 -lbinding_cl_lin64 -lclblast -lOpenCL
import "C"
