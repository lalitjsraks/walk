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

type CheckBox struct {
	Button
}

func NewCheckBox(parent IContainer) (*CheckBox, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		BS_AUTOCHECKBOX /*|BS_NOTIFY*/ |WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	cb := &CheckBox{Button: Button{Widget: Widget{hWnd: hWnd, parent: parent}}}
	cb.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = cb

	parent.Children().Add(cb)

	return cb, nil
}

func (*CheckBox) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (cb *CheckBox) PreferredSize() drawing.Size {
	return cb.dialogBaseUnitsToPixels(drawing.Size{50, 10})
}
