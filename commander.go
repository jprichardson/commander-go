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
	Hidden bool //won't show up in Usage() / '--help'
	Arg string
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
	TinyOpts map[string]*Option
	LongOpts map[string]*Option
}


// Initialize a new commander with `args`
func Init(version string) *Commander {
	p := &Commander{
		Name: filepath.Base(os.Args[0]),
		Version: version,
		Opts: make(map[string]*Option),
		TinyOpts: make(map[string]*Option),
		LongOpts: make(map[string]*Option),
	}

	p.Add(&Option{
		Name: "help",
		Tiny: "-h",
		Verbose: "--help",
		Hidden: false,
		Description: "display usage",
		Required: false,
		StringValue: "",
		Callback: func(args ...string) {
			p.Usage()
		},
	})

	p.Add(&Option{
		Name: "version",
		Tiny: "-v",
		Verbose: "--version",
		Hidden: false,
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
func (commander *Commander) oldParse() {
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
			} /*else {
				//fmt.Printf("arg: %s\n", arg)
			}*/
		}
		if option.Required && !found {
			// Option is required and not found
			fmt.Fprintf(os.Stderr, "%s, %s is required.", option.Tiny, option.Description)
		}
	}
}

func (prog *Commander) Parse() {
	args := os.Args[1:]//prog.explode(os.Args[1:])
	opts := make(map[string]*Option)

	//first pass
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.Index(arg, "--") == 0 { 
			if prog.LongOpts[arg] != nil {
				opts[arg] = prog.LongOpts[arg]
			}
		} else if strings.Index(arg, "-") == 0 {
			if prog.TinyOpts[arg] != nil {
				opts[arg] = prog.TinyOpts[arg]
			}
		}
	}

	//second pass
	for i:= 0; i < len(args); i++ {
		arg := args[i]

		if opt,ok := opts[arg]; ok {
			nextArg := ""
			if i < len(args) - 1 {
				nextArg = args[i + 1]
			}
			if opts[nextArg] == nil && nextArg != "" {
				opt.Callback(nextArg)
			} else {
				opt.Callback()
			}
		}
	}
}

func (commander *Commander) explode(args []string) []string {
	//fmt.Printf("%#v\n",args)
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

	//fmt.Printf("%#v\n",newargs)
	return newargs
}


//`Add` `option` to the commander instance
func (commander *Commander) Add(options ...*Option) {
	for i := range options {
		commander.Options = append(commander.Options, options[i])
		commander.Opts[options[i].Name] = options[i];
		commander.TinyOpts[options[i].Tiny] = options[i];
		commander.LongOpts[options[i].Verbose] = options[i];
	}
}


//Display the usage of `commander`
func (commander *Commander) Usage() {
	fmt.Fprintf(os.Stdout, "\n  Usage: %s [options]\n\n", commander.Name)
	fmt.Fprintf(os.Stdout, "  Options:\n");

	options := commander.Options

	longest := -1
	for i := range options {
		if !options[i].Hidden {
			s := fmt.Sprintf("    %s, %s %s", options[i].Tiny, options[i].Verbose, options[i].Arg)
			if len(s) > longest {
				longest = len(s)
			}
		}
	}

	longest += 3 //add a few more spaces for padding
	ofmt := fmt.Sprintf("%%-%ds %%s\n", longest)

	for i := range options {
		if !options[i].Hidden {
			left := fmt.Sprintf("    %s, %s %s", options[i].Tiny, options[i].Verbose, options[i].Arg)
			fmt.Fprintf(os.Stdout, ofmt, left, options[i].Description)
		}
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
	longStuff := strings.Split(strings.TrimSpace(switches), " ") //clear param if exists
	longName := longStuff[0]
	arg := ""
	if len(longStuff) > 1 {
		arg = longStuff[1]
	}
	name := strings.TrimLeft(longName, "--")
	tiny := ""

	if strings.Contains(switches, ",") { //"-c, --cheese <type>"
		ss := strings.Split(switches, ",")
		tiny = ss[0] //"-c"
		longStuff = strings.Split(strings.TrimSpace(ss[1]), " ") //"--cheese <type>"
		longName = longStuff[0] //"--cheese"
		if len(longStuff) > 1 {
			arg = longStuff[1] //"<type>"
		} else {
			arg = ""
		}

		name = strings.TrimLeft(longName, "--") //"cheese"
	} 

	option := &Option{
		Name: name,
		Tiny: tiny,
		Verbose: longName,
		Arg: arg,
		Hidden: false,
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



