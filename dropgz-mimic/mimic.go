package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	fmt.Println("Starting mimic")

	fmt.Println("Stat'ing file")
	const absolutePath = "C:\\Users\\paujohns\\sand\\repro-dropgz\\repro.exe"
	if _, err := os.Stat(absolutePath); err == nil {
		fmt.Println("File exists")
		fmt.Println("Renaming file")
		var newName = absolutePath + ".old"
		if err = os.Rename(absolutePath, newName); err != nil {
			fmt.Println(errors.Wrapf(err, "failed to rename the %s to %s", absolutePath, newName))
			return
		}
	}

	fmt.Println("Open file")
	_, err := os.OpenFile(absolutePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755) //nolint:gomnd // executable file bitmask
	if err != nil {
		fmt.Println(errors.Wrapf(err, "failed to create file %s", absolutePath))
		return
	}

	fmt.Println("io.Copy (write) to file")
}
