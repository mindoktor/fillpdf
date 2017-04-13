package fillpdf

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mindoktor/mderrors"
)

// Multistamp stamps one PDF ontop of another, returns a reader to bytes generated.
func Multistamp(stampontoPDFFile, stampPDFFile string) (io.Reader, error) {
	var err error

	// Check if the pdftk utility exists.
	if _, err := exec.LookPath("pdftk"); err != nil {
		return nil, mderrors.Wrap(err)
	}

	if stampontoPDFFile, err = getAbs(stampontoPDFFile); err != nil {
		return nil, mderrors.Wrap(err)
	}

	stampPDFFile, err = getAbs(stampPDFFile)
	if err != nil {
		return nil, mderrors.Wrap(err)
	}

	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "fillpdf-")
	if err != nil {
		return nil, mderrors.Wrap(err)
	}

	// Remove the temporary directory on defer again.
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	// Create the temporary output file path.
	outputFile := filepath.Clean(tmpDir + "/output.pdf")

	// Create the pdftk command line arguments.
	args := []string{
		stampontoPDFFile,
		"multistamp", stampPDFFile,
		"output", outputFile,
	}

	// Run the pdftk utility.
	err = runCommandInPath(tmpDir, "pdftk", args...)
	if err != nil {
		return nil, mderrors.Wrap(err)
	}

	fb, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return nil, mderrors.Wrap(err)
	}

	return bytes.NewReader(fb), nil
}
