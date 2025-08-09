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
	isElements := true
	var parentStack []*Element

	nestCount := 0
	for _, line := range lines {
		fmt.Println("len(ps)", len(parentStack))
		// we allow for switching between the elements and children slices

		isListItem, generation, name := identifyListedItem(line)
		if isListItem {
			if nestCount < generation {
				// we just nested so create a parent
				fmt.Println("was elements:", isElements)
				parent := &Element{name: name, children: []*Element{}}
				parentStack = append(parentStack, parent)
				*target = append(*target, parent)
				target = &parent.children
				isElements = false
				nestCount = generation
			} else if nestCount > generation {
				fmt.Println("len of parentstack:", len(parentStack))
				fmt.Println("jumps back:", nestCount-generation)
				for i := 0; i < nestCount-generation; i++ {
					// pop a parent from the stack
					parentStack = parentStack[:len(parentStack)-1]
				}
				nestCount = generation

				if len(parentStack) > 0 {
					target = &parentStack[len(parentStack)-1].children
				} else {
					target = &elements
				}

				isElements = true
			}
			// fmt.Println("found Item!", isListItem, generation, line)
		} else {
			// fmt.Println("not list Item")
			// wipe the parent stack?
			parentStack = []*Element{}
			nestCount = 0
			target = &elements
			isElements = true
		}

		fmt.Println("isElements:", isElements, "|", line)

		if processHeader(target, line) ||
			processHorizontalRule(target, line) ||
			processVariableLine(target, line) {
			continue
		}
	}

	return elements
}

func identifyListedItem(line string) (bool, int, string) {
	// check for the first character being a number or hyphen
	trimmed := strings.TrimLeft(line, " ")
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
		// we have a list item
		// so we need to count the left hand spaces and divide by 2
		spaceCount := 0
		for i, r := range line {
			if r != ' ' {
				spaceCount = i
				break
			}
		}

		return true, (spaceCount / 2) + 1, name

	}

	// not a list item
	return false, 0, ""
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
