//go:build !clblast && !cuda
// +build !clblast,!cuda

package main

var (
	deviceType = "cpu"
)
