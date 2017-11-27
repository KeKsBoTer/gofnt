# gofnt
A Go library that reads bitmap-font files, described by
http://www.angelcode.com/products/bmfont/doc/file_format.html

Font files can be created using:
https://github.com/libgdx/libgdx/wiki/Hiero

## How To Use
```
content, err := ioutil.ReadFile("path_to_font.fnt")
if err != nil{
		log.Fatalln(err)
}
font := gofnt.Parse(string(content))
```
