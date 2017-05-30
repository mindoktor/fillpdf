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

package fillpdf

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"unicode/utf16"
)

// Form represents the PDF form.
// This is a key value map.
type Form map[string]interface{}

// Fill a PDF form with the specified form values and create a final filled PDF file.
// One variadic boolean specifies, whenever to overwrite the destination file if it exists.
func Fill(form Form, formPDFFile, destPDFFile string, overwrite bool) error {
	var err error

	// Check if the pdftk utility exists.
	if _, err := exec.LookPath("pdftk"); err != nil {
		return err
	}

	// Get the absolute paths.
	if formPDFFile, err = getAbs(formPDFFile); err != nil {
		return err
	}

	if destPDFFile, err = filepath.Abs(destPDFFile); err != nil {
		return err
	}

	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "fillpdf-")
	if err != nil {
		return err
	}

	// Remove the temporary directory on defer again.
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	// Create the temporary output file path.
	outputFile := filepath.Clean(tmpDir + "/output.pdf")

	// Create the fdf data file.
	fdfFile := filepath.Clean(tmpDir + "/data.fdf")
	if err := createFdfFile(form, fdfFile); err != nil {
		return err
	}

	// Create the pdftk command line arguments.
	args := []string{
		formPDFFile,
		"fill_form", fdfFile,
		"output", outputFile,
		"flatten",
	}

	// Run the pdftk utility.
	if err := runCommandInPath(tmpDir, "pdftk", args...); err != nil {
		return fmt.Errorf("pdftk error: %v", err)
	}

	// Check if the destination file exists.
	e, err := exists(destPDFFile)
	if err != nil {
		return err
	} else if e {
		if !overwrite {
			return fmt.Errorf("destination PDF file already exists: '%s'", destPDFFile)
		}

		if err := os.Remove(destPDFFile); err != nil {
			return err
		}
	}

	// On success, copy the output file to the final destination.
	if err := copyFile(outputFile, destPDFFile); err != nil {
		return err
	}

	return nil
}

// createFdfFile with 16 bit encoded utf to enable creation of pdf with special characters
func createFdfFile(form Form, path string) error {
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
				valStr = "Yes"
			} else {
				valStr = "Off"
			}
		default:
			valStr = fmt.Sprintf("%v", value)
		}

		b.WriteString("<<\n")
		b.WriteString("/T (")
		b.Write(EncodeUTF16(key, true))
		b.WriteString(")\n")
		b.WriteString("/V (")
		b.Write(EncodeUTF16(valStr, true))
		b.WriteString(")\n")
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

// Taken from https://gist.github.com/ik5/65de721ca495fa1bf451
// EncodeUTF16 get a utf8 string and translate it into a slice of bytes of ucs2
func EncodeUTF16(s string, addBom bool) []byte {
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
