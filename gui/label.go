// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
)

import (
	"walk/drawing"
	. "walk/winapi/user32"
)

type Label struct {
	Widget
}

func NewLabel(parent IContainer) (*Label, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("STATIC"), nil,
		WS_CHILD|WS_VISIBLE,
		0, 0, 80, 24, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	l := &Label{Widget: Widget{hWnd: hWnd, parent: parent}}
	l.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = l

	parent.Children().Add(l)

	return l, nil
}

func (*Label) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (l *Label) PreferredSize() drawing.Size {
	return l.dialogBaseUnitsToPixels(drawing.Size{50, 8})
}
