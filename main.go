// main
package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	inFile  = "template.touchosc"
	outFile = "oneinten.touchosc"
)

func b64decode(e string) string {
	b, err := base64.StdEncoding.DecodeString(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func b64encode(s string) string {
	e := base64.StdEncoding.EncodeToString([]byte(s))
	return e
}

type PageCustomization struct {
	PageName        *string
	MasterFaderName **string
	MixName         **string
}

type Control struct {
	XMLName xml.Name `xml:"control"`

	Name string `xml:"name,attr"` // base64
	Type string `xml:"type,attr"`

	X     int    `xml:"x,attr"`
	Y     int    `xml:"y,attr"`
	W     int    `xml:"w,attr"`
	H     int    `xml:"h,attr"`
	Color string `xml:"color,attr"`

	ScaleF *float64 `xml:"scalef,attr,omitempty"`
	ScaleT *float64 `xml:"scalet,attr,omitempty"`

	// faderh
	Response *string `xml:"response,attr,omitempty"`
	Inverted *bool   `xml:"inverted,attr,omitempty"`
	Centered *bool   `xml:"centered,attr,omitempty"`

	// labelv
	Text       *string `xml:"text,attr,omitempty"`
	Size       *int    `xml:"size,attr,omitempty"`
	Background *bool   `xml:"background,attr,omitempty"`
	Outline    *bool   `xml:"outline,attr,omitempty"`

	// rotaryh
	NoRollover *bool `xml:"norollover,attr,omitempty"`
}

type TabPage struct {
	XMLName xml.Name `xml:"tabpage"`

	Name     string    `xml:"name,attr"` // base64
	ScaleF   float64   `xml:"scalef,attr"`
	ScaleT   float64   `xml:"scalet,attr"`
	Controls []Control `xml:"control"`
}

type Layout struct {
	XMLName xml.Name `xml:"layout"`

	Version     int       `xml:"version,attr"`
	Mode        int       `xml:"mode,attr"`
	Width       int       `xml:"w,attr"`
	Height      int       `xml:"h,attr"`
	Orientation string    `xml:"orientation,attr"`
	TabPages    []TabPage `xml:"tabpage"`
}

func intPtr(i int) *int             { return &i }
func strPtr(s string) *string       { return &s }
func boolPtr(b bool) *bool          { return &b }
func float64Ptr(f float64) *float64 { return &f }

func makePage(tp *TabPage, bank int, pagename string, mixname string, masterlabel string) {
	const tracks = 10
	const controls = 3

	const fader_min = 0.483
	const fader_max_0 = 0.716878
	const fader_max = 0.86

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

	tp.Name = pagename
	tp.ScaleF = 0.0
	tp.ScaleT = 1.0
	tp.Controls = make([]Control, 0, tracks*controls)

	// Mix name label:
	tp.Controls = append(tp.Controls, Control{
		Name:       "mixname",
		Type:       "labelv",
		X:          672,
		Y:          483,
		W:          40,
		H:          260,
		Color:      "orange",
		Text:       strPtr(mixname),
		Size:       intPtr(30),
		Background: boolPtr(false),
		Outline:    boolPtr(false),
	})

	bank_track := (bank - 1) * 12
	track := track_setup[0].track + bank_track

	// Master label
	tp.Controls = append(tp.Controls, Control{
		Name:       fmt.Sprintf("%d/label", track),
		Type:       "labelv",
		X:          20,
		Y:          50,
		W:          40,
		H:          100,
		Color:      track_setup[0].color,
		Text:       strPtr(masterlabel),
		Size:       intPtr(20),
		Background: boolPtr(false),
		Outline:    boolPtr(true),
	})

	// Master fader:
	tp.Controls = append(tp.Controls, Control{
		Name:     fmt.Sprintf("%d/volume", track),
		Type:     "faderh",
		X:        66,
		Y:        50,
		W:        600,
		H:        100,
		Color:    track_setup[0].color,
		ScaleF:   float64Ptr(0.0),
		ScaleT:   float64Ptr(fader_max_0),
		Response: strPtr("absolute"),
		Inverted: boolPtr(false),
		Centered: boolPtr(false),
	})

	// Controls per track:
	for i := 1; i < tracks; i++ {
		track := track_setup[i].track + bank_track
		// Label
		tp.Controls = append(tp.Controls, Control{
			Name:       fmt.Sprintf("%d/label", track),
			Type:       "labelv",
			X:          20,
			Y:          50 + (i * 120),
			W:          40,
			H:          100,
			Color:      track_setup[i].color,
			Text:       strPtr(track_setup[i].name),
			Size:       intPtr(20),
			Background: boolPtr(false),
			Outline:    boolPtr(true),
		})

		// Fader
		tp.Controls = append(tp.Controls, Control{
			Name:     fmt.Sprintf("%d/volume", track),
			Type:     "faderh",
			X:        66,
			Y:        50 + (i * 120),
			W:        600,
			H:        100,
			Color:    track_setup[i].color,
			ScaleF:   float64Ptr(fader_min),
			ScaleT:   float64Ptr(fader_max),
			Response: strPtr("absolute"),
			Inverted: boolPtr(false),
			Centered: boolPtr(false),
		})

		//// Pan
		//tp.Controls = append(tp.Controls, Control{
		//	Name:       fmt.Sprintf("%d/pan", track),
		//	Type:       "rotaryh",
		//	X:          566,
		//	Y:          50 + (i * 120),
		//	W:          100,
		//	H:          100,
		//	Color:      "yellow",
		//	ScaleF:     float64Ptr(0.25),
		//	ScaleT:     float64Ptr(0.75),
		//	Response:   strPtr("absolute"),
		//	Inverted:   boolPtr(false),
		//	Centered:   boolPtr(true),
		//	NoRollover: boolPtr(true),
		//})
	}
}

func translateTemplate(layout *Layout) {
	// Open ZIP file:
	zf, err := zip.OpenReader(inFile)
	if err != nil {
		panic(err)
	}

	var xf io.ReadCloser
	for _, f := range zf.File {
		if filepath.Ext(f.Name) != ".xml" {
			continue
		}

		xf, err = f.Open()
		if err != nil {
			panic(err)
		}
		break
	}

	// Read all of index.xml file:
	b, err := ioutil.ReadAll(xf)
	if err != nil {
		panic(err)
	}
	xf.Close()
	zf.Close()

	err = xml.Unmarshal(b, layout)
	if err != nil {
		panic(err)
	}

	pc := make([]PageCustomization, 5)

	// Base64 decode Name properties:
	for pi := range layout.TabPages {
		p := &layout.TabPages[pi]
		p.Name = b64decode(p.Name)
		for ci := range p.Controls {
			c := &p.Controls[ci]
			c.Name = b64decode(c.Name)
			if c.Text != nil {
				*c.Text = b64decode(*c.Text)
			}
		}
	}

	// Copy template TabPage:
	tp := layout.TabPages[0]

	layout.TabPages = make([]TabPage, 0, 5)
	for pi := 0; pi < 5; pi++ {
		tmp := tp
		tmp.Controls = make([]Control, len(tp.Controls))
		copy(tmp.Controls[:], tp.Controls)
		layout.TabPages = append(layout.TabPages, tmp)

		p := &layout.TabPages[pi]

		// Create the PageCustomization:
		pc[pi].PageName = &p.Name
		for ci := range p.Controls {
			if p.Controls[ci].Name == "master/label" {
				pc[pi].MasterFaderName = &p.Controls[ci].Text
			} else if p.Controls[ci].Name == "mixname" {
				pc[pi].MixName = &p.Controls[ci].Text
			}
		}
	}

	*pc[0].PageName = "PA"
	*pc[0].MixName = strPtr("PA System")
	*pc[0].MasterFaderName = strPtr("PA Master")

	*pc[1].PageName = "JD"
	*pc[1].MixName = strPtr("Monitor for JD")
	*pc[1].MasterFaderName = strPtr("JD Master")

	*pc[2].PageName = "MG"
	*pc[2].MixName = strPtr("Monitor for MG")
	*pc[2].MasterFaderName = strPtr("MG Master")

	*pc[3].PageName = "MB"
	*pc[3].MixName = strPtr("Monitor for MB")
	*pc[3].MasterFaderName = strPtr("MB Master")

	*pc[4].PageName = "AS"
	*pc[4].MixName = strPtr("Monitor for AS")
	*pc[4].MasterFaderName = strPtr("AS Master")
}

func createLayout(layout *Layout) {
	// Create layout:
	layout.Version = 15
	layout.Mode = 3
	layout.Width = 752
	layout.Height = 1280
	layout.Orientation = "vertical"

	layout.TabPages = make([]TabPage, 5)
	makePage(&layout.TabPages[0], 1, "PA", "PA System", "Master")
	makePage(&layout.TabPages[1], 2, "JD", "Monitor for JD", "Master")
	makePage(&layout.TabPages[2], 3, "MG", "Monitor for MG", "Master")
	makePage(&layout.TabPages[3], 4, "MB", "Monitor for MB", "Master")
	makePage(&layout.TabPages[4], 5, "AS", "Monitor for AS", "Master")
}

func main() {
	var layout Layout

	//translateTemplate(&layout)

	createLayout(&layout)

	// Dump layout XML to stdout:
	b, err := xml.MarshalIndent(&layout, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b))

	// Base64 encode Name properties:
	for pi := range layout.TabPages {
		p := &layout.TabPages[pi]
		p.Name = b64encode(p.Name)
		for ci := range p.Controls {
			c := &p.Controls[ci]
			c.Name = b64encode(c.Name)
			if c.Text != nil {
				tmp := b64encode(*c.Text)
				c.Text = &tmp
			}
		}
	}

	b, err = xml.Marshal(&layout)

	func() {
		ozf, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer ozf.Close()
		zw := zip.NewWriter(ozf)
		xw, err := zw.Create("index.xml")
		if err != nil {
			panic(err)
		}
		// Write XML header:
		xw.Write([]byte(xml.Header))
		// Write XML document:
		xw.Write(b)
		err = zw.Close()
		if err != nil {
			panic(err)
		}
	}()
}
