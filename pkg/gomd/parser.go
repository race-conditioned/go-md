package gomd

import (
	"strings"
)

func ParseMD(md string, format MarkDownFormat) []*Element {
	lines := strings.Split(md, "\n")
	var elements *[]*Element = &[]*Element{}
	target := elements
	var parentStack *[]*Element = &[]*Element{}

	nestCount := 0
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 && i < len(lines)-1 {
			// is the next line a rule?
			if i <= len(lines)-1 && len(lines) > 0 {
				if lines[i+1] == "---" {
					continue
				}
			}

			*elements = append(*elements, &Element{Kind: KNewLine, LineBreak: true})
		}
		// we allow for switching between the elements and Children slices
		isListItem, generation, name, trimmed := identifyListedItem(line)
		if isListItem {
			line = trimmed
			handleListItem(name, &nestCount, &generation, parentStack, elements, &target)
		} else {
			parentStack = &[]*Element{}
			nestCount = 0
			target = elements
		}

		if processHeader(target, line) ||
			processHorizontalRule(target, line, &i) ||
			processVariableLine(target, line) {
			continue
		}
	}

	return *elements
}

func identifyListedItem(line string) (bool, int, ListType, string) {
	trimmed := strings.TrimLeft(line, " \t")
	listType := ListNone
	if len(trimmed) >= 2 && trimmed[0] == '-' && trimmed[1] == ' ' {
		listType = ListUnordered
	}
	if len(trimmed) >= 3 &&
		strings.Contains("123456789", string(trimmed[0])) &&
		trimmed[1] == '.' &&
		trimmed[2] == ' ' {
		listType = ListOrdered
	}
	if listType != ListNone {
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

		if listType == ListUnordered {
			trimmed = trimmed[2:]
		} else {
			trimmed = trimmed[3:]
		}

		return true, (spaceCount / 2) + 1, listType, trimmed
	}

	return false, 0, ListNone, ""
}

func handleListItem(listType ListType, nestCount, generation *int, parentStack, elements *[]*Element, target **[]*Element) {
	if *nestCount < *generation {
		var targetParent *Element
		var rootParent *Element
		// create as many parents as required and link them in lineage order, and keep a pointer to the root parent
		for i := 0; i < *generation-*nestCount; i++ {
			parent := &Element{Kind: KList, ListKind: listType, Children: []*Element{}}
			if i == 0 {
				rootParent = parent
			} else {
				targetParent.Children = append(targetParent.Children, parent)
			}
			targetParent = parent

			*parentStack = append(*parentStack, parent)
		}
		// add the parent built with chidlren (if required) to the current target
		**target = append(**target, rootParent)

		// now we need to move the pointer to the most junior child
		*target = &(*parentStack)[len(*parentStack)-1].Children
		*nestCount = *generation
	} else if *nestCount > *generation {
		if len(*parentStack) > 0 {
			for i := 0; i < *nestCount-*generation; i++ {
				*parentStack = (*parentStack)[:len(*parentStack)-1]
			}
		}
		*nestCount = *generation

		if len(*parentStack) > 0 {
			*target = &(*parentStack)[len(*parentStack)-1].Children
		} else {
			*target = elements
		}
	}
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
			*elements = append(*elements, &Element{Kind: KHeading, Level: hashCount, LineBreak: true, Text: strings.TrimLeft(line, "# ")})
			return isHeader
		}
	}
	return false
}

func processHorizontalRule(elements *[]*Element, line string, index *int) bool {
	if strings.HasPrefix(line, "---") {
		invalidChar := strings.ContainsFunc(line, func(r rune) bool {
			if r != ' ' && r != '-' {
				return true
			}
			return false
		})

		if !invalidChar {
			*elements = append(*elements, &Element{Kind: KRule, LineBreak: true, Text: "\n" + line + "\n"})
			*index = *index + 1
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
	ruleString := "!*`[()]_"
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
			switch char {
			case "*":
				handleBold(elements, line, &specialChars, &cache, &basePointer, &lookAheadPointer)
				continue
			case "_":
				handleItalic(elements, line, &specialChars, &cache, &basePointer, &lookAheadPointer)
				continue
			case "[":
				handleLink(elements, line, &specialChars, &cache, &basePointer, &lookAheadPointer)
				continue
			case "!":
				handleImage(elements, line, &specialChars, &cache, &basePointer, &lookAheadPointer)
				continue
			case "`":
				handleCode(elements, line, &specialChars, &cache, &basePointer, &lookAheadPointer)
				continue
			}
		}
		cache = append(cache, line[basePointer])
	}
	if len(cache) != 0 || len(line) == 1 {
		*elements = append(*elements, &Element{Kind: KText, LineBreak: true, Text: string(cache) + string(line[len(line)-1])})
	}

	return false
}

