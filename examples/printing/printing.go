// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"runtime"
)

import (
	"walk/drawing"
	"walk/printing"
)

func panicIfErr(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.LockOSThread()

	doc := printing.NewDocument("Walk Printing Example")
	defer doc.Dispose()

	doc.InsertPageBreak()

	text := "Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua."
	font, err := drawing.NewFont("Times New Roman", 12, 0)
	panicIfErr(err)
	color := drawing.RGB(0, 0, 0)
	preferredSize := drawing.Size{1000, 0}
	format := drawing.TextWordbreak

	for i := 0; i < 20; i++ {
		panicIfErr(doc.AddText(fmt.Sprintf("%d) %s", i, text), font, color, preferredSize, format))
	}

	panicIfErr(doc.Print())
}
