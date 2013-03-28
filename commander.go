package commander

import (
	"os"
	"fmt"
	"unicode/utf8"
	"path/filepath"
	"strings"
)


//Option type, contains data for a specific option
type Option struct {
	Name string
	Tiny string
	Verbose string
	Description string
	Required bool
	StringValue string
	Value interface{}
	Callback func(...string)
}


// Commander type, contains all program data 
type Commander struct {
	Name string
	Version string
	Options []*Option
	Opts map[string]*Option //only temporary, eventually replace Options
}


// Initialize a new commander with `args`
func Init(version string) *Commander {
	p := &Commander{
		Name: filepath.Base(os.Args[0]),
		Version: version,
		Opts: make(map[string]*Option),
	}

	p.Add(&Option{
		Name: "help",
		Tiny: "-h",
		Verbose: "--help",
		Description: "display usage",
		Required: false,
		StringValue: "",
		Callback: func(args ...string) {
			p.Usage()
		},
	})

	p.Add(&Option{
		Name: "version",
		Tiny: "-V",
		Verbose: "--version",
		Description: "display version",
		Required: false,
		StringValue: "",
		Callback: func(args ...string) {
			fmt.Fprintf(os.Stdout, "  Version: %s\n", p.Version)
			os.Exit(0);
		},
	})

	return p
}


//`Parse` arguments
func (commander *Commander) Parse() {
	args := commander.explode(os.Args[1:])
	//fmt.Println(args)

	for i, l := 0, len(commander.Options); i < l; i++ {
		option := commander.Options[i]

		found := false
		for j, l := 0, len(args); j < l; j++ {
			arg := args[j]

			if option.Tiny == arg || option.Verbose == arg {
				if j != l  - 1 && args[j + 1][0] != '-' {
					option.Callback(args[j + 1])
					j++
				} else {
					option.Callback()
				}
				found = true
			}
		}
		if option.Required && !found {
			// Option is required and not found
			fmt.Fprintf(os.Stderr, "%s, %s is required.", option.Tiny, option.Description)
		}
	}
}

func (commander *Commander) explode(args []string) []string {
	newargs := make([]string, 0, len(args))

	for i := range args {
		arg := args[i]
		l := utf8.RuneCountInString(arg)

		if l > 2 && arg[0] == '-' && arg[1] != '-' {
			for i := 1; i < l; i++ {
				newargs = append(newargs, "-" + string(arg[i]))
			}
		} else {
			newargs = append(newargs, arg)
		}
	}

	return newargs
}


//`Add` `option` to the commander instance
func (commander *Commander) Add(options ...*Option) {
	for i := range options {
		commander.Options = append(commander.Options, options[i])
		commander.Opts[options[i].Name] = options[i];
	}
}


//Display the usage of `commander`
func (commander *Commander) Usage() {
	fmt.Fprintf(os.Stdout, "\n  Usage: %s [options]\n\n", commander.Name)
	fmt.Fprintf(os.Stdout, "  Options:\n");

	options := commander.Options
	for i := range options {
		fmt.Fprintf(os.Stdout, "    %s, %s %s\n",
			options[i].Tiny, options[i].Verbose, options[i].Description)
	}
	fmt.Fprintf(os.Stdout, "\n")
	os.Exit(0)
}


//Return the total number of arguments
func (commander *Commander) Len() int {
	return len(os.Args)
}

//Add option in clean way.
func (commander *Commander) Option(switches string, description string) {
	commander.OptionWithDefault(switches, description, "")
}

func (commander *Commander) OptionWithDefault (switches string, description string, defaultVal string) {
	longArg := strings.Split(strings.TrimSpace(switches), " ")[0] //clear param if exists
	name := strings.TrimLeft(longArg, "--")
	tiny := ""

	if strings.Contains(switches, ",") {
		ss := strings.Split(switches, ",")
		tiny = ss[0]
		longArg = strings.Split(strings.TrimSpace(ss[1]), " ")[0] //clear param if exists
		name = strings.TrimLeft(longArg, "--")
	} 

	option := &Option{
		Name: name,
		Tiny: tiny,
		Verbose: longArg,
		Description: description,
		Required: false,
		StringValue: defaultVal,
		Value: nil,
		Callback: nil,
	}

	cb := func(args ...string) {
		if len(args) > 0 {
			option.StringValue = args[0]
		} else {
			option.StringValue = "true"
		}
	}

	option.Callback = cb
	
	commander.Add(option)
}



