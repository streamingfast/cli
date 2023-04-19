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

var noCast = ""
var viperUnsupported = ""

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
	{"BoolSlice", "[]bool", viperUnsupported, noCast},
	{"Uint8", "uint8", "Uint16", "uint8"},
	{"StringSlice", "[]string", "StringSlice", noCast},
	{"IPSlice", "[]net.IP", viperUnsupported, noCast},
	{"StringToString", "map[string]string", "StringMapString", noCast},
	{"Float64", "float64", "Float64", noCast},
	{"Uint32", "uint32", "Uint32", noCast},
	{"DurationSlice", "[]time.Duration", viperUnsupported, noCast},
	{"Uint16", "uint16", "Uint16", noCast},
	{"Float32Slice", "[]float32", viperUnsupported, noCast},
	{"Duration", "time.Duration", "Duration", noCast},
	{"Int64", "int64", "Int64", noCast},
	{"UintSlice", "[]uint", viperUnsupported, noCast},
	{"Bool", "bool", "Bool", noCast},
	{"Int32Slice", "[]int32", viperUnsupported, noCast},
	{"Int32", "int32", "Int32", noCast},
	{"StringToInt", "map[string]int", viperUnsupported, noCast},
	{"Int16", "int16", "Int32", "int16"},
	{"IP", "net.IP", viperUnsupported, noCast},
	{"IPNet", "net.IPNet", viperUnsupported, noCast},
	{"Uint64", "uint64", "Uint64", noCast},
	{"StringToInt64", "map[string]int64", viperUnsupported, noCast},
	{"Float32", "float32", "Float64", "float32"},
	{"IPv4Mask", "net.IPMask", viperUnsupported, noCast},
	{"Count", "int", viperUnsupported, noCast},
	{"Int", "int", "Int", noCast},
	{"Uint", "uint", "Uint", noCast},
	{"Float64Slice", "[]float64", viperUnsupported, noCast},
	{"Int8", "int8", "Int32", "int8"},
	{"BytesHex", "[]byte", viperUnsupported, noCast},
	{"BytesBase64", "[]byte", viperUnsupported, noCast},
	{"IntSlice", "[]int", "IntSlice", noCast},
	{"String", "string", "String", noCast},
	{"StringArray", "[]string", "StringSlice", noCast},
	{"Int64Slice", "[]int64", viperUnsupported, noCast},
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
	Name      string
	Type      string
	ViperName string
	ViperCast string
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"lower":      strings.ToLower,
		"pascalCase": strcase.ToCamel,
		"camelCase":  strcase.ToLowerCamel,
	}
}
