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
		1. **Ol** inside?
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
	for i, line := range lines {
		if len(line) == 0 && i < len(lines)-1 {
			elements = append(elements, &Element{name: "nl", content: "\n"})
		}
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
			*elements = append(*elements, &Element{name: fmt.Sprintf("h%d", hashCount), content: strings.TrimLeft(line, "# ") + "\n"})
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

type IndexChar struct {
	i int
	c rune
}

func processVariableLine(elements *[]*Element, line string) bool {
	ruleString := "!*`[()]"
	specialChars := []*IndexChar{}
	for i, r := range line {
		if strings.Contains(ruleString, string(r)) {
			specialChars = append(specialChars, &IndexChar{i: i, c: r})
		}
	}

	cache := []byte{}
	for basePointer, lookAheadPointer := 0, 0; basePointer < len(line)-1; basePointer++ {
		char := string(line[basePointer])
		if strings.Contains(ruleString, char) {
			// look ahead for closer?
			if char == "*" {
				// is the next char a *? if so it's bold opener
				if line[basePointer+1] == '*' {

					// let's unshift 2 elements from specialChars
					specialChars = specialChars[2:]
					// it's bold
					// let's look ahead
					found := false

					// look for closer
					for _, indexChar := range specialChars {
						if indexChar.c == '*' {
							lookAheadPointer = indexChar.i
							found = true
							break
						}
					}

					if found && line[lookAheadPointer+1] == '*' {
						// now we have a range of text in the bold between the basePointer and lookAheadPointer that is bold
						boldText := line[basePointer+2 : lookAheadPointer]
						if len(cache) != 0 {
							*elements = append(*elements, &Element{name: "text", content: string(cache)})
							cache = []byte{}
						}
						// are there more chars?
						if lookAheadPointer+1 == len(line)-1 {
							// no
							*elements = append(*elements, &Element{name: "boldln", content: boldText + "\n"})
							break
						} else {
							// yes
							*elements = append(*elements, &Element{name: "bold", content: boldText})
						}
						// remove the further double * from specialChars

						newSpecialChars := []*IndexChar{}
						count := 0
						for i, specialChar := range specialChars {
							if count == 2 {
								newSpecialChars = append(newSpecialChars, specialChars[i:]...)
								break
							}
							if specialChar.c == '*' {
								count++
							} else {
								newSpecialChars = append(newSpecialChars, specialChar)
							}
						}
						specialChars = newSpecialChars
						basePointer = lookAheadPointer + 1
					} else {
						// we can ignore the * character, and the following one
						cache = append(cache, line[basePointer], line[basePointer+1])
						basePointer++
					}
					continue
				} else {
					// it's italic
					// let's unshift 2 elements from specialChars
					specialChars = specialChars[1:]
					// it's bold
					// let's look ahead
					found := false
					for _, indexChar := range specialChars {
						if indexChar.c == '*' {
							lookAheadPointer = indexChar.i
							found = true
							break
						}
					}
					if found {
						italicText := line[basePointer+1 : lookAheadPointer]
						if len(cache) != 0 {
							*elements = append(*elements, &Element{name: "text", content: string(cache)})
							cache = []byte{}
						}
						// are there more chars?
						if lookAheadPointer == len(line)-1 {
							// no
							*elements = append(*elements, &Element{name: "italicln", content: italicText + "\n"})
							break
						} else {
							// yes
							*elements = append(*elements, &Element{name: "italic", content: italicText})
						}
						newSpecialChars := []*IndexChar{}
						count := 0
						for i, specialChar := range specialChars {
							if count == 1 {
								newSpecialChars = append(newSpecialChars, specialChars[i:]...)
								break
							}
							if specialChar.c == '*' {
								count++
							} else {
								newSpecialChars = append(newSpecialChars, specialChar)
							}
						}
						specialChars = newSpecialChars
						basePointer = lookAheadPointer
					} else {
						// we can ignore the * character, and the following one
						cache = append(cache, line[basePointer])
					}
					continue
				}
			} else {
				// not bold and not italic
				if char == "[" {
					// link
					specialChars = specialChars[1:]
					// WARN: I will ignore internal bold and italics for now

					found := false
					tempLookAheadPointer := 0
					display := ""
					for _, indexChar := range specialChars {
						if indexChar.c == ']' {
							found = true
							tempLookAheadPointer = indexChar.i
							break
						}
					}
					if found && string(line[tempLookAheadPointer+1]) == "(" {
						// ok we have a square brackets enclosed text
						display = line[basePointer+1 : tempLookAheadPointer]

						// continue the same again but for parenthesis
						foundLink := false
						link := ""

						for _, indexChar := range specialChars {
							if indexChar.i <= tempLookAheadPointer+1 {
								continue
							}
							if indexChar.c == ')' {
								foundLink = true
								lookAheadPointer = indexChar.i
								break
							}
						}
						if foundLink {
							link = line[tempLookAheadPointer+2 : lookAheadPointer]
							// we definitely have a link
							if lookAheadPointer == len(line)-1 {
								// no
								*elements = append(*elements, &Element{name: "linkln", content: fmt.Sprintf("[%s](%s)\n", display, link)})
								break
							} else {
								// yes
								*elements = append(*elements, &Element{name: "link", content: fmt.Sprintf("[%s](%s) ", display, link)})
							}
							newSpecialChars := []*IndexChar{}
							count := 0
							foundClosingSquare := false
							for i, specialChar := range specialChars {
								if count == 1 {
									newSpecialChars = append(newSpecialChars, specialChars[i:]...)
									break
								}
								if specialChar.c == ']' {
									foundClosingSquare = true
								}
								if specialChar.c == ')' && foundClosingSquare {
									count++
								} else {
									newSpecialChars = append(newSpecialChars, specialChar)
								}
							}
							specialChars = newSpecialChars
							basePointer = lookAheadPointer
						} else {
							cache = append(cache, line[basePointer])
							continue
						}
					} else {
						cache = append(cache, line[basePointer])
						continue
					}

					// and then check for immediate re open parenthesis
					// and then check for close
				} else if char == "!" {
					// TODO: images
				}
			}
		}
		cache = append(cache, line[basePointer])
	}
	if len(cache) != 0 || len(line) == 1 {
		*elements = append(*elements, &Element{name: "textln", content: string(cache) + string(line[len(line)-1]) + "\n"})
	}

	return false
}
