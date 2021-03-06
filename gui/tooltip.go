// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	"walk/drawing"
	. "walk/winapi"
	. "walk/winapi/comctl32"
	. "walk/winapi/user32"
)

type ToolTip struct {
	Widget
}

func NewToolTip(parent IContainer) (*ToolTip, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	hWnd := CreateWindowEx(
		WS_EX_TOPMOST, syscall.StringToUTF16Ptr("tooltips_class32"), nil,
		TTS_ALWAYSTIP|TTS_BALLOON|WS_POPUP,
		CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT,
		parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	tt := &ToolTip{Widget: Widget{hWnd: hWnd, parent: parent}}
	tt.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = tt

	parent.Children().Add(tt)

	SetWindowPos(hWnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE)

	return tt, nil
}

func (*ToolTip) LayoutFlags() LayoutFlags {
	return 0
}

func (tt *ToolTip) PreferredSize() drawing.Size {
	return drawing.Size{0, 0}
}

func (tt *ToolTip) Title() string {
	var gt TTGETTITLE

	buf := make([]uint16, 128)

	gt.DwSize = uint(unsafe.Sizeof(gt))
	gt.Cch = uint(len(buf))
	gt.PszTitle = &buf[0]

	SendMessage(tt.hWnd, TTM_GETTITLE, 0, uintptr(unsafe.Pointer(&gt)))

	return syscall.UTF16ToString(buf)
}

func (tt *ToolTip) SetTitle(value string) os.Error {
	if FALSE == SendMessage(tt.hWnd, TTM_SETTITLE, uintptr(TTI_INFO), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value)))) {
		return newError("TTM_SETTITLE failed")
	}

	return nil
}

func (tt *ToolTip) AddWidget(widget IWidget, text string) os.Error {
	var ti TOOLINFO

	ti.CbSize = uint(unsafe.Sizeof(ti))
	parent := widget.Parent()
	if parent != nil {
		ti.Hwnd = parent.Handle()
	}
	ti.UFlags = TTF_IDISHWND | TTF_SUBCLASS
	ti.UId = uintptr(widget.Handle())
	ti.LpszText = syscall.StringToUTF16Ptr(text)

	if FALSE == SendMessage(tt.hWnd, TTM_ADDTOOL, 0, uintptr(unsafe.Pointer(&ti))) {
		return newError("TTM_ADDTOOL failed")
	}

	return nil
}

func (tt *ToolTip) RemoveWidget(widget IWidget) os.Error {
	panic("not implemented")
}
