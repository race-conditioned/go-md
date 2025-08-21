package gomd

import (
	"context"
	"strings"
)

// reset clears the OnePassParser's state, allowing it to be reused for a new parse operation.
func (p *OnePassParser) reset() {
	p.text = ""
	p.elements = []*Element{}
	p.leafNode = &p.elements
	p.parentStack = []*Element{}
	p.err = nil
}

// reset resets the variableLineCtx to its initial state, clearing pointers and caches.
func (ctx *variableLineCtx) reset() {
	ctx.basePointer = 0
	ctx.lookAheadPointer = 0
	ctx.specialChars = []indexChar{}
	ctx.cache = []byte{}
}

// seek searches for the next occurrence of a rune in the specialChars slice.
func (ctx *variableLineCtx) seek(r rune) (found bool) {
	for _, indexChar := range ctx.specialChars {
		if indexChar.c == r {
			ctx.lookAheadPointer = indexChar.i
			found = true
			break
		}
	}
	return found
}

// canceled checks if the context has been canceled or has an error.
func (p *OnePassParser) canceled() bool {
	if p.ctx == nil {
		return false
	}
	if err := p.ctx.Err(); err != nil {
		p.err = err
		return true
	}
	return false
}

// appendElement adds the element pointer to the leaf node
func (p *OnePassParser) appendElement(e *Element) {
	*p.leafNode = append(*p.leafNode, e)
}

// flushCtxCache checks if there is any cached text in the lineCtx and appends it as a text Element.
func (p *OnePassParser) flushCtxCache() {
	if len(p.lineCtx.cache) != 0 {
		p.appendElement(&Element{Kind: EKText, Text: string(p.lineCtx.cache)})
		p.lineCtx.cache = []byte{}
	}
}

// Parse parses the provided Markdown string and returns a slice of Elements.
func (p *OnePassParser) Parse(md string) *Document {
	document, _ := p.ParseCtx(context.Background(), md)
	return document
}

// ParseCtx parses the provided Markdown string in the context of the provided context.Context.
func (p *OnePassParser) ParseCtx(ctx context.Context, md string) (*Document, error) {
	p.ctx = ctx
	p.reset()
	lines := strings.Split(md, "\n")
	nestCount := 0

	for i := 0; i < len(lines); i++ {
		if p.canceled() {
			return nil, p.err
		}

		p.text = lines[i]
		if len(p.text) == 0 && i < len(lines)-1 {
			// is the next line a rule?
			if i <= len(lines)-1 && len(lines) > 0 && lines[i+1] == "---" {
				continue
			}

			p.appendElement(&Element{Kind: EKNewLine, LineBreak: true})
		}

		// we allow for switching between the elements and Children slices
		isListItem, generation, listType := p.identifyListedItem()
		if isListItem {
			nestCount = p.handleListItem(listType, nestCount, generation)
		} else {
			p.parentStack = (p.parentStack)[:0]
			nestCount = 0
			p.leafNode = &p.elements
		}

		if p.processHeader() ||
			p.processHorizontalRule(&i) ||
			p.processVariableLine() {
			continue
		}
	}

	return &Document{Elements: p.elements}, nil
}

