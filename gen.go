package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var abbrevs = []struct {
	kind        string
	what        string
	replacement string
	description string
	charsback   int
}{

	{"go", "aa", "append(,)", "", -2},
	{"go", "ife", "if err != nil {\n\t\n\t}\n", "", -4},
	{"go", "ie", "if err := ", "", 0},
	{"go", ";e", "err != nil {\t\n\t\t\n\t}", "", -3},
	{"go", "re", "return err", "", 0},
	{"go", "rf", "return fmt.Errorf(\": %w\", err)", "", -11},
	{"go", "rn", "return nil", "", 0},
	{"go", "ew", "errors.Wrap(err, \"\")", "", -2},
	{"go", "pp", "println(\"=== \")", "", -2},
	{"go", "pps", "println(\"=== \", fmt.Sprintf(\"%\",))", "", -4},
	{"go", "fu", "func () {\n}\n", "", -7},

	{"testing", "tt", "func Test(t*testing.T)", "", -13},
	{"testing", "tl", "t.Log(\"\",)", "", -2},
	{"testing", "tlf", "t.Logf(\"\",)", "", -3},
	{"testing", "tf", "t.Fatal(err)", "", 0},

	{"fmt", "ff", "fmt.Printf(\"\",)", "", -3},
	{"fmt", "fff", "fmt.Fprintf(w, \"\",)", "", -3},
	{"fmt", "ffp", "fmt.Fprintln(w, )", "", -1},
	{"fmt", "sp", "fmt.Sprintf(\"\",)", "", -3},
	{"fmt", "fe", "fmt.Errorf(\"\",)", "", -3},
	{"fmt", "fp", "fmt.Println()", "", -1},

	{"log", "lp", "log.Println()", "", -1},
	{"log", "lpf", "log.Printf(\"\",)", "", -3},
	{"log", "lf", "log.Fatal(err)", "", -1},

	{"strings", "ss.", "strings.", "", 0},

	{"io", "iora", "io.ReadAll()", "", -1},

	{"os", "osrf", "os.ReadFile(fname)", "", -1},
	{"os", "oswf", "os.WriteFile(fname,,0600)", "", -6},
	{"os", "osc", "os.Create()", "", -1},
	{"os", "oso", "os.Open()", "", -1},

	{"http", "he", "http.Error(w, err.Error(), http.StatusInternalServerError)\nreturn\n", "", 0},
	{"http", "hw", "w http.ResponseWriter", "", 0},
	{"http", "hr", "r*http.Request", "", 0},
	{"http", "hg", "", "http.Get() and error handling", 0},

	{"sql", "qsel", "SELECT *\nFROM\nWHERE\nGROUP BY\nORDER BY\nLIMIT 1000\n;\n", "select *", -1},

	{"zap", "zs", "zap.String(\"\", )", "", -4},
	{"zap", "zi", "zap.Int64(\"\", )", "", -4},
	{"zap", "ze", "zap.Error(err)", "", 0},
}

func getSnippet(abbrev string) (string, error) {
	buf, err := ioutil.ReadFile("gen.go")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(buf), "\n")
	var text string
	var dump bool

	for _, ln := range lines {
		if strings.HasPrefix(ln, "func Snippet_"+abbrev+"()") {
			dump = true
			continue
		}
		if !dump {
			continue
		}
		if ln == "}" {
			break
		}
		text += ln + "\n"
	}

	return text, nil
}

func Snippet_hg() error {
	resp, err := http.Get("url")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %s", resp.Status)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = buf
	return nil
}

func main() {
	f, err := os.Create("abbrevs.lua")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fh, err := os.Create("help/goabbrevs.md")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	fmt.Fprintln(f, "-- This file is generated by gen.go. Do not edit.")
	fmt.Fprintln(f, "-- Run `go run gen.go` to regenerate.")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "function expand(str)")

	for _, a := range abbrevs {
		if a.replacement == "" {
			a.replacement, err = getSnippet(a.what)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Fprintf(f, "\tif str == %q then\n", a.what)
		fmt.Fprintf(f, "\t\treturn %q, %d\n", a.replacement, a.charsback)
		fmt.Fprintln(f, "\tend")
	}

	fmt.Fprintf(f, "\treturn %q, 0\n", "")
	fmt.Fprintln(f, "end")

	fmt.Fprint(fh, "# goabbrevs\n========\n\n")
	fmt.Fprintln(fh, "# kind    abbrev     replacement")
	repl := strings.NewReplacer("\t", "", "\n", " ")
	for _, a := range abbrevs {
		if a.description != "" {
			a.replacement = a.description
		}
		fmt.Fprintf(fh, "%-10s %-10s %s\n", a.kind, a.what, repl.Replace(a.replacement))
	}
	fmt.Fprintln(fh)
	fmt.Fprintln(fh, "# generated by\n", os.Getenv("HOME")+"/.config/micro/plug/goabbrevs/gen.go")
}
