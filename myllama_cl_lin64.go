//go:build linux && clblast
// +build linux,clblast

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cl -lbinding_cl
import "C"
