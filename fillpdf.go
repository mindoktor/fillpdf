/*
 *  FillPDF - Fill PDF forms
 *  Copyright DesertBit
 *  Author: Roland Singer
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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	iconv "gopkg.in/iconv.v1"
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
	// defer func() {
	// 	os.RemoveAll(tmpDir)
	// }()

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
	b.Write([]byte("%FDF-1.2\n"))
	b.Write([]byte("\xE2\xE3\xCF\xD3\n"))
	b.Write([]byte("1 0 obj \n"))
	b.Write([]byte("<<\n"))
	b.Write([]byte("/FDF \n"))
	b.Write([]byte("<<\n"))
	b.Write([]byte("/Fields [\n"))

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

		b.Write([]byte("<<\n"))
		b.Write([]byte("/T (" + toUTF16(key) + ")\n"))
		b.Write([]byte("/V (" + toUTF16(valStr) + ")\n"))
		b.Write([]byte(">>\n"))
	}

	// Footer
	b.Write([]byte("]\n"))
	b.Write([]byte(">>\n"))
	b.Write([]byte(">>\n"))
	b.Write([]byte("endobj \n"))
	b.Write([]byte("trailer\n"))
	b.Write([]byte("\n"))
	b.Write([]byte("<<\n"))
	b.Write([]byte("/Root 1 0 R\n"))
	b.Write([]byte(">>\n"))
	b.Write([]byte("%%EOF\n"))

	// Flush everything.
	return b.Flush()
}

// Each field is interptreted as one UTF-16 entity and thus needs to have the BOM bytes on each seperate string.
// Achieve this by opening a new handle to all convertions
func toUTF16(input string) string {
	cd, err := iconv.Open("utf-16", "utf-8")
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return ""
	}
	defer cd.Close()
	return cd.ConvString(input)
}