// identifyListedItem checks if the current line starts with a list item marker (either unordered or ordered).
func (p *OnePassParser) identifyListedItem() (bool, int, ListType) {
	trimmed := strings.TrimLeft(p.text, " \t")
	listType := ListNone

	if len(trimmed) >= 2 && trimmed[0] == '-' && trimmed[1] == ' ' {
		listType = ListUnordered
	}

	// TODO: update for parsing lists longer than 9
	if len(trimmed) >= 3 &&
		strings.Contains("123456789", string(trimmed[0])) &&
		trimmed[1] == '.' &&
		trimmed[2] == ' ' {
		listType = ListOrdered
	}

	if listType == ListNone {
		return false, 0, ListNone
	}

	spaceCount := 0
	for _, r := range p.text {
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

	p.text = trimmed
	return true, (spaceCount / 2) + 1, listType
}

// handleListItem processes a list item based on its type and nesting level.
func (p *OnePassParser) handleListItem(listType ListType, nestCount int, generation int) int {
	if nestCount < generation {
		var targetParent *Element
		var rootParent *Element
		// create as many parents as required and link them in lineage order, and keep a pointer to the root parent
		for i := 0; i < generation-nestCount; i++ {
			parent := &Element{Kind: EKList, ListKind: listType, Children: []*Element{}}
			if i == 0 {
				rootParent = parent
			} else {
				targetParent.Children = append(targetParent.Children, parent)
			}
			targetParent = parent
			p.parentStack = append(p.parentStack, parent)
		}
		// add the parent built with chidlren (if required) to the current target
		p.appendElement(rootParent)

		// now we need to move the pointer to the most junior child
		p.leafNode = &(p.parentStack)[len(p.parentStack)-1].Children
		nestCount = generation
		return nestCount
	}

	if len(p.parentStack) == 0 {
		p.leafNode = &p.elements
		return generation
	}

	for i := 0; i < nestCount-generation; i++ {
		p.parentStack = (p.parentStack)[:len(p.parentStack)-1]
	}

	p.leafNode = &(p.parentStack)[len(p.parentStack)-1].Children

	return nestCount
}

// processHeader determines if the line has a valid header by counting hashes and checking for a space that follows immediately.
// it returns true and appends the Element pointer to the dereferenced *[]*Element slice if it identifies a valid header.
// it returns false and does nothing if no valid header is found.
// it identifies the header type based on the hash count
func (p *OnePassParser) processHeader() bool {
	if !strings.HasPrefix(p.text, "#") {
		return false
	}
	level := 0
	isHeader := true

	for i, r := range p.text {
		if r == '#' {
			level = i + 1
		} else if level > 0 {
			if r != ' ' {
				isHeader = false
			}
			break
		} else {
			isHeader = false
		}
	}

	if isHeader {
		p.appendElement(&Element{Kind: EKHeading, Level: level, LineBreak: true, Text: strings.TrimRight(p.text[level+1:], " ")})
	}
	return isHeader
}

// processHorizontalRule checks if the line starts with "---" and contains only valid characters.
func (p *OnePassParser) processHorizontalRule(index *int) bool {
	if !strings.HasPrefix(p.text, "---") {
		return false
	}

	invalidChar := strings.ContainsFunc(p.text, func(r rune) bool {
		if r != ' ' && r != '-' {
			return true
		}
		return false
	})

	if !invalidChar {
		p.appendElement(&Element{Kind: EKRule, LineBreak: true, Text: "\n" + p.text + "\n"})
		*index = *index + 1
	}
	return !invalidChar
}

// processVariableLine processes a line of text for Markdown syntax elements such as bold, italic, links, images, and code spans.
func (p *OnePassParser) processVariableLine() bool {
	p.lineCtx.reset()
	for i, r := range p.text {
		if strings.Contains(p.lineCtx.ruleString, string(r)) {
			p.lineCtx.specialChars = append(p.lineCtx.specialChars, indexChar{i: i, c: r})
		}
	}

	for ; p.lineCtx.basePointer < len(p.text)-1; p.lineCtx.basePointer++ {
		char := string(p.text[p.lineCtx.basePointer])
		if strings.Contains(p.lineCtx.ruleString, char) {
			switch char {
			case "*":
				p.handleBold(&p.lineCtx)
				continue
			case "_":
				p.handleItalic(&p.lineCtx)
				continue
			case "[":
				p.handleLink(&p.lineCtx)
				continue
			case "!":
				p.handleImage(&p.lineCtx)
				continue
			case "`":
				p.handleCode(&p.lineCtx)
				continue
			}
		}
		p.lineCtx.cache = append(p.lineCtx.cache, p.text[p.lineCtx.basePointer])
	}
	if len(p.lineCtx.cache) != 0 || len(p.text) == 1 {
		last := ""
		if len(p.text) > 0 {
			last = string(p.text[len(p.text)-1])
		}
		p.appendElement(&Element{Kind: EKText, LineBreak: true, Text: string(p.lineCtx.cache) + last})
	}

	return false
}

// handleBold processes bold text enclosed in double asterisks "**...**".
func (p *OnePassParser) handleBold(ctx *variableLineCtx) {
	if p.err != nil {
		return
	}
	if p.text[ctx.basePointer+1] != '*' {
		return
	}

	ctx.specialChars = ctx.specialChars[2:]

	if !ctx.seek('*') || p.text[ctx.lookAheadPointer+1] != '*' {
		return
	}

	// now we have a range of text in the bold between the basePointer and lookAheadPointer that is bold
	boldText := p.text[ctx.basePointer+2 : ctx.lookAheadPointer]
	// clear the cache to a text element
	p.flushCtxCache()
	// if there are no more chars then we use the ln version
	if ctx.lookAheadPointer+1 == len(p.text)-1 {
		p.appendElement(&Element{Kind: EKBold, LineBreak: true, Text: inlineWrap("**", boldText)})
	} else {
		p.appendElement(&Element{Kind: EKBold, Text: inlineWrap("**", boldText)})
	}

	// remove the further double * from specialChars
	newSpecialChars := []indexChar{}
	count := 0
	for i, specialChar := range ctx.specialChars {
		if count == 2 {
			newSpecialChars = append(newSpecialChars, ctx.specialChars[i:]...)
			break
		}
		if specialChar.c == '*' {
			count++
		} else {
			newSpecialChars = append(newSpecialChars, specialChar)
		}
	}
	ctx.specialChars = newSpecialChars
	// we can increment the lookaheadpointer by one as it is currently targeting the first closing *
	ctx.lookAheadPointer = ctx.lookAheadPointer + 1
	// shift the pointer up, it gets incremented by 1 in the loop
	ctx.basePointer = ctx.lookAheadPointer
}

// handleItalic processes italic text enclosed in single underscores "_..._".
func (p *OnePassParser) handleItalic(ctx *variableLineCtx) {
	if p.err != nil {
		return
	}
	// let's unshift 1 from specialChars
	ctx.specialChars = ctx.specialChars[1:]

	if !ctx.seek('_') {
		return
	}

	italicText := p.text[ctx.basePointer+1 : ctx.lookAheadPointer]
	// clear the cache to a text element

	p.flushCtxCache()
	// if there are no more chars we use the ln version
	if ctx.lookAheadPointer == len(p.text)-1 {
		p.appendElement(&Element{Kind: EKItalic, LineBreak: true, Text: inlineWrap("_", italicText)})
	} else {
		p.appendElement(&Element{Kind: EKItalic, Text: inlineWrap("_", italicText)})
	}

	// let's remove special chars that are no more of use to us
	newSpecialChars := []indexChar{}
	count := 0
	for i, specialChar := range ctx.specialChars {
		if count == 1 {
			newSpecialChars = append(newSpecialChars, ctx.specialChars[i:]...)
			break
		}
		if specialChar.c == '_' {
			count++
		} else {
			newSpecialChars = append(newSpecialChars, specialChar)
		}
	}

	ctx.specialChars = newSpecialChars
	// shift the pointer up, it gets incremented by 1 in the loop
	ctx.basePointer = ctx.lookAheadPointer
}

// handleItalic processes italic text enclosed in single underscores "_..._".
func (p *OnePassParser) handleLink(ctx *variableLineCtx) {
	if p.err != nil {
		return
	}
	// let's unshift 1 from specialChars
	ctx.specialChars = ctx.specialChars[1:]

	found := false
	tempLookAheadPointer := 0
	for _, indexChar := range ctx.specialChars {
		if indexChar.c == ']' {
			found = true
			tempLookAheadPointer = indexChar.i
			break
		}
	}
	if !found || string(p.text[tempLookAheadPointer+1]) != "(" {
		return
	}
	// ok we have a square brackets enclosed text
	display := p.text[ctx.basePointer+1 : tempLookAheadPointer]

	// continue the same again but for parenthesis
	foundLink := false
	for _, indexChar := range ctx.specialChars {
		if indexChar.i <= tempLookAheadPointer+1 {
			continue
		}
		if indexChar.c == ')' {
			foundLink = true
			ctx.lookAheadPointer = indexChar.i
			break
		}
	}

	if !foundLink {
		return
	}

	// we definitely have a link
	link := p.text[tempLookAheadPointer+2 : ctx.lookAheadPointer]
	// flush the cache
	p.flushCtxCache()
	// if there are no more chars we use the ln version
	if ctx.lookAheadPointer == len(p.text)-1 {
		p.appendElement(&Element{Kind: EKLink, LineBreak: true, Text: display, Href: link})
	} else {
		p.appendElement(&Element{Kind: EKLink, Text: display, Href: link})
	}
	// let's clear out the special chars that are of no use to us
	newSpecialChars := []indexChar{}
	count := 0
	foundClosingSquare := false
	for i, specialChar := range ctx.specialChars {
		if count == 1 {
			newSpecialChars = append(newSpecialChars, ctx.specialChars[i:]...)
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
	ctx.specialChars = newSpecialChars
	// shift the pointer up, it gets incremented by 1 in the loop
	ctx.basePointer = ctx.lookAheadPointer + 1
}

// handleImage processes images in Markdown syntax, which are similar to links but start with an exclamation mark.
func (p *OnePassParser) handleImage(ctx *variableLineCtx) {
	if p.err != nil {
		return
	}

	ctx.specialChars = ctx.specialChars[1:]

	if string(p.text[ctx.basePointer+1]) != "[" {
		return
	}

	// so we have an ! and a [
	found := false
	tempLookAheadPointer := 0
	for _, indexChar := range ctx.specialChars {
		if indexChar.c == ']' {
			found = true
			tempLookAheadPointer = indexChar.i
			break
		}
	}

	if !found || string(p.text[tempLookAheadPointer+1]) != "(" {
		return
	}

	// ok we have a square brackets enclosed text
	alt := p.text[ctx.basePointer+2 : tempLookAheadPointer]

	// continue the same again but for parenthesis
	foundImage := false

	for _, indexChar := range ctx.specialChars {
		if indexChar.i <= tempLookAheadPointer+1 {
			continue
		}
		if indexChar.c == ')' {
			foundImage = true
			ctx.lookAheadPointer = indexChar.i
			break
		}
	}

	if !foundImage {
		return
	}

	// we have
	href := p.text[tempLookAheadPointer+2 : ctx.lookAheadPointer]
	// flush the cache to a text element
	p.flushCtxCache()
	// images are always "ln"
	p.appendElement(&Element{Kind: EKImage, LineBreak: true, Alt: alt, Href: href})

	// now get rid of the special chars that are of no use to us
	newSpecialChars := []indexChar{}
	count := 0
	foundClosingSquare := false
	for i, specialChar := range ctx.specialChars {
		if count == 1 {
			newSpecialChars = append(newSpecialChars, ctx.specialChars[i:]...)
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
	ctx.specialChars = newSpecialChars
	// shift the pointer up, it gets incremented by 1 in the loop
	ctx.basePointer = ctx.lookAheadPointer + 1
}

// handleCode processes inline code spans enclosed in backticks "`...`".
func (p *OnePassParser) handleCode(ctx *variableLineCtx) {
	if p.err != nil {
		return
	}
	// let's unshift one char from specialChars
	ctx.specialChars = ctx.specialChars[1:]

	if !ctx.seek('`') {
		return
	}

	// we have a code block
	codeText := p.text[ctx.basePointer+1 : ctx.lookAheadPointer]
	// flush the cache to a text element
	p.flushCtxCache()
	// if there are no more chars we create an ln element
	if ctx.lookAheadPointer == len(p.text)-1 {
		p.appendElement(&Element{Kind: EKCodeSpan, LineBreak: true, Text: "`" + codeText + "`"})
	} else {
		p.appendElement(&Element{Kind: EKCodeSpan, Text: "`" + codeText + "`"})
	}

	// let's get rid of the specialChars that are of no use to us
	newSpecialChars := []indexChar{}
	count := 0
	for i, specialChar := range ctx.specialChars {
		if count == 1 {
			newSpecialChars = append(newSpecialChars, ctx.specialChars[i:]...)
			break
		}
		if specialChar.c == '`' {
			count++
		} else {
			newSpecialChars = append(newSpecialChars, specialChar)
		}
	}
	ctx.specialChars = newSpecialChars
	// shift the pointer up, it gets incremented by 1 in the loop
	ctx.basePointer = ctx.lookAheadPointer
}
