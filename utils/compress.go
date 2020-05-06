package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func ZipFile(filename string, zipName string) error {
	srcFile, err := os.Open(filename)
	if nil != err {
		return fmt.Errorf("failed to open source file, %w", err)
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if nil != err {
		return fmt.Errorf("failed to stat source file, %w", err)
	}


	destFile, err := os.Create(zipName)
	if nil != err {
		return fmt.Errorf("failed to create dest zip file, %w", err)
	}
	defer destFile.Close()

	zipWriter := zip.NewWriter(destFile)
	defer zipWriter.Close()

	header, err := zip.FileInfoHeader(info)
	if nil != err {
		return fmt.Errorf("failed to create zip header, %w", err)
	}
	header.Name = "/" + header.Name
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if nil != err {
		return fmt.Errorf("failed to create new zip header, %w", err)
	}

	n, err := io.Copy(writer, srcFile)
	fmt.Println(n)
	if nil != err {
		return fmt.Errorf("failed to copy data, %w", err)
	}

	return nil
}


