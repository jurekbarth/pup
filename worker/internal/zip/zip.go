package zip

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Zip the file
func Zip(dir string, dest string) error {
	baseFolder := dir
	if filepath.Ext(baseFolder) == "" {
		baseFolder = filepath.Clean(baseFolder) + "/"
	}

	// Get a Buffer to Write To
	outFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	err = addFiles(w, baseFolder, "")

	if err != nil {
		return err
	}
	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func addFiles(w *zip.Writer, basePath string, baseInZip string) error {

	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			err = addFile(w, basePath, baseInZip, file.Name())
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			bs := baseInZip + file.Name() + "/"
			addFiles(w, newBase, bs)
		}
	}
	return nil
}

func addFile(w *zip.Writer, basePath string, baseInZip string, fileName string) error {
	data, err := ioutil.ReadFile(basePath + fileName)
	if err != nil {
		return err
	}
	return addData(w, basePath, baseInZip, fileName, data)
}

func addData(w *zip.Writer, basePath string, baseInZip string, fileName string, data []byte) error {
	// Add some files to the archive.
	f, err := w.Create(baseInZip + fileName)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
