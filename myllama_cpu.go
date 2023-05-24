//go:build !clblast
// +build !clblast

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama -lbinding
import "C"
