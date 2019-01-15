package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/lindell/mockay/mockgen"
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [path]\n", os.Args[0])

	flag.PrintDefaults()
}

var position = regexp.MustCompile("^(\\d+):(\\d+)$")

func main() {
	verbose := flag.Bool("verbose", false, "print logging statements")
	outputFile := flag.String("o", "", "output file (otherwise stdout is used)")
	pos := flag.String("pos", "", "the position of the interface to be mocked, as line:column")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		return
	}
	path := args[0]

	options := []mockgen.Option{}
	if *verbose {
		options = append(options, mockgen.WithLogger(&logger{}))
	}
	if *outputFile != "" {
		file, err := os.OpenFile(*outputFile, os.O_WRONLY|os.O_CREATE, 0660)
		defer file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open file: %s\n", *outputFile)
		}
		options = append(options, mockgen.WithWriter(file))
	}
	posMatch := position.FindStringSubmatch(*pos)
	if posMatch != nil {
		x, _ := strconv.Atoi(posMatch[1])
		y, _ := strconv.Atoi(posMatch[2])
		options = append(options, mockgen.WithPosition(mockgen.Position{X: x, Y: y}))
	}
	generator := mockgen.New(options...)

	err := generator.Generate(path)
	if err != nil {
		fmt.Println(err)
	}
}

type logger struct{}

func (l *logger) Info(str string) {
	fmt.Println(str)
}