func handleItalic(elements *[]*Element, line string, specialChars *[]*IndexChar, cache *[]byte, basePointer, lookAheadPointer *int) {
	// let's unshift 1 from specialChars
	*specialChars = (*specialChars)[1:]

	found := false
	for _, indexChar := range *specialChars {
		if indexChar.c == '_' {
			*lookAheadPointer = indexChar.i
			found = true
			break
		}
	}

	if found {
		italicText := line[*basePointer+1 : *lookAheadPointer]
		// clear the cache to a text element
		if len(*cache) != 0 {
			*elements = append(*elements, &Element{Kind: KText, Text: string(*cache)})
			*cache = []byte{}
		}
		// if there are no more chars we use the ln version
		if *lookAheadPointer == len(line)-1 {
			*elements = append(*elements, &Element{Kind: KItalic, LineBreak: true, Text: "_" + italicText + "_"})
		} else {
			*elements = append(*elements, &Element{Kind: KItalic, Text: "_" + italicText + "_"})
		}

		// let's remove special chars that are no more of use to us
		newSpecialChars := []*IndexChar{}
		count := 0
		for i, specialChar := range *specialChars {
			if count == 1 {
				newSpecialChars = append(newSpecialChars, (*specialChars)[i:]...)
				break
			}
			if specialChar.c == '_' {
				count++
			} else {
				newSpecialChars = append(newSpecialChars, specialChar)
			}
		}

		*specialChars = newSpecialChars
		// shift the pointer up, it gets incremented by 1 in the loop
		*basePointer = *lookAheadPointer
	} else {
		// we can ignore the underscore character as this is not an italic block
		*cache = append(*cache, line[*basePointer])
	}
}

func handleBold(elements *[]*Element, line string, specialChars *[]*IndexChar, cache *[]byte, basePointer, lookAheadPointer *int) {
	// is the next char a *? if so it's bold opener
	if line[*basePointer+1] == '*' {
		// let's unshift 2 from specialChars
		*specialChars = (*specialChars)[2:]

		found := false
		for _, indexChar := range *specialChars {
			if indexChar.c == '*' {
				*lookAheadPointer = indexChar.i
				found = true
				break
			}
		}

		if found && line[*lookAheadPointer+1] == '*' {
			// now we have a range of text in the bold between the basePointer and lookAheadPointer that is bold
			boldText := line[*basePointer+2 : *lookAheadPointer]
			// clear the cache to a text element
			if len(*cache) != 0 {
				*elements = append(*elements, &Element{Kind: KText, Text: string(*cache)})
				*cache = []byte{}
			}
			// if there are no more chars then we use the ln version
			if *lookAheadPointer+1 == len(line)-1 {
				*elements = append(*elements, &Element{Kind: KBold, LineBreak: true, Text: "**" + boldText + "**"})
			} else {
				*elements = append(*elements, &Element{Kind: KBold, Text: "**" + boldText + "**"})
			}

			// remove the further double * from specialChars
			newSpecialChars := []*IndexChar{}
			count := 0
			for i, specialChar := range *specialChars {
				if count == 2 {
					newSpecialChars = append(newSpecialChars, (*specialChars)[i:]...)
					break
				}
				if specialChar.c == '*' {
					count++
				} else {
					newSpecialChars = append(newSpecialChars, specialChar)
				}
			}
			*specialChars = newSpecialChars
			// we can increment the lookaheadpointer by one as it is currently targeting the first closing *
			*lookAheadPointer = *lookAheadPointer + 1
			// shift the pointer up, it gets incremented by 1 in the loop
			*basePointer = *lookAheadPointer
		} else {
			// we can ignore the * character, and the following one as there is no bold block
			*cache = append(*cache, line[*basePointer], line[*basePointer+1])
			*basePointer++
		}
	}

	return
}

func handleLink(elements *[]*Element, line string, specialChars *[]*IndexChar, cache *[]byte, basePointer, lookAheadPointer *int) { // link
	// let's unshift 1 from specialChars
	*specialChars = (*specialChars)[1:]

	found := false
	tempLookAheadPointer := 0
	display := ""
	for _, indexChar := range *specialChars {
		if indexChar.c == ']' {
			found = true
			tempLookAheadPointer = indexChar.i
			break
		}
	}
	if found && string(line[tempLookAheadPointer+1]) == "(" {
		// ok we have a square brackets enclosed text
		display = line[*basePointer+1 : tempLookAheadPointer]

		// continue the same again but for parenthesis
		foundLink := false
		link := ""

		for _, indexChar := range *specialChars {
			if indexChar.i <= tempLookAheadPointer+1 {
				continue
			}
			if indexChar.c == ')' {
				foundLink = true
				*lookAheadPointer = indexChar.i
				break
			}
		}
		if foundLink {
			// we definitely have a link
			link = line[tempLookAheadPointer+2 : *lookAheadPointer]
			// is there cache to append?
			if len(*cache) != 0 {
				*elements = append(*elements, &Element{Kind: KText, Text: string(*cache)})
				*cache = []byte{}
			}
			// if there are no more chars we use the ln version
			if *lookAheadPointer == len(line)-1 {
				*elements = append(*elements, &Element{Kind: KLink, LineBreak: true, Text: display, Href: link})
			} else {
				*elements = append(*elements, &Element{Kind: KLink, Text: display, Href: link})
			}
			// let's clear out the special chars that are of no use to us
			newSpecialChars := []*IndexChar{}
			count := 0
			foundClosingSquare := false
			for i, specialChar := range *specialChars {
				if count == 1 {
					newSpecialChars = append(newSpecialChars, (*specialChars)[i:]...)
					break
				}
				if specialChar.c == ']' {
					foundClosingSquare = true
				} else if specialChar.c == ')' && foundClosingSquare {
					count++
				} else if strings.Contains(string(specialChar.c), "[]()") {
					newSpecialChars = append(newSpecialChars, specialChar)
				}
			}
			*specialChars = newSpecialChars
			// shift the pointer up, it gets incremented by 1 in the loop
			*basePointer = *lookAheadPointer + 1
		} else {
			*cache = append(*cache, line[*basePointer])
		}
	} else {
		*cache = append(*cache, line[*basePointer])
	}
}

