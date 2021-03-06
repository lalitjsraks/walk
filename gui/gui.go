// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi/gdi32"
	. "walk/winapi/user32"
)

var defaultFont *drawing.Font

func init() {
	// Initialize default font
	var ncm NONCLIENTMETRICS
	ncm.CbSize = uint(unsafe.Sizeof(ncm))

	if !SystemParametersInfo(SPI_GETNONCLIENTMETRICS, ncm.CbSize, unsafe.Pointer(&ncm), 0) {
		panic("SystemParametersInfo failed")
	}

	hdc := GetDC(0)
	defer ReleaseDC(0, hdc)
	dpi := GetDeviceCaps(hdc, LOGPIXELSY)

	// FIXME: Find out how to get dialog item font and use that.
	var err os.Error
	defaultFont, err = drawing.NewFontFromLOGFONT(&ncm.LfMenuFont, dpi)
	if err != nil {
		panic("failed to create default font")
	}
}
