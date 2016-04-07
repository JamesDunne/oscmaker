// main
package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type ts struct {
	name  string
	color string
	track int
}

var track_setup = []ts{
	{"Master", "red", 1},
	{"Vox MG", "yellow", 2},
	{"Vox JD", "yellow", 3},
	{"Vox AS", "yellow", 4},
	{"Gtr MG", "purple", 6},
	{"Gtr JD", "purple", 7},
	{"Bass", "purple", 8},
	{"Kick", "green", 10},
	{"Snare", "green", 11},
	{"Overheads", "green", 12},
}

var colors = Colors{
	ForeColor:   -50116,
	BackColor:   -14803426,
	TextColor:   -1,
	BorderColor: -12566464,
	ForeAlpha:   255,
	BackAlpha:   255,
	TextAlpha:   255,
	BorderAlpha: 255,
}
var colors2 = Colors{
	ForeColor:   16727100,
	BackColor:   1973790,
	TextColor:   -1775042766,
	BorderColor: 16711422,
	ForeAlpha:   0,
	BackAlpha:   0,
	TextAlpha:   150,
	BorderAlpha: 0,
}

var colors_for = map[string]Colors{
	"red":    colors,
	"yellow": colors,
	"purple": colors,
	"green":  colors,
}

func createLayout(bank int, layoutname string, mixname string, masterlabel string) {
	var layout Layout

	// Create layout:
	layout.Version = layout_version
	layout.LayoutName = layoutname + " - One In Ten"
	layout.LayoutSummary = layoutname + " mix"
	layout.Width = 1280
	layout.Height = 752
	layout.Rotation = 1
	layout.Controls = make([]*Control, 0, 1)

	bank_track := (bank - 1) * 12
	const spacing = 128

	for t, ts := range track_setup {
		track := ts.track + bank_track

		// Fader:
		c := newFader(fmt.Sprintf("/track/%d/volume", track), ts.name)
		c.X = t*spacing + 0
		c.Y = 130
		c.Width = 120
		c.Height = 570
		c.SmoothingFactor = 12.0
		c.Colors = colors_for[ts.color]
		c.IsTouchable = true
		c.IsRelative = true
		layout.Controls = append(layout.Controls, c)

		// VU L:
		c = newFader(fmt.Sprintf("/track/%d/vu", track), "VU L")
		c.X = t*spacing + 0
		c.Y = 107
		c.Width = 13
		c.Height = 570
		layout.Controls = append(layout.Controls, c)

		// VU R:
	}

	// Dump layout XML to stdout:
	b, err := xml.MarshalIndent(&layout, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b))

	func() {
		of, err := os.OpenFile(layoutname+".oc.xml", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		// Write XML document:
		of.Write(b)
		err = of.Close()
		if err != nil {
			panic(err)
		}
	}()
}

func main() {
	createLayout(1, "PA", "PA System", "Master")
	//	createLayout(2, "JD", "Monitor for JD", "Master")
	//	createLayout(3, "MG", "Monitor for MG", "Master")
	//	createLayout(4, "MB", "Monitor for MB", "Master")
	//	createLayout(5, "AS", "Monitor for AS", "Master")
}
