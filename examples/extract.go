package main

import "github.com/kelvyne/swf"
import "os"
import "fmt"
import "io/ioutil"

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("os.open: %v\n", err)
	}
	defer file.Close()

	parser := swf.NewParser(file)
	s, err := parser.Parse()
	if err != nil {
		fmt.Printf("swf.Parse: %v\n", err)
	}

	for _, tag := range s.Tags {
		if tag.Code() == swf.CodeTagDoABC {
			doAbc := tag.(*swf.TagDoABC)
			filename := fmt.Sprintf("./%v.abc", doAbc.Name)
			err := ioutil.WriteFile(filename, doAbc.ABCData, 0644)
			if err != nil {
				fmt.Printf("ioutil.WriteFile: %v", err)
			}
		}
	}
}
