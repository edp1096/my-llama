//go:build cpu || !clblast
// +build cpu !clblast

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama -lbinding
import "C"
