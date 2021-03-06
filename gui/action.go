// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"container/vector"
	"os"
)

import (
	"walk/drawing"
)

type actionChangedHandler interface {
	onActionChanged(action *Action) (err os.Error)
}

var (
	// ISSUE: When pressing enter resp. escape,
	// WM_COMMAND with wParam=1 resp. 2 is sent.
	// Maybe there is more to consider.
	nextActionId uint16             = 3
	actionsById  map[uint16]*Action = make(map[uint16]*Action)
)

type Action struct {
	menu              *Menu
	triggeredHandlers vector.Vector
	changedHandlers   vector.Vector
	text              string
	toolTip           string
	image             *drawing.Bitmap
	enabled           bool
	visible           bool
	checkable         bool
	checked           bool
	exclusive         bool
	id                uint16
}

func NewAction() *Action {
	a := &Action{
		enabled: true,
		id:      nextActionId,
		visible: true,
	}

	actionsById[a.id] = a

	nextActionId++

	return a
}

func (a *Action) Checkable() bool {
	return a.checkable
}

func (a *Action) SetCheckable(value bool) (err os.Error) {
	if value != a.checkable {
		old := a.checkable

		a.checkable = value

		err = a.raiseChanged()
		if err != nil {
			a.checkable = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Checked() bool {
	return a.checked
}

func (a *Action) SetChecked(value bool) (err os.Error) {
	if value != a.checked {
		old := a.checked

		a.checked = value

		err = a.raiseChanged()
		if err != nil {
			a.checked = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Enabled() bool {
	return a.enabled
}

func (a *Action) SetEnabled(value bool) (err os.Error) {
	if value != a.enabled {
		old := a.enabled

		a.enabled = value

		err = a.raiseChanged()
		if err != nil {
			a.enabled = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Exclusive() bool {
	return a.exclusive
}

func (a *Action) SetExclusive(value bool) (err os.Error) {
	if value != a.exclusive {
		old := a.exclusive

		a.exclusive = value

		err = a.raiseChanged()
		if err != nil {
			a.exclusive = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Image() *drawing.Bitmap {
	return a.image
}

func (a *Action) SetImage(value *drawing.Bitmap) (err os.Error) {
	if value != a.image {
		old := a.image

		a.image = value

		err = a.raiseChanged()
		if err != nil {
			a.image = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Text() string {
	return a.text
}

func (a *Action) SetText(value string) (err os.Error) {
	if value != a.text {
		old := a.text

		a.text = value

		err = a.raiseChanged()
		if err != nil {
			a.text = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) ToolTip() string {
	return a.toolTip
}

func (a *Action) SetToolTip(value string) (err os.Error) {
	if value != a.toolTip {
		old := a.toolTip

		a.toolTip = value

		err = a.raiseChanged()
		if err != nil {
			a.toolTip = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) Visible() bool {
	return a.visible
}

func (a *Action) SetVisible(value bool) (err os.Error) {
	if value != a.visible {
		old := a.visible

		a.visible = value

		err = a.raiseChanged()
		if err != nil {
			a.visible = old
			a.raiseChanged()
		}
	}

	return
}

func (a *Action) AddTriggeredHandler(handler EventHandler) {
	a.triggeredHandlers.Push(handler)
}

func (a *Action) RemoveTriggeredHandler(handler EventHandler) {
	for i, h := range a.triggeredHandlers {
		if h.(EventHandler) == handler {
			a.triggeredHandlers.Delete(i)
			break
		}
	}
}

func (a *Action) raiseTriggered() {
	for _, handlerIface := range a.triggeredHandlers {
		handler := handlerIface.(EventHandler)
		handler(&eventArgs{a})
	}
}

func (a *Action) addChangedHandler(handler actionChangedHandler) {
	a.changedHandlers.Push(handler)
}

func (a *Action) removeChangedHandler(handler actionChangedHandler) {
	for i, h := range a.changedHandlers {
		if h.(actionChangedHandler) == handler {
			a.changedHandlers.Delete(i)
			break
		}
	}
}

func (a *Action) raiseChanged() (err os.Error) {
	for _, handlerIface := range a.changedHandlers {
		handler := handlerIface.(actionChangedHandler)
		err = handler.onActionChanged(a)
		if err != nil {
			return
		}
	}

	return
}
