package gomd

import (
	"fmt"
	"strings"
)

var exampleText string = `# Header

**Great...** my links are all **Big Players**
[Amazon](https://amazon.com) [Google](https://google.com)
*This* is my *great list:*
- 1
- 2
  - a
  - b
		1. Ol inside?
    2. yes indeed
- [MSFT](https://microsoft.com)
- c

---
ooh noo
make sure you use the ` + "`builder.Generate()`" + ` function
1. first
2. second
3. third
  - really
  - great
    1. first again
    2. second again
      - so
      - good`

func ParseMD(md string, format MarkDownFormat) []*Element {
	lines := strings.Split(md, "\n")
	var elements []*Element
	target := &elements
	var parentStack []*Element

	nestCount := 0
	for _, line := range lines {
		// we allow for switching between the elements and children slices
		isListItem, generation, name, trimmed := identifyListedItem(line)
		if isListItem {
			line = trimmed
			if nestCount < generation {
				var targetParent *Element
				var rootParent *Element
				// create as many parents as required and link them in lineage order, and keep a pointer to the root parent
				for i := 0; i < generation-nestCount; i++ {
					parent := &Element{name: name, children: []*Element{}}
					if i == 0 {
						rootParent = parent
					} else {
						targetParent.children = append(targetParent.children, parent)
					}
					targetParent = parent

					parentStack = append(parentStack, parent)
				}
				// add the parent built with chidlren (if required) to the current target
				*target = append(*target, rootParent)

				// now we need to move the pointer to the most junior child
				target = &parentStack[len(parentStack)-1].children
				nestCount = generation
			} else if nestCount > generation {
				for i := 0; i < nestCount-generation; i++ {
					parentStack = parentStack[:len(parentStack)-1]
				}
				nestCount = generation

				if len(parentStack) > 0 {
					target = &parentStack[len(parentStack)-1].children
				} else {
					target = &elements
				}
			}
		} else {
			parentStack = []*Element{}
			nestCount = 0
			target = &elements
		}

		if processHeader(target, line) ||
			processHorizontalRule(target, line) ||
			processVariableLine(target, line) {
			continue
		}
	}

	return elements
}

func identifyListedItem(line string) (bool, int, string, string) {
	trimmed := strings.TrimLeft(line, " \t")
	name := ""
	if len(trimmed) >= 2 && trimmed[0] == '-' && trimmed[1] == ' ' {
		name = "ul"
	}
	if len(trimmed) >= 3 &&
		strings.Contains("123456789", string(trimmed[0])) &&
		trimmed[1] == '.' &&
		trimmed[2] == ' ' {
		name = "ol"
	}
	if len(name) > 0 {
		spaceCount := 0
		for _, r := range line {
			if string(r) == "\t" {
				spaceCount += 2
			} else if r == ' ' {
				spaceCount += 1
			} else {
				break
			}
		}

		if name == "ul" {
			trimmed = trimmed[2:]
		} else {
			trimmed = trimmed[3:]
		}

		return true, (spaceCount / 2) + 1, name, trimmed
	}

	return false, 0, "", ""
}

// processHeader determines if the line has a valid header by counting hashes and checking for a space that follows immediately.
// it returns true and appends the Element pointer to the dereferenced *[]*Element slice if it identifies a valid header.
// it returns false and does nothing if no valid header is found.
// it identifies the header type based on the hash count
func processHeader(elements *[]*Element, line string) bool {
	if strings.HasPrefix(line, "#") {
		hashCount := 0
		isHeader := true

		for i, r := range line {
			if r == '#' {
				hashCount = i + 1
			} else if hashCount > 0 {
				if r != ' ' {
					isHeader = false
				}
				break
			} else {
				isHeader = false
			}
		}

		if isHeader {
			*elements = append(*elements, &Element{name: fmt.Sprintf("h%d", hashCount), content: strings.TrimLeft(line, "# ")})
			return isHeader
		}
	}
	return false
}

func processHorizontalRule(elements *[]*Element, line string) bool {
	if strings.HasPrefix(line, "---") {
		invalidChar := strings.ContainsFunc(line, func(r rune) bool {
			if r != ' ' && r != '-' {
				return true
			}
			return false
		})

		if !invalidChar {
			*elements = append(*elements, &Element{name: "rule", content: line})
			return true
		}
	}
	return false
}

func processVariableLine(elements *[]*Element, line string) bool {
	// check if it's a bland line?
	// if !strings.ContainsAny(line, "*`[()]") {
	// it's a clean text?
	// no OL and UL check
	*elements = append(*elements, &Element{name: "line", content: line})
	//}
	return false
}
