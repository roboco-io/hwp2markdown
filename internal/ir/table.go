package ir

// TableBlock represents a table region in the document.
type TableBlock struct {
	Rows      int      `json:"rows"`
	Cols      int      `json:"cols"`
	Cells     [][]Cell `json:"cells,omitempty"`
	RawText   string   `json:"raw_text,omitempty"`   // fallback: tab/newline separated text
	Caption   string   `json:"caption,omitempty"`    // table caption if any
	HasHeader bool     `json:"has_header,omitempty"` // first row is header
}

// Cell represents a single cell in a table.
type Cell struct {
	Text    string    `json:"text"`
	RowSpan int       `json:"row_span,omitempty"` // number of rows this cell spans
	ColSpan int       `json:"col_span,omitempty"` // number of columns this cell spans
	Style   CellStyle `json:"style,omitempty"`
}

// CellStyle contains cell-level styling hints.
type CellStyle struct {
	Bold      bool   `json:"bold,omitempty"`
	Alignment string `json:"alignment,omitempty"` // left, center, right
	IsHeader  bool   `json:"is_header,omitempty"`
}

// NewTable creates a new table with the specified dimensions.
func NewTable(rows, cols int) *TableBlock {
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
		for j := range cells[i] {
			cells[i][j] = Cell{
				RowSpan: 1,
				ColSpan: 1,
			}
		}
	}
	return &TableBlock{
		Rows:  rows,
		Cols:  cols,
		Cells: cells,
	}
}

// NewTableFromRawText creates a table from raw text (tab/newline separated).
func NewTableFromRawText(rawText string, rows, cols int) *TableBlock {
	return &TableBlock{
		Rows:    rows,
		Cols:    cols,
		RawText: rawText,
	}
}

// SetCell sets the content of a specific cell.
func (t *TableBlock) SetCell(row, col int, text string) {
	if row >= 0 && row < t.Rows && col >= 0 && col < t.Cols && t.Cells != nil {
		t.Cells[row][col].Text = text
	}
}

// GetCell returns the cell at the specified position.
func (t *TableBlock) GetCell(row, col int) *Cell {
	if row >= 0 && row < t.Rows && col >= 0 && col < t.Cols && t.Cells != nil {
		return &t.Cells[row][col]
	}
	return nil
}

// SetHeaderRow marks the first row as a header row.
func (t *TableBlock) SetHeaderRow() {
	t.HasHeader = true
	if t.Cells != nil && t.Rows > 0 {
		for j := 0; j < t.Cols; j++ {
			t.Cells[0][j].Style.IsHeader = true
		}
	}
}
