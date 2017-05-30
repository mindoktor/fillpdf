package fillpdf

/*
 *  FillPDF - Fill PDF forms, with extentions
 *  Copyright
 *  PKG Author: Roland Singer
 *  File Author: Alexander Félix
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

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Merge concatenates all input <files> and outputs one single pdf in <output>
func Merge(files ...string) (io.Reader, error) {
	args := []string{}

	// Get abs path for all input files while verifying their existence
	for _, f := range files {
		fAbsPath, err := getAbs(f)
		if err != nil {
			return nil, err
		}
		args = append(args, fAbsPath)
	}

	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "fillpdf-")
	if err != nil {
		return nil, err
	}

	// Remove the temporary directory on defer again.
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	// Create the temporary output file path.
	outputFile := filepath.Join(tmpDir, fmt.Sprintf("%d.pdf", time.Now().Unix()))

	// Create the pdftk command line arguments.
	args = append(args, "cat", "output", outputFile)

	// Run the pdftk utility.
	err = runCommandInPath(tmpDir, "pdftk", args...)
	if err != nil {
		return nil, fmt.Errorf("pdftk error: %v", err)
	}

	fb, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(fb), nil
}
