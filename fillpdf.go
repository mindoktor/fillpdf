/*
 *  FillPDF - Fill PDF forms
 *  Copyright DesertBit
 *  Authors: Roland Singer, Alexander Félix
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"unicode/utf16"
)

// Form represents the PDF form.
// This is a key value map.
type Form map[string]interface{}

// createFdfFile with 16 bit encoded utf to enable creation of pdf with special characters
func createFdfFile(form Form, path, checkedString, uncheckedString string) error {
	// Create the file.
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new writer.
	b := bufio.NewWriter(file)

	// Header
	b.WriteString("%FDF-1.2\n")
	b.WriteString("\xE2\xE3\xCF\xD3\n")
	b.WriteString("1 0 obj \n")
	b.WriteString("<<\n")
	b.WriteString("/FDF \n")
	b.WriteString("<<\n")
	b.WriteString("/Fields [\n")

	// Write the form data.
	for key, value := range form {
		var valStr string
		switch v := value.(type) {
		case bool:
			if v {
				valStr = checkedString
			} else {
				valStr = uncheckedString
			}
		default:
			valStr = fmt.Sprintf("%v", value)
		}

		b.WriteString("<<\n")
		b.WriteString("/T <")
		b.WriteString(hex.EncodeToString(encodeUTF16(key, true)))
		b.WriteString(">\n")
		b.WriteString("/V <")
		b.WriteString(hex.EncodeToString(encodeUTF16(valStr, true)))
		b.WriteString(">\n")
		b.WriteString(">>\n")
	}

	// Footer
	b.WriteString("]\n")
	b.WriteString(">>\n")
	b.WriteString(">>\n")
	b.WriteString("endobj \n")
	b.WriteString("trailer\n")
	b.WriteString("\n")
	b.WriteString("<<\n")
	b.WriteString("/Root 1 0 R\n")
	b.WriteString(">>\n")
	b.WriteString("%%EOF\n")

	// Flush everything.
	return b.Flush()
}

// encodeUTF16 translates a utf8 string into a slice of bytes of ucs2.
// Taken from https://gist.github.com/ik5/65de721ca495fa1bf451
func encodeUTF16(s string, addBom bool) []byte {
	r := []rune(s)
	iresult := utf16.Encode(r)
	var bytes []byte
	if addBom {
		bytes = make([]byte, 2)
		bytes = []byte{254, 255}
	}
	for _, i := range iresult {
		temp := make([]byte, 2)
		binary.BigEndian.PutUint16(temp, i)
		bytes = append(bytes, temp...)
	}
	return bytes
}
