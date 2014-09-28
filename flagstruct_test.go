package flagstruct

import (
	"fmt"
	"log"
)

func Example() {
	flags := struct {
		//the rune: rest pattern redefines the split character from , to rune
		Name        string `flag:"|:name|Johnson, Rick|your name"`
		Simple_flag bool   `flag:""` //flag will be named simple-flag
		Sub         struct {
			On   bool   `flag:"on,,activate?"` //empty default value is zero
			Skip string //explicit tag required
		}
	}{}

	//this will be overwritten by flag parsing, if the flag is set
	flags.Name = "Angus"

	flags.Sub.Skip = "this is not set by flagstruct"

	//create a new flagstruct parser named example that fills flags.
	parser, err := New("example", &flags)
	if err != nil {
		log.Fatalln(err)
	}

	//this sets the passed in struct after normal flag parsing
	err = parser.Parse([]string{"-on", "-name", "John", "foo"})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Name:", flags.Name)
	fmt.Println("Simple_flag:", flags.Simple_flag)
	fmt.Println("On:", flags.Sub.On)
	fmt.Println("Skip:", flags.Sub.Skip)
	fmt.Println("Args:", parser.Args())
	// Output:
	// Name: John
	// Simple_flag: false
	// On: true
	// Skip: this is not set by flagstruct
	// Args: [foo]
}
