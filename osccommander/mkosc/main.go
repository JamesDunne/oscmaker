// main
package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

const (
	outFile = "OneInTen.oc.xml"
)

type AxisScale struct {
	XMLName  xml.Name `xml:"axis"`
	Number   int      `xml:"number"`
	Function int      `xml:"function"`
}

type AxisActive struct {
	XMLName xml.Name `xml:"axis"`
	Number  int      `xml:"number"`
	Active  bool     `xml:"active"`
}

type ActiveAxes struct {
	Axes []AxisActive
}
type ScalingAxes struct {
	Axes []AxisScale
}

type MaybeFloat struct {
	float64
}

func (f *MaybeFloat) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if f == nil {
		return nil
	}
	e.Encode(f.float64)
	return nil
}

type OSCArgument struct {
	Version string `xml:"version,attr"`

	ActiveAxes  ActiveAxes  `xml:"activeAxes"`
	ScalingAxes ScalingAxes `xml:"scalingAxes"`

	ValueOFloat   *MaybeFloat `xml:"valueOFloat,omitempty"`
	ValueOfString *string     `xml:"valueOfString,omitempty"`

	RefOFloat   *float64 `xml:"refOFloat,omitempty"`
	RefOfString *string  `xml:"refOfString,omitempty"`

	MaxOFloat   *float64 `xml:"maxOFloat,omitempty"`
	MaxOfString *string  `xml:"maxOfString,omitempty"`

	DefOfString *string `xml:"defOfString,omitempty"`

	MinOFloat   *float64 `xml:"minOFloat,omitempty"`
	MinOfString *string  `xml:"minOfString,omitempty"`

	IsListening     bool `xml:"isListening"`
	DisplayOnWidget bool `xml:"displayOnWidget"`
}

type OSCMessage struct {
	OSCAddress   string        `xml:"OSCAddress"`
	OSCArguments []OSCArgument `xml:",any"`
}

type OSCBundle struct {
	OSCMessages []OSCMessage `xml:"OSCMessage"`
}

type Colors struct {
	ForeColor   int `xml:"foreColor"`
	BackColor   int `xml:"backColor"`
	TextColor   int `xml:"textColor"`
	BorderColor int `xml:"borderColor"`
	ForeAlpha   int `xml:"foreAlpha"`
	BackAlpha   int `xml:"backAlpha"`
	TextAlpha   int `xml:"textAlpha"`
	BorderAlpha int `xml:"borderAlpha"`
}

type Control struct {
	XMLName xml.Name

	Version string `xml:"version,attr"`
	Text    string `xml:"text,attr"`

	OSCBundle OSCBundle `xml:"OSCBundle"`

	Colors Colors `xml:"colors"`

	X               int     `xml:"X"`
	Y               int     `xml:"Y"`
	Z               int     `xml:"Z"`
	Width           int     `xml:"width"`
	Height          int     `xml:"height"`
	Borderwidth     int     `xml:"borderwidth"`
	SmoothingFactor float64 `xml:"smoothingFactor"`
	Rotation        float64 `xml:"rotation"`
	IsTouchable     bool    `xml:"isTouchable"`
	DisplayName     bool    `xml:"displayName"`
	LocalFeedback   bool    `xml:"localFeedback"`
	IsSliding       bool    `xml:"isSliding"`
	IsRelative      bool    `xml:"isRelative"`
	ContForeImg     bool    `xml:"contForeImg"`
	Font            string  `xml:"font"`
	BackgroundImage string  `xml:"backgroundImage"`
	ForegroundImage string  `xml:"foregroundImage"`
}

type Layout struct {
	XMLName xml.Name `xml:"main"`
	Version string   `xml:"version,attr"`

	LayoutName    string `xml:"layoutName"`
	LayoutSummary string `xml:"layoutSummary"`
	Width         int    `xml:"width"`
	Height        int    `xml:"height"`
	Rotation      int    `xml:"rotation"`

	Controls []*Control `xml:",any"`
}

func intPtr(i int) *int             { return &i }
func strPtr(s string) *string       { return &s }
func boolPtr(b bool) *bool          { return &b }
func float64Ptr(f float64) *float64 { return &f }

func newControl(name xml.Name) *Control {
	return &Control{
		XMLName: name,
		Version: control_version,
		Colors: Colors{
			ForeColor:   16727100,
			BackColor:   1973790,
			TextColor:   -1775042766,
			BorderColor: 16711422,
			ForeAlpha:   0,
			BackAlpha:   0,
			TextAlpha:   150,
			BorderAlpha: 0,
		},
		Z:               0,
		Width:           13,
		Height:          25,
		Borderwidth:     0,
		SmoothingFactor: 0.0,
		Rotation:        0.0,
		IsTouchable:     false,
		DisplayName:     true,
		LocalFeedback:   true,
		IsSliding:       false,
		IsRelative:      false,
		ContForeImg:     false,
		Font:            "Dialog-plain-36",
		BackgroundImage: " ",
		ForegroundImage: " ",
	}
}

const (
	control_version  = "1.8"
	argument_version = "1.2"
)

var (
	fader = xml.Name{Local: "fader"}
	pad   = xml.Name{Local: "pad"}

	activeAxes1 = ActiveAxes{
		Axes: []AxisActive{
			AxisActive{
				Number: 1,
				Active: true,
			},
			AxisActive{
				Number: 2,
				Active: false,
			},
			AxisActive{
				Number: 3,
				Active: false,
			},
		},
	}

	scalingAxes3 = ScalingAxes{
		Axes: []AxisScale{
			AxisScale{
				Number:   1,
				Function: 3,
			},
			AxisScale{
				Number:   2,
				Function: 3,
			},
			AxisScale{
				Number:   3,
				Function: 3,
			},
		},
	}
)

func newFader(osc string, text string) *Control {
	c := newControl(fader)
	c.Text = text
	c.OSCBundle = OSCBundle{
		[]OSCMessage{
			OSCMessage{
				OSCAddress: "/track/1/volume",
				OSCArguments: []OSCArgument{
					OSCArgument{
						Version:         argument_version,
						ActiveAxes:      activeAxes1,
						ScalingAxes:     scalingAxes3,
						ValueOFloat:     &MaybeFloat{1.0},
						RefOFloat:       float64Ptr(0.0),
						MaxOFloat:       float64Ptr(0.0),
						DefOfString:     strPtr("1"),
						MinOFloat:       float64Ptr(0.0),
						IsListening:     true,
						DisplayOnWidget: false,
					},
				},
			},
		},
	}
	return c
}

func createLayout(layout *Layout) {
	// Create layout:
	layout.Version = "1.3"
	layout.LayoutName = "One In Ten - PA"
	layout.LayoutSummary = "PA mixer"
	layout.Width = 752
	layout.Height = 1280
	layout.Rotation = 1
	layout.Controls = make([]*Control, 0, 1)

	var c *Control
	c = newFader("/track/1/volume", "Vol 1")
	c.Width = 150
	c.Height = 560
	layout.Controls = append(layout.Controls, c)

}

func main() {
	var layout Layout

	createLayout(&layout)

	// Dump layout XML to stdout:
	b, err := xml.MarshalIndent(&layout, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b))

	func() {
		of, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		// Write XML header:
		//xw.Write([]byte(xml.Header))
		// Write XML document:
		of.Write(b)
		err = of.Close()
		if err != nil {
			panic(err)
		}
	}()
}
