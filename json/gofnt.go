package json

import (
	"fmt"
	"errors"
	"encoding/json"
)

// allowing numbers 0 and 1 to be interpreted as bool
type jbool bool
type jfloat = float32

func (bit jbool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		bit = true
	} else if asString == "0" || asString == "false" {
		bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}

func ParseJSON(file []byte) (*Font, error) {
	fnt := new(Font)
	err := json.Unmarshal(file, fnt)

	return fnt, err
}

type Font struct {
	Info   Info     `json:"info"`
	Common Common   `json:"common"`
	Pages  []string `json:"pages"`
	Chars  []Char   `json:"chars"`
}

type Info struct {
	Name     string   `json:"face"`
	Size     int      `json:"size"`
	Bold     jbool    `json:"bold"`
	Italic   jbool    `json:"italic"`
	Charset  []string `json:"charset"`
	Unicode  jbool    `json:"unicode"`
	StretchH int      `json:"stretchH"`
	Smooth   jbool    `json:"smooth"`
	AA       int      `json:"aa"`
	Padding  [4]int   `json:"padding"`
	Spacing  [2]int   `json:"spacing"`
	Outline  int      `json:"outline"`
}

type Common struct {
	LineHeight   jfloat `json:"lineHeight"`
	Base         jfloat `json:"base"`
	ScaleW       int    `json:"scaleW"`
	ScaleH       int    `json:"scaleH"`
	Pages        int    `json:"pages"`
	Packed       jbool  `json:"packed"`
	AlphaChannel int8   `json:"alphaChnl"`
	RedChannel   int8   `json:"redChnl"`
	GreenChannel int8   `json:"greenChnl"`
	BlueChannel  int8   `json:"blueChnl"`
}
type Page struct {
	Id   int    `json:"id"`
	File string `json:"file"`
}

type Char struct {
	Id        int32  `json:"id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	XOffset   jfloat `json:"xoffset"`
	YOffset   jfloat `json:"yoffset"`
	XAdvanced jfloat `json:"xadvance"`
	Page      int    `json:"page"`
	Chnl      int    `json:"chnl"`
}
