package main

import "encoding/xml"

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
	hasValue bool
	value    float64
}

func HaveFloat(value float64) MaybeFloat {
	return MaybeFloat{hasValue: true, value: value}
}

func NoFloat() MaybeFloat {
	return MaybeFloat{}
}

func (f MaybeFloat) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !f.hasValue {
		return nil
	}
	return e.EncodeElement(f.value, start)
}

type OSCArgument struct {
	XMLName xml.Name `xml:"OSCArgument"`

	Version string `xml:"version,attr"`

	ActiveAxes  ActiveAxes  `xml:"activeAxes"`
	ScalingAxes ScalingAxes `xml:"scalingAxes"`

	ValueOFloat   MaybeFloat `xml:"valueOFloat"`
	ValueOfString *string    `xml:"valueOfString,omitempty"`

	RefOFloat   MaybeFloat `xml:"refOFloat"`
	RefOfString *string    `xml:"refOfString,omitempty"`

	MaxOFloat   MaybeFloat `xml:"maxOFloat"`
	MaxOfString *string    `xml:"maxOfString,omitempty"`

	DefOfString *string `xml:"defOfString,omitempty"`

	MinOFloat   MaybeFloat `xml:"minOFloat"`
	MinOfString *string    `xml:"minOfString,omitempty"`

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
		DisplayName:     false,
		LocalFeedback:   true,
		IsSliding:       false,
		IsRelative:      false,
		ContForeImg:     false,
		Font:            "Dialog-plain-15",
		BackgroundImage: " ",
		ForegroundImage: " ",
	}
}

const (
	layout_version   = "1.3"
	control_version  = "1.8"
	argument_version = "1.2"
)

const (
	reaper_zero_dB = 0.716878
)

var (
	fader        = xml.Name{Local: "fader"}
	toggleButton = xml.Name{Local: "toggleButton"}
	pad          = xml.Name{Local: "pad"}

	activeAxesNone = ActiveAxes{
		Axes: []AxisActive{
			AxisActive{
				Number: 1,
				Active: false,
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

	argFloat = OSCArgument{
		Version:         argument_version,
		ActiveAxes:      activeAxes1,
		ScalingAxes:     scalingAxes3,
		ValueOFloat:     HaveFloat(1.0),
		RefOFloat:       HaveFloat(0.0),
		MaxOFloat:       HaveFloat(reaper_zero_dB),
		DefOfString:     strPtr("1"),
		MinOFloat:       HaveFloat(0.0),
		IsListening:     true,
		DisplayOnWidget: false,
	}

	argString = OSCArgument{
		Version:         argument_version,
		ActiveAxes:      activeAxesNone,
		ScalingAxes:     scalingAxes3,
		ValueOfString:   strPtr("On"),
		RefOfString:     strPtr("Off"),
		MaxOfString:     strPtr("On"),
		DefOfString:     strPtr("1"),
		MinOfString:     strPtr("Off"),
		IsListening:     false,
		DisplayOnWidget: false,
	}
)

func newFader(osc string, text string) *Control {
	c := newControl(fader)
	c.Text = text
	c.OSCBundle = OSCBundle{
		[]OSCMessage{
			OSCMessage{
				OSCAddress: osc,
				OSCArguments: []OSCArgument{
					argFloat,
				},
			},
		},
	}
	return c
}

func newPad(text string) *Control {
	c := newControl(fader)
	c.Text = text
	c.OSCBundle = OSCBundle{}
	return c
}

func newToggleButton(osc string, text string, arg OSCArgument) *Control {
	c := newControl(toggleButton)
	c.Text = text

	c.OSCBundle = OSCBundle{
		[]OSCMessage{
			OSCMessage{
				OSCAddress: osc,
				OSCArguments: []OSCArgument{
					arg,
				},
			},
		},
	}
	return c
}
