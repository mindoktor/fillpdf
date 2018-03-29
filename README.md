# FillPDF

FillPDF is a golang library to easily fill PDF forms. This library uses the pdftk utility to fill the PDF forms with fdf data.
Currently this library only supports PDF text and checkbox field values. Feel free to add support to more form types (Send pull request to original developer)
This fork extends with some more pdftk commands
* Multistamp
* Ability to generate PDF's with special characters (with flatten) with pdftk. (Limited by font in PDF)

## Documentation 

Check the Documentation at [GoDoc.org](https://godoc.org/github.com/desertbit/fillpdf).


## Sample

There is an example:

```go
package main

import (
	"log"

	"github.com/phelian/fillpdf"
)

func main() {
	// Create the form values.
	form := fillpdf.Form{
		"field_1": "Hello",
		"field_2": "World",
		"field_3": "with special åäöéè",
	}

	// Fill the form PDF with our values.
	err := fillpdf.Fill(form, "form.pdf", "filled.pdf", "On", "Off", true)
	if err != nil {
		log.Fatal(err)
	}
}
```

Run the example as following:

```
cd sample
go build
./sample
```
