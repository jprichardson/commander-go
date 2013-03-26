package commander

import (
	"os"
	"testing"
	//"fmt"
)

//fuck Google and their encouragement of verbosity
func T (expr bool) {
	if !expr {
		panic("Not True")
	}
}

func F (expr bool) {
	if expr {
		panic("True")
	}
}

func TestInit(t *testing.T) {
	prog := Init("0.0.3")
	if prog.Version != "0.0.3" {
		t.Error("Init Version failure.")
	}

	if len(prog.Name) == 0 {
		t.Error("Init Name failure.")
	}
}

func TestOption(t *testing.T) {
	prog := Init("0.0.1")
	prog.Option("-c, --cheese <type>", "Add the specified type of cheese.")
	os.Args = []string{"bogus", "-c", "mozzarella"}
	prog.Parse()

	T (prog.Opts["cheese"].Name == "cheese")
	T (prog.Opts["cheese"].StringValue == "mozzarella")

	os.Args = []string{"bogus", "--cheese", "cheddar"}
	prog.Parse()
	//fmt.Println(prog.Opts["cheese"])
	T (prog.Opts["cheese"].StringValue == "cheddar")
}

func TestArgData(t *testing.T) {
	exec := false

	doWork := func(s ...string) {
		if len(s) < 1 || len(s) > 1 {
			t.Error("Got an invalid number of arguments")
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

	if !exec {
		t.Error("Did not execute `doWork()`")
	}
}

func TestArgs(t *testing.T) {
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

	if exec != 2 {
		t.Error("Did not run all callbacks")
	}
}
