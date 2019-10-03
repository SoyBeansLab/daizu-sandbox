package main

import (
	"flag"
	"io/ioutil"
	"os"

	diff "github.com/sergi/go-diff/diffmatchpatch"
)

func checkErr(err error) {
	if err != nil {
		os.Exit(2)
	}
}

func main() {
	flag.Parse()

	f1, err := os.Open(flag.Arg(0))
	defer f1.Close()
	checkErr(err)
	text1, err := ioutil.ReadAll(f1)
	checkErr(err)

	f2, err := os.Open(flag.Arg(1))
	defer f2.Close()
	checkErr(err)
	text2, err := ioutil.ReadAll(f2)
	checkErr(err)

	dmp := diff.New()
	diffs := dmp.DiffMain(string(text1), string(text2), false)

	// 一緒の場合は1が帰ってくる
	if len(diffs) != 1 {
		os.Exit(1)
	}
}
