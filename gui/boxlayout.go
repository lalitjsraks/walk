// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	//    "log"
	"os"
)

import (
	"walk/drawing"
)

type BoxLayout struct {
	container IContainer
	margins   *Margins
	spacing   int
	vertical  bool
}

func NewHBoxLayout() *BoxLayout {
	return &BoxLayout{margins: &Margins{}}
}

func NewVBoxLayout() *BoxLayout {
	return &BoxLayout{margins: &Margins{}, vertical: true}
}

func (l *BoxLayout) Container() IContainer {
	return l.container
}

func (l *BoxLayout) SetContainer(value IContainer) {
	if value != l.container {
		if l.container != nil {
			l.container.SetLayout(nil)
		}

		l.container = value

		if value != nil && value.Layout() != Layout(l) {
			value.SetLayout(l)

			l.Update(true)
		}
	}
}

func (l *BoxLayout) Margins() *Margins {
	return l.margins
}

func (l *BoxLayout) SetMargins(value *Margins) os.Error {
	if value == nil {
		return newError("margins cannot be nil")
	}

	l.margins = value

	return nil
}

func (l *BoxLayout) Spacing() int {
	return l.spacing
}

func (l *BoxLayout) SetSpacing(value int) os.Error {
	if value != l.spacing {
		if value < 0 {
			return newError("spacing cannot be negative")
		}

		l.spacing = value

		l.Update(false)
	}

	return nil
}

func (l *BoxLayout) Update(reset bool) (err os.Error) {
	if l.container == nil {
		return
	}

	//    log.Stdout("*BoxLayout.Update")

	widgets := make([]IWidget, 0, l.container.Children().Len())

	children := l.container.Children()
	j := 0
	for i := 0; i < cap(widgets); i++ {
		widget := children.At(i)

		ps := widget.PreferredSize()
		if ps.Width == 0 && ps.Height == 0 && widget.LayoutFlags() == 0 {
			continue
		}

		widgets = widgets[0 : j+1]
		widgets[j] = widget
		j++
	}

	widgetCount := len(widgets)

	if widgetCount == 0 {
		return
	}

	// We will start by collecting some valuable information.
	flags := make([]LayoutFlags, widgetCount)
	prefSizes := make([]drawing.Size, widgetCount)
	var prefSizeSum drawing.Size
	var shrinkHorzCount, growHorzCount, shrinkVertCount, growVertCount int

	for i := 0; i < widgetCount; i++ {
		widget := widgets[i]

		ps := widget.PreferredSize()

		maxSize, err := widget.MaxSize()
		if err != nil {
			return err
		}

		lf := widget.LayoutFlags()
		if maxSize.Width > 0 {
			lf &^= GrowHorz
			ps.Width = maxSize.Width
		}
		if maxSize.Height > 0 {
			lf &^= GrowVert
			ps.Height = maxSize.Height
		}

		if lf&ShrinkHorz > 0 {
			shrinkHorzCount++
		}
		if lf&GrowHorz > 0 {
			growHorzCount++
		}
		if lf&ShrinkVert > 0 {
			shrinkVertCount++
		}
		if lf&GrowVert > 0 {
			growVertCount++
		}
		flags[i] = lf

		prefSizeSum.Width += ps.Width
		prefSizeSum.Height += ps.Height
		prefSizes[i] = ps
	}

	cb, err := l.container.ClientBounds()
	if err != nil {
		return
	}

	spacingSum := (widgetCount - 1) * l.spacing

	// Now do the actual layout thing.
	if l.vertical {
		diff := cb.Height - l.margins.Top - prefSizeSum.Height - spacingSum - l.margins.Bottom

		reqW := 0

		for i, s := range prefSizes {
			if s.Width > reqW && (flags[i]&ShrinkHorz == 0) {
				reqW = s.Width
			}
		}
		//        if reqW == 0 {
		reqW = cb.Width - l.margins.Left - l.margins.Right
		//        }

		var change int
		if diff < 0 {
			if shrinkVertCount > 0 {
				change = diff / shrinkVertCount
			}
		} else {
			if growVertCount > 0 {
				change = diff / growVertCount
			}
		}

		//        log.Stdoutf("*BoxLayout.Update: widgetCount: %d, cb: %+v, prefSizeSum: %+v, diff: %d, change: %d, reqW: %d", widgetCount, cb, prefSizeSum, diff, change, reqW)

		y := cb.Y + l.margins.Top
		for i := 0; i < widgetCount; i++ {
			widget := widgets[i]

			h := prefSizes[i].Height

			switch {
			case change < 0:
				if flags[i]&ShrinkVert > 0 {
					h += change
				}

			case change > 0:
				if flags[i]&GrowVert > 0 {
					h += change
				}
			}

			bounds := drawing.Rectangle{cb.X + l.margins.Left, y, reqW, h}

			//            log.Stdoutf("*BoxLayout.Update: bounds: %+v", bounds)

			widget.SetBounds(bounds)

			y += h + l.spacing
		}
	} else {
		diff := cb.Width - l.margins.Left - prefSizeSum.Width - spacingSum - l.margins.Right
		reqH := 0

		for i, s := range prefSizes {
			if s.Height > reqH && (flags[i]&ShrinkVert == 0) {
				reqH = s.Height
			}
		}
		//        if reqH == 0 {
		reqH = cb.Height - l.margins.Top - l.margins.Bottom
		//        }

		var change int
		if diff < 0 {
			if shrinkHorzCount > 0 {
				change = diff / shrinkHorzCount
			}
		} else {
			if growHorzCount > 0 {
				change = diff / growHorzCount
			}
		}

		//        log.Stdoutf("*BoxLayout.Update: widgetCount: %d, cb: %+v, prefSizeSum: %+v, diff: %d, change: %d, reqH: %d", widgetCount, cb, prefSizeSum, diff, change, reqH)

		x := cb.X + l.margins.Left
		for i := 0; i < widgetCount; i++ {
			widget := widgets[i]

			w := prefSizes[i].Width

			switch {
			case change < 0:
				if flags[i]&ShrinkHorz > 0 {
					w += change
				}

			case change > 0:
				if flags[i]&GrowHorz > 0 {
					w += change
				}
			}

			bounds := drawing.Rectangle{x, cb.Y + l.margins.Top, w, reqH}

			//            log.Stdoutf("*BoxLayout.Update: bounds: %+v", bounds)

			widget.SetBounds(bounds)

			x += w + l.spacing
		}
	}

	return
}
