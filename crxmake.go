package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcuadros/go-crxmake/crxmake"

	"github.com/jessevdk/go-flags"
)

var (
	ErrMissingFolder  = errors.New("missing folder extension")
	ErrFolderNotFound = errors.New("unable to find folder extension")

	outputFormat = "%s.crx"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)
	cmd, err := parser.AddCommand("crxmake", "", "", &Command{})
	//it replace the defualt command
	parser.Command = cmd

	_, err = parser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}

			fmt.Println(err)
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}

type Command struct {
	Options struct {
		Folder string `positional-arg-name:"folder" description:"folder where the extension is located."`
		Output string `positional-arg-name:"output" description:"output file name"`
	} `positional-args:"yes"`
}

func (c *Command) Execute(args []string) error {
	if err := c.init(); err != nil {
		return err
	}

	b := crxmake.NewBuilder()

	fmt.Printf("Reading files from %q\n", c.Options.Folder)
	err := b.BuildZip(c.Options.Folder)
	if err != nil {
		return err
	}

	fmt.Printf("Writing file %q ... ", c.Options.Output)
	file, err := os.Create(c.Options.Output)
	if err != nil {
		fmt.Println("FAIL")
		return err
	}

	defer file.Close()

	err = b.WriteToFile(file)
	if err != nil {
		fmt.Println("FAIL")
		return err
	}

	fmt.Println("DONE")
	return nil
}

func (c *Command) init() error {
	if c.Options.Folder == "" {
		return ErrMissingFolder
	}

	if _, err := os.Stat(c.Options.Folder); err != nil {
		return ErrFolderNotFound
	}

	if c.Options.Output == "" {
		c.Options.Output = fmt.Sprintf(outputFormat, filepath.Base(c.Options.Folder))
	}

	return nil
}
