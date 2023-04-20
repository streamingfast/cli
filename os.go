package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(inPath, outPath string) {
	inFile, err := os.Open(inPath)
	NoError(err, "Unable to open actual file %q", inPath)
	defer inFile.Close()

	outFile, err := os.Create(outPath)
	NoError(err, "Unable to open expected file %q", outPath)
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	NoError(err, "Unable to copy file %q to %q", inPath, outPath)
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		// For this script, we don't care
		return false
	}

	return !stat.IsDir()
}

func DirectoryExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		// For this script, we don't care
		return false
	}

	return stat.IsDir()
}

// WriteFile is a quick version `os.WriteFile` where [NoError] is used to
// ensure no error occur.
func WriteFile(name string, content string, args ...any) {
	NoError(os.WriteFile(name, []byte(fmt.Sprintf(content, args...)), os.ModePerm), "Unable to write file")
}

func ReadFile(name string) string {
	content, err := os.ReadFile(name)
	NoError(err, "Unable to read file %q", name)

	return string(content)
}

func WorkingDirectory() string {
	directory, err := os.Getwd()
	NoError(err, "Unable to get working directory")

	return directory
}

func UserHomeDirectory() string {
	home, err := os.UserHomeDir()
	NoError(err, "Unable to get user home directory")

	return home
}

func AbsolutePath(in string) string {
	out, err := filepath.Abs(in)
	NoError(err, "Unable to make path %q absolute", in)

	return out
}
