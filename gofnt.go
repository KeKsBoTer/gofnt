// Copyright  2017 Simon Niedermayr.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Parsing font files for further use

package gofnt

import (
	"strings"
	"reflect"
	"strconv"
	"errors"
)

type Font struct {
	Info   Info   `json:"info"`
	Common Common `json:"common"`
	Pages  []Page `json:"page"`
	Chars  []Char `json:"char"`
}

type Info struct {
	Name     string `json:"face"`
	Size     int    `json:"size"`
	Bold     bool   `json:"bold"`
	Italic   bool   `json:"italic"`
	Charset  string `json:"charset"`
	Unicode  bool   `json:"unicode"`
	StretchH int    `json:"stretchH"`
	Smooth   bool   `json:"smooth"`
	AA       int    `json:"aa"`
	Padding  [4]int `json:"padding"`
	Spacing  [2]int `json:"spacing"`
	Outline  int    `json:"outline"`
}

type Common struct {
	LineHeight   int  `json:"lineHeight"`
	Base         int  `json:"base"`
	ScaleW       int  `json:"scaleW"`
	ScaleH       int  `json:"scaleH"`
	Pages        int  `json:"pages"`
	Packed       bool `json:"packed"`
	AlphaChannel int8 `json:"alphaChnl"`
	RedChannel   int8 `json:"redChnl"`
	GreenChannel int8 `json:"greenChnl"`
	BlueChannel  int8 `json:"blueChnl"`
}

type Page struct {
	Id   int    `json:"id"`
	File string `json:"file"`
}

type Char struct {
	Id        int32 `json:"id"`
	X         int `json:"x"`
	Y         int `json:"y"`
	Width     int `json:"width"`
	Height    int `json:"height"`
	XOffset   int `json:"xoffset"`
	YOffset   int `json:"yoffset"`
	XAdvanced int `json:"xadvance"`
	Page      int `json:"page"`
	Chnl      int `json:"chnl"`
}

// This function takes a font file as string and converts it into a Font struct.
func Parse(file string) (*Font, error) {
	font := Font{}
	textLines := strings.Split(file, "\n")
	lines := map[string][]map[string]string{}
	for _, l := range textLines {
		//TODO allow space in quotation marks
		cells := strings.Split(l, " ")
		if len(cells) > 1 {
			index := 0
			if lines[cells[0]] == nil {
				lines[cells[0]] = make([]map[string]string, 1)
			} else {
				lines[cells[0]] = append(lines[cells[0]], make(map[string]string))
				index = len(lines[cells[0]]) - 1
			}
			lines[cells[0]][index] = map[string]string{}
			for j := 1; j < len(cells); j++ {
				if len(cells[j]) < 3 {
					continue
				}
				key, value := parsePair(cells[j])
				if key == "" {
					continue
				}
				lines[cells[0]][index][key] = value
			}
		}
	}
	fontType := reflect.ValueOf(&font)
	elm := fontType.Elem()
	for i := 0; i < elm.NumField(); i++ {
		field := elm.Field(i)
		fieldTag := elm.Type().Field(i).Tag
		if tag, exists := fieldTag.Lookup("json"); exists {
			if field.Kind() == reflect.Slice {
				unmarshalSlice(field, lines[tag])
			} else {
				unmarshal(field, lines[tag][0])
			}
		}
	}
	return &font, nil //TOTO error handling
}

// This function takes a list of maps and copies the values in the given slice
// The field needs to be of type slice
// See function unmarshal for closer description
func unmarshalSlice(field reflect.Value, values []map[string]string) {
	buffer := reflect.MakeSlice(field.Type(), len(values), len(values))
	for i, v := range values {
		newInstance := reflect.New(field.Type().Elem()).Elem()
		unmarshal(newInstance, v)
		buffer.Index(i).Set(newInstance)
	}
	field.Set(reflect.AppendSlice(field, buffer))
}

// This function copies the values from the map into the corresponding struct fields in the given field.
// Which value belongs to which field in the struct is determined by the "json" tag of a field.
// e.g.
// type Test struct{
// 		field1 string `json="fieldOne"`
// }
// testMap := map[string]string{"fieldOne":"valueOne"}
//
// In this example the value "valueOne" would be copied to the field one in Test
func unmarshal(field reflect.Value, values map[string]string) {
	for j := 0; j < field.NumField(); j++ {
		subField := field.Field(j)
		subFieldTag := field.Type().Field(j).Tag
		if name, exists := subFieldTag.Lookup("json"); exists {
			copyValue(subField, []byte(values[name]))
		}

	}
}

// This functions extracts a key-value-pair from a string.
// The key and value must be separated by a "="
// The first returned value is the key and the second one the value
// e.g The string "key=value" returns ("key","value")
func parsePair(pair string) (string, string) {
	noBreak := strings.Replace(pair, "\r", "", -1)
	split := strings.Split(noBreak, "=")
	if len(split) != 2 {
		return "", ""
	}
	valueLen := len(split[1])
	// Remove quotation marks
	if split[1][0] == '"' && split[1][valueLen-1] == '"' {
		split[1] = split[1][1:valueLen-1]
	}
	return split[0], split[1]
}

// copied from golang.org/encoding/xml
func copyValue(dst reflect.Value, src []byte) (err error) {
	dst0 := dst

	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}

	// Save accumulated data.
	switch dst.Kind() {
	case reflect.Invalid:
		// Probably a comment.
	default:
		return errors.New("cannot unmarshal into " + dst0.Type().String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(src) == 0 {
			dst.SetInt(0)
			return nil
		}
		itmp, err := strconv.ParseInt(string(src), 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetInt(itmp)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if len(src) == 0 {
			dst.SetUint(0)
			return nil
		}
		utmp, err := strconv.ParseUint(string(src), 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetUint(utmp)
	case reflect.Float32, reflect.Float64:
		if len(src) == 0 {
			dst.SetFloat(0)
			return nil
		}
		ftmp, err := strconv.ParseFloat(string(src), dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetFloat(ftmp)
	case reflect.Bool:
		if len(src) == 0 {
			dst.SetBool(false)
			return nil
		}
		value, err := strconv.ParseBool(strings.TrimSpace(string(src)))
		if err != nil {
			return err
		}
		dst.SetBool(value)
	case reflect.String:
		dst.SetString(string(src))
	case reflect.Slice:
		if len(src) == 0 {
			// non-nil to flag presence
			src = []byte{}
		}
		dst.SetBytes(src)
	case reflect.Array:
		values := strings.Split(string(src), ",")
		for i := 0; i < dst.Len(); i++ {
			copyValue(dst.Index(i), []byte(values[i]))
		}
	}
	return nil
}
