// main
package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

var colors_for = map[string]Colors{
	"red": Colors{
		ForeColor:   -50116,
		BackColor:   -14803426,
		TextColor:   -1,
		BorderColor: -12566464,
		ForeAlpha:   255,
		BackAlpha:   255,
		TextAlpha:   255,
		BorderAlpha: 255,
	},
	"yellow": Colors{
		ForeColor:   -256,
		BackColor:   -14803456,
		TextColor:   -1,
		BorderColor: -12566464,
		ForeAlpha:   255,
		BackAlpha:   255,
		TextAlpha:   255,
		BorderAlpha: 255,
	},
	"purple": Colors{
		ForeColor:   -3392513,
		BackColor:   -15597543,
		TextColor:   -1,
		BorderColor: -12566464,
		ForeAlpha:   255,
		BackAlpha:   255,
		TextAlpha:   255,
		BorderAlpha: 255,
	},
	"green": Colors{
		ForeColor:   -13382656,
		BackColor:   -16377600,
		TextColor:   -1,
		BorderColor: -12566464,
		ForeAlpha:   255,
		BackAlpha:   255,
		TextAlpha:   255,
		BorderAlpha: 255,
	},
	"blue": Colors{
		ForeColor:   -16751361,
		BackColor:   -14803426,
		TextColor:   -1,
		BorderColor: -16751361,
		ForeAlpha:   255,
		BackAlpha:   255,
		TextAlpha:   255,
		BorderAlpha: 255,
	},
}

type ts struct {
	name  string
	color string
	track int
}

var track_setup = [10]ts{
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

func createLayout(bank int, layoutname string, mixname string, masterlabel string) {
	var layout Layout

	// Create layout:
	layout.Version = layout_version
	layout.LayoutName = fmt.Sprintf("%d - %s", bank, layoutname)
	layout.LayoutSummary = layoutname + " mix"
	layout.Width = 1280
	layout.Height = 752
	layout.Rotation = 1
	layout.Controls = make([]*Control, 0, 1)

	// Top label:
	c := newPad(mixname)
	c.X = 402
	c.Y = 0
	c.Z = 0
	c.Width = 470
	c.Height = 70
	c.Borderwidth = 0
	c.IsTouchable = false
	c.DisplayName = true
	c.Font = "Dialog-plain-38"
	layout.Controls = append(layout.Controls, c)

	bank_track := (bank - 1) * 12
	const spacing = 128

	// Set master label:
	track_setup[0].name = masterlabel

	for t, ts := range track_setup {
		track := ts.track + bank_track

		// Mute:
		c := newToggleButton(fmt.Sprintf("/track/%d/mute", track), "Mute", argFloat)
		c.X = t*spacing + 0
		c.Y = 68
		c.Width = 120
		c.Height = 60
		c.Colors = colors_for["blue"]
		c.Borderwidth = 2
		c.SmoothingFactor = 12.0
		c.IsTouchable = true
		c.DisplayName = true
		c.LocalFeedback = false
		c.IsSliding = false
		c.IsRelative = false
		layout.Controls = append(layout.Controls, c)

		// Fader:
		c = newFader(fmt.Sprintf("/track/%d/volume", track), ts.name)
		c.X = t*spacing + 0
		c.Y = 130
		c.Width = 120
		c.Height = 570
		c.SmoothingFactor = 12.0
		c.Colors = colors_for[ts.color]
		c.IsTouchable = true
		c.IsRelative = true
		// TODO: is Default misnamed as Max? Editor seems to confuse the two.
		c.OSCBundle.OSCMessages[0].OSCArguments[0].MaxOFloat = HaveFloat(reaper_zero_dB)
		layout.Controls = append(layout.Controls, c)

		// VU L:
		// VU R:

		// Label:
		c = newToggleButton(fmt.Sprintf("/track/%d/name", track), ts.name, argString)
		c.X = t*spacing + 0
		c.Y = 700
		c.Width = 120
		c.Height = 52
		c.Borderwidth = 1
		c.SmoothingFactor = 12.0
		c.IsTouchable = false
		c.DisplayName = true
		c.LocalFeedback = true
		c.Font = "Dialog-plain-20"
		layout.Controls = append(layout.Controls, c)
	}

	// Dump layout XML to stdout:
	b, err := xml.MarshalIndent(&layout, "", "   ")
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s\n", string(b))

	// Create XML file:
	of, err := os.OpenFile(fmt.Sprintf("OneInTen - %d - %s.oc.xml", bank, layoutname), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	// Write XML document:
	of.Write(b)
	err = of.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	createLayout(1, "PA", "PA System", "PA Master")
	createLayout(2, "JD", "Monitor for JD", "JD Master")
	createLayout(3, "MG", "Monitor for MG", "MG Master")
	createLayout(4, "MB", "Monitor for MB", "MB Master")
	createLayout(5, "AS", "Monitor for AS", "AS Master")
}
