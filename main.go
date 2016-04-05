// main
package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/xml"
	// "fmt"
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

// <control name="bWFzdGVyL2xldmVs" x="66" y="40" w="500" h="100" color="red"
//   scalef="0.0" scalet="1.0" type="faderh" response="absolute" inverted="false"
//   centered="false"></control>

// <control name="bGFiZWwz" x="21" y="160" w="45" h="100" color="yellow"
//   type="labelv" text="Vm94IE1H" size="14" background="false" outline="true"></control>

// <control name="bWFzdGVyL3Bhbg==" x="566" y="40" w="100" h="100" color="red"
//   scalef="0.0" scalet="1.0" type="rotaryh" response="absolute" inverted="false"
//   centered="true" norollover="true" ></control>

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

func main() {
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

	var layout Layout
	err = xml.Unmarshal(b, &layout)
	if err != nil {
		panic(err)
	}

	// Base64 decode Name properties:
	for pi := range layout.TabPages {
		p := &layout.TabPages[pi]
		p.Name = b64decode(p.Name)
		for ci := range p.Controls {
			c := &p.Controls[ci]
			c.Name = b64decode(c.Name)
		}
	}

	// Copy template TabPage:
	tp := layout.TabPages[0]

	layout.TabPages = make([]TabPage, 0, 5)
	for i := 0; i < 5; i++ {
		p := tp
		p.Controls = make([]Control, len(tp.Controls))
		copy(p.Controls[:], tp.Controls)

		layout.TabPages = append(layout.TabPages, p)
	}

	layout.TabPages[1].Name = "JD"
	layout.TabPages[2].Name = "MG"
	layout.TabPages[3].Name = "MB"
	layout.TabPages[4].Name = "AS"

	// Base64 encode Name properties:
	for pi := range layout.TabPages {
		p := &layout.TabPages[pi]
		p.Name = b64encode(p.Name)
		for ci := range p.Controls {
			c := &p.Controls[ci]
			c.Name = b64encode(c.Name)
		}
	}

	//b, err = xml.MarshalIndent(&layout, "", "  ")
	b, err = xml.Marshal(&layout)
	//fmt.Printf("%s\n", string(b))

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
