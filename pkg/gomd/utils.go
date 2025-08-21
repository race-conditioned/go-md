package gomd

// Walk traverses the elements and applies the visit function to each element.
func Walk(elems []*Element, visit func(*Element)) {
	for _, el := range elems {
		if el == nil {
			continue
		}
		visit(el)
		Walk(el.Children, visit)
	}
}

// DeepCopy returns a deep copy of the element tree.
func DeepCopy(e *Element) *Element {
	if e == nil {
		return nil
	}
	// copy the element struct
	cp := *e

	// copy children recursively
	if len(e.Children) > 0 {
		cp.Children = make([]*Element, len(e.Children))
		for i, child := range e.Children {
			cp.Children[i] = DeepCopy(child)
		}
	}
	return &cp
}

// DeepCopySlice copies a slice of elements.
func DeepCopySlice(elems []*Element) []*Element {
	if len(elems) == 0 {
		return nil
	}
	out := make([]*Element, len(elems))
	for i, e := range elems {
		out[i] = DeepCopy(e)
	}
	return out
}

// inlineWrap concatenates the wrap to each side of s.
func inlineWrap(wrap, s string) string {
	return wrap + s + wrap
}

// btoi converts a boolean to an integer.
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
