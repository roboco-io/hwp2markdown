package ir

// ListBlock represents a list (ordered or unordered) in the document.
type ListBlock struct {
	Ordered bool       `json:"ordered"` // true = numbered list, false = bullet list
	Items   []ListItem `json:"items"`
	Start   int        `json:"start,omitempty"` // starting number for ordered lists
}

// ListItem represents a single item in a list.
type ListItem struct {
	Text     string     `json:"text"`
	Level    int        `json:"level,omitempty"`    // nesting level (0 = top level)
	Children []ListItem `json:"children,omitempty"` // nested items
}

// NewList creates a new list block.
func NewList(ordered bool) *ListBlock {
	return &ListBlock{
		Ordered: ordered,
		Items:   make([]ListItem, 0),
		Start:   1,
	}
}

// NewOrderedList creates a new ordered (numbered) list.
func NewOrderedList() *ListBlock {
	return NewList(true)
}

// NewUnorderedList creates a new unordered (bullet) list.
func NewUnorderedList() *ListBlock {
	return NewList(false)
}

// AddItem adds an item to the list.
func (l *ListBlock) AddItem(text string) {
	l.Items = append(l.Items, ListItem{
		Text:  text,
		Level: 0,
	})
}

// AddItemWithLevel adds an item with a specific nesting level.
func (l *ListBlock) AddItemWithLevel(text string, level int) {
	l.Items = append(l.Items, ListItem{
		Text:  text,
		Level: level,
	})
}

// IsEmpty returns true if the list has no items.
func (l *ListBlock) IsEmpty() bool {
	return len(l.Items) == 0
}
