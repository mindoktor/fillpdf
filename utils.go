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
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mindoktor/mderrors"
)

func getAbs(path string) (string, error) {
	// Get the absolute paths.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", mderrors.Wrap(err)
	}

	// Check if the form file exists.
	e, err := exists(path)
	if err != nil {
		return "", mderrors.Wrap(err)
	} else if !e {
		return "", mderrors.New("file does not exists", path)
	}

	return absPath, nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, mderrors.Wrap(err)
}

// copyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return mderrors.Wrap(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return mderrors.Wrap(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return mderrors.Wrap(err)
	}
	return mderrors.Wrap(out.Sync())
}

// runCommandInPath runs a command and waits for it to exit.
// The working directory is also set.
// The stderr error message is returned on error.
func runCommandInPath(dir, name string, args ...string) error {
	// Create the command.
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Dir = dir

	// Start the command and wait for it to exit.
	err := cmd.Run()
	if err != nil {
		return mderrors.Wrap(err, "stderr", strings.TrimSpace(stderr.String()))
	}

	return nil
}
