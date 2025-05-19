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

	if !helper.IsDirectory(src) {
		fmt.Println("Source not a directory")
		os.Exit(1)
	}

	write := helper.IsDirectoryWritable(target)
	if !write {
		fmt.Println("Target not writeable")
		os.Exit(1)
	}

	filesSource, err := helper.CollectFileData(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filesTarget, err := helper.CollectFileData(target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diff := helper.CompareFileData(filesSource, filesTarget)
	if !helper.IsSyncRequired(*deleteFlag, diff) {
		fmt.Println("No sync required")
		os.Exit(0)
	}

	helper.ExplainSyncActions(diff)

	var proceed string
	fmt.Print("Do you want to proceed? (Y/n): ")
	fmt.Scanln(&proceed)

	if proceed != "Y" && proceed != "y" && proceed != "" {
		fmt.Println("Operation canceled.")
		os.Exit(0)
	}

	helper.SyncFiles(diff, src, target, *deleteFlag)

	fmt.Println("Done")
}
