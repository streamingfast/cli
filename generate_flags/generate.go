package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/streamingfast/cli"
)

//go:embed *.gotmpl
var templates embed.FS

// To generate the definition list, you can use the following command. But first,
// find the path to pflag package or download locally a version. The rest of the
// command assume enviornment variable 'PFLAG_PATH' points to there.
//
// I checked in my editor where 'cmd.Flags().GetString(<name>)' was defined and used
// path.
//
// find "$PFLAG_PATH" -name "*.go" | grep -v "_test.go" | xargs -n1 -I{} grep -E 'func \(f \*FlagSet\) Get' {}
//
// This will list all 'GetXXX' function, then perform multi-cursor selection in your editor
// and extract the part after 'Get' and the actual type.
var definitions = []Definition{
	{"BoolSlice", "[]bool"},
	{"Uint8", "uint8"},
	{"StringSlice", "[]string"},
	{"IPSlice", "[]net.IP"},
	{"StringToString", "map[string]string"},
	{"Float64", "float64"},
	{"Uint32", "uint32"},
	{"DurationSlice", "[]time.Duration"},
	{"Uint16", "uint16"},
	{"Float32Slice", "[]float32"},
	{"Duration", "time.Duration"},
	{"Int64", "int64"},
	{"UintSlice", "[]uint"},
	{"Bool", "bool"},
	{"Int32Slice", "[]int32"},
	{"Int32", "int32"},
	{"StringToInt", "map[string]int"},
	{"Int16", "int16"},
	{"IP", "net.IP"},
	{"IPNet", "net.IPNet"},
	{"Uint64", "uint64"},
	{"StringToInt64", "map[string]int64"},
	{"Float32", "float32"},
	{"IPv4Mask", "net.IPMask"},
	{"Count", "int"},
	{"Int", "int"},
	{"Uint", "uint"},
	{"Float64Slice", "[]float64"},
	{"Int8", "int8"},
	{"BytesHex", "[]byte"},
	{"BytesBase64", "[]byte"},
	{"IntSlice", "[]int"},
	{"String", "string"},
	{"StringArray", "[]string"},
	{"Int64Slice", "[]int64"},
}

func main() {
	cli.Ensure(len(os.Args) == 3, "go run ./flags <output_file> <package_name>")

	output := os.Args[1]
	packageName := os.Args[2]

	tmpl, err := template.New("flags").Funcs(templateFunctions()).ParseFS(templates, "*.gotmpl")
	cli.NoError(err, "Unable to instantiate template")

	var out io.Writer = os.Stdout
	if output != "-" {
		cli.NoError(os.MkdirAll(filepath.Dir(output), os.ModePerm), "Unable to create output file directories")

		file, err := os.Create(output)
		cli.NoError(err, "Unable to open output file")

		bufferedOut := bufio.NewWriter(file)
		out = bufferedOut

		defer func() {
			bufferedOut.Flush()
			file.Close()
		}()
	}

	err = tmpl.ExecuteTemplate(out, "template.gotmpl", map[string]any{
		"Package":     packageName,
		"Definitions": definitions,
	})
	cli.NoError(err, "Unable to render template")

	fmt.Println("Done")
}

type Definition struct {
	Name string
	Type string
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"lower":      strings.ToLower,
		"pascalCase": strcase.ToCamel,
		"camelCase":  strcase.ToLowerCamel,
	}
}
