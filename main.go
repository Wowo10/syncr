package main

import (
	"flag"
	"fmt"
	"os"
	"syncr/helper"
)

func main() {
	deleteFlag := flag.Bool("delete-missing", false, "Delete Missing files in target directory")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <src> <target>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	src := args[0]
	target := args[1]

	fmt.Println("Source:", src, "->", helper.IsDirectory(src))
	fmt.Println("Target:", target, "->", helper.IsDirectory(target))
	fmt.Println("Skip:", *deleteFlag)
}
