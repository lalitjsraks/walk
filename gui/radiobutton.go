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

type RadioButton struct {
	Button
}

func NewRadioButton(parent IContainer) (*RadioButton, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		0, syscall.StringToUTF16Ptr("BUTTON"), nil,
		BS_AUTORADIOBUTTON /*|BS_NOTIFY*/ |WS_CHILD|WS_TABSTOP|WS_VISIBLE,
		0, 0, 120, 24, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	rb := &RadioButton{Button: Button{Widget: Widget{hWnd: hWnd, parent: parent}}}
	rb.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = rb

	parent.Children().Add(rb)

	return rb, nil
}

func (*RadioButton) LayoutFlags() LayoutFlags {
	return ShrinkHorz | GrowHorz
}

func (rb *RadioButton) PreferredSize() drawing.Size {
	return rb.dialogBaseUnitsToPixels(drawing.Size{50, 10})
}
