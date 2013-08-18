golang: commander 
=================

Note: Fork of [commander.go](https://github.com/vanetix/commander.go)

<!-- 

[![Build Status](https://travis-ci.org/jprichardson/go-commander.png?branch=master)](https://travis-ci.org/jprichardson/go-commander)

-->

Unix like argument parsing for go. Heavily inspired by: [commander.c](https://github.com/visionmedia/commander.c) and [commander.js](https://github.com/visionmedia/commander.js)


Installing
----------

    go get github.com/jprichardson/commander-go


Usage
-----

```go
func doWork(args ...string) {
	// Do some work
}

func main() {
	program := commander.Init("0.1.3")

	program.Add(&commander.Option{
	    Required: false,
	    Name: "Work",
	    Tiny: "-w",
	    Verbose: "--work",
	    Description: "do some work",
	    Callback: doWork,
	})

	program.Parse()
}
```

or...

```go
func main() {
  program := commander.Init("0.1.3")

  program.Option("-p, --peppers", "Add peppers")
  program.Option("-P, --pineapple", "Add pineapple")
  program.Option("-b, --bbq", "Add bbq sauce")
  program.Option("-c, --cheese [type]", "Add the specified type of cheese")
  program.Parse()

  //access
  Program.Opts["cheese"].Value; //contains value passed like: "-c cheddar" or "--cheese cheddar" 
}
```



License (MIT)
-------------
Copyright (c) 2013 Matt McFarland, JP Richardson


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
