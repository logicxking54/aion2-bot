package main

import "gocv.io/x/gocv"

type RECT struct {
	Left, Top, Right, Bottom int32
}

type MonsterMat struct {
	Path string
	Mat  gocv.Mat
}
