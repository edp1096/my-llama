//go:build linux && clblast
// +build linux,clblast

package myllama

// #cgo LDFLAGS: -L. -lstdc++ -lllama_cl_lin64 -lbinding_cl_lin64 -lclblast_lin64 -lOpenCL_lin64
import "C"