func handleImage(elements *[]*Element, line string, specialChars *[]*IndexChar, cache *[]byte, basePointer, lookAheadPointer *int) { // link
	*specialChars = (*specialChars)[1:]

	if string(line[*basePointer+1]) == "[" {
		// so we have an ! and a [
		found := false
		tempLookAheadPointer := 0
		display := ""
		for _, indexChar := range *specialChars {
			if indexChar.c == ']' {
				found = true
				tempLookAheadPointer = indexChar.i
				break
			}
		}
		if found && string(line[tempLookAheadPointer+1]) == "(" {
			// ok we have a square brackets enclosed text
			display = line[*basePointer+2 : tempLookAheadPointer]

			// continue the same again but for parenthesis
			foundImage := false
			link := ""

			for _, indexChar := range *specialChars {
				if indexChar.i <= tempLookAheadPointer+1 {
					continue
				}
				if indexChar.c == ')' {
					foundImage = true
					*lookAheadPointer = indexChar.i
					break
				}
			}
			if foundImage {
				// we have
				link = line[tempLookAheadPointer+2 : *lookAheadPointer]
				// flush the cache to a text element
				if len(*cache) != 0 {
					*elements = append(*elements, &Element{Kind: KText, Text: string(*cache)})
					*cache = []byte{}
				}
				// images are always "ln"
				*elements = append(*elements, &Element{Kind: KImage, LineBreak: true, Alt: display, Href: link})

				// now get rid of the special chars that are of no use to us
				newSpecialChars := []*IndexChar{}
				count := 0
				foundClosingSquare := false
				for i, specialChar := range *specialChars {
					if count == 1 {
						newSpecialChars = append(newSpecialChars, (*specialChars)[i:]...)
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
				*specialChars = newSpecialChars
				// shift the pointer up, it gets incremented by 1 in the loop
				*basePointer = *lookAheadPointer + 1
			} else {
				*cache = append(*cache, line[*basePointer])
			}
		} else {
			*cache = append(*cache, line[*basePointer])
		}
	} else {
		*cache = append(*cache, line[*basePointer])
	}
}

func handleCode(elements *[]*Element, line string, specialChars *[]*IndexChar, cache *[]byte, basePointer, lookAheadPointer *int) {
	// let's unshift one char from specialChars
	*specialChars = (*specialChars)[1:]

	found := false
	for _, indexChar := range *specialChars {
		if indexChar.c == '`' {
			*lookAheadPointer = indexChar.i
			found = true
			break
		}
	}

	if found {
		// we have a code block
		codeText := line[*basePointer+1 : *lookAheadPointer]
		// flush the cache to a text element
		if len(*cache) != 0 {
			*elements = append(*elements, &Element{Kind: KText, Text: string(*cache)})
			*cache = []byte{}
		}
		// if there are no more chars we create an ln element
		if *lookAheadPointer == len(line)-1 {
			*elements = append(*elements, &Element{Kind: KCodeSpan, LineBreak: true, Text: "`" + codeText + "`"})
		} else {
			*elements = append(*elements, &Element{Kind: KCodeSpan, Text: "`" + codeText + "`"})
		}

		// let's get rid of the specialChars that are of no use to us
		newSpecialChars := []*IndexChar{}
		count := 0
		for i, specialChar := range *specialChars {
			if count == 1 {
				newSpecialChars = append(newSpecialChars, (*specialChars)[i:]...)
				break
			}
			if specialChar.c == '`' {
				count++
			} else {
				newSpecialChars = append(newSpecialChars, specialChar)
			}
		}
		*specialChars = newSpecialChars
		// shift the pointer up, it gets incremented by 1 in the loop
		*basePointer = *lookAheadPointer
	} else {
		// we can ignore the ` character
		*cache = append(*cache, line[*basePointer])
	}
}
