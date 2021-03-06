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
	. "walk/winapi/user32"
)

type Menu struct {
	hMenu   HMENU
	hWnd    HWND
	actions *ActionList
}

func newMenuBar() (*Menu, os.Error) {
	hMenu := CreateMenu()
	if hMenu == 0 {
		return nil, lastError("CreateMenu")
	}

	m := &Menu{hMenu: hMenu}
	m.actions = newActionList(m)

	return m, nil
}

func NewMenu() (*Menu, os.Error) {
	hMenu := CreatePopupMenu()
	if hMenu == 0 {
		return nil, lastError("CreatePopupMenu")
	}

	var mi MENUINFO
	mi.CbSize = uint(unsafe.Sizeof(mi))

	if !GetMenuInfo(hMenu, &mi) {
		return nil, lastError("GetMenuInfo")
	}

	mi.FMask |= MIM_STYLE
	mi.DwStyle = MNS_CHECKORBMP

	if !SetMenuInfo(hMenu, &mi) {
		return nil, lastError("SetMenuInfo")
	}

	m := &Menu{hMenu: hMenu}
	m.actions = newActionList(m)

	return m, nil
}

func (m *Menu) Dispose() {
	if m.hMenu != 0 {
		DestroyMenu(m.hMenu)
		m.hMenu = 0
	}
}

func (m *Menu) IsDisposed() bool {
	return m.hMenu == 0
}

func (m *Menu) Actions() *ActionList {
	return m.actions
}

func (m *Menu) initMenuItemInfoFromAction(mii *MENUITEMINFO, action *Action) {
	mii.CbSize = uint(unsafe.Sizeof(*mii))
	mii.FMask = MIIM_FTYPE | MIIM_ID | MIIM_STRING
	if action.image != nil {
		mii.FMask |= MIIM_BITMAP
		mii.HbmpItem = action.image.Handle()
	}
	mii.FType = MFT_STRING
	mii.WID = uint(action.id)
	mii.DwTypeData = syscall.StringToUTF16Ptr(action.Text())
	mii.Cch = uint(len([]int(action.Text())))

	menu := action.menu
	if menu != nil {
		mii.FMask |= MIIM_SUBMENU
		mii.HSubMenu = menu.hMenu
	}
}

func (m *Menu) onActionChanged(action *Action) (err os.Error) {
	var mii MENUITEMINFO

	m.initMenuItemInfoFromAction(&mii, action)

	if !SetMenuItemInfo(m.hMenu, uint(m.actions.IndexOf(action)), true, &mii) {
		err = newError("SetMenuItemInfo failed")
	}

	return
}

func (m *Menu) onInsertingAction(index int, action *Action) (err os.Error) {
	var mii MENUITEMINFO

	m.initMenuItemInfoFromAction(&mii, action)

	if !InsertMenuItem(m.hMenu, uint(index), true, &mii) {
		return newError("wingui.Menu.onInsertingAction: win32.InsertMenuItem failed")
	}

	action.addChangedHandler(m)

	menu := action.menu
	if menu != nil {
		menu.hWnd = m.hWnd
	}

	if m.hWnd != 0 {
		DrawMenuBar(m.hWnd)
	}

	return
}

func (m *Menu) onRemovingAction(index int, action *Action) (err os.Error) {
	panic("not implemented")
}

func (m *Menu) onClearingActions() (err os.Error) {
	panic("not implemented")
}
