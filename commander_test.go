package commander

import (
	"os"
	"testing"
	"fmt"
	"github.com/jprichardson/goatee-go"
)

func TestStart(*testing.T) {
	fmt.Println("Tests starting... ") //Go will bitch if we don't use it somewhere
}

func TestInit(*testing.T) {
	prog := Init("0.0.3")
	t.EQ (prog.Version, "0.0.3") 
	t.NEQ (len(prog.Name), 0) 
}

func TestOption(*testing.T) {
	prog := Init("0.0.1")
	prog.Option("-c, --cheese <type>", "Add the specified type of cheese.")
	os.Args = []string{"bogus", "-c", "mozzarella"}
	prog.Parse()

	t.EQ (prog.Opts["cheese"].Name, "cheese")
	t.EQ (prog.Opts["cheese"].StringValue, "mozzarella")

	os.Args = []string{"bogus", "--cheese", "cheddar"}
	prog.Parse()
	t.EQ (prog.Opts["cheese"].StringValue, "cheddar")
}

func TestOptionWithDefault (*testing.T) {
	prog := Init("")
	prog.OptionWithDefault("-l, --lines <number>", "number of lines, default 10", "10")
	t.EQ (prog.Opts["lines"].StringValue, "10")
}

func TestOptionWithoutTiny (*testing.T) {
	prog := Init("")
	prog.Option("--chunks <number>", "number of chunks to parse")

	os.Args = []string{"progname", "--chunks", "500"}
	prog.Parse()
	
	t.EQ (prog.Opts["chunks"].StringValue, "500")
	t.EQ (prog.Opts["chunks"].Tiny, "")
}

func TestArgData(*testing.T) {
	exec := false

	doWork := func(s ...string) {
		if len(s) < 1 || len(s) > 1 {
			panic("Got an invalid number of arguments")
		} else {
			exec = true
		}
	}

	prog := Init("0.0.1")
	prog.Add(&Option{
		Name:        "work",
		Tiny:        "-w",
		Verbose:     "--work",
		Description: "do work",
		Required:    true,
		Callback:    doWork,
	})

	// Mock arguments to parse
	os.Args = []string{"bogus", "-w", "work"}
	prog.Parse()

	t.T (exec)
}

func TestArgs(*testing.T) {
	exec := 0

	doWork := func(s ...string) {
		exec++
	}

	prog := Init("0.0.1")
	prog.Add(&Option{
		Name:        "work",
		Tiny:        "-w",
		Verbose:     "--work",
		Description: "do work",
		Required:    true,
		Callback:    doWork,
	})

	prog.Add(&Option{
		Name:        "more work",
		Tiny:        "-m",
		Verbose:     "--more",
		Description: "do work",
		Required:    true,
		Callback:    doWork,
	})

	// Mock arguments to parse
	os.Args = []string{"bogus", "-w", "-m"}
	prog.Parse()

	t.EQ(exec, 2) 
}
