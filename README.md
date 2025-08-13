# gomd

## Motivation

Welcome to race-conditioned/gomd, a markdown library written in go. The motivation for this project is to be able to programatically create markdown files that will serve like github tickets but with less overhead for personal projects.

I realised that in order to do this I would be best served with a library that can build and parse markdown. However I quickly found that markdown is a very loose grammar markup language and has many edge cases. I want to be able to round trip some markdown so that I can save a file by my specification, then load it again, edit it, and save it again.

This project is a WIP, and in these early stages, breaking changes may be released, so please contact me if you are interested in using it or contributing.

## Usage

There are two main ways to build markdown from this library. One way is to use the Compounder, which is good for simple creation of markdown. Another way is to use the Builder for more fine grained control.

This README was written using the `gomd.Compounder`

You can see the usage in the `example` directory.

Here is An example of using the Builder for raw constructions:

```go
 func main() {
	brandName := "X Company"
	b := gomd.Builder{}
	header := []*gomd.Element{
		b.H1(fmt.Sprintf("My %s Document", brandName)),
		b.NL(),
		b.Textln("great!"),
		b.NL(),
		b.UL(
			b.Textln("first"),
			b.Textln("second"),
			b.OL(
				b.Bold("first"),
				b.Textln(" element"),
			),
		),
	}

	body := []*gomd.Element{
		b.Text("This is the body"),
	}

	template := []*gomd.Element{}
	template = append(template, header...)
	template = append(template, b.NL())
	template = append(template, body...)

	md := b.Generate(template...)
	err := Write("xcompany.md", md)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
```

As you can see templates can be composed and mixed and matched at your discretion, or you can input builder functions directly into the `Build` function

## Features

- markdown builder ✅
- markdown compounder ✅
- Read and Write markdown ✅
- Basic Markdown syntax supported ✅
- Builder and composer tested ✅
- Common Mark compatible - planned but a long way off
- Deep nesting - partial support
- Round trip support - planned
- Serve to a viewer - planned
- Conversion to HTML - planned
