package ir

import (
	"encoding/json"
	"testing"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument()

	if doc.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", doc.Version)
	}
	if len(doc.Content) != 0 {
		t.Errorf("expected empty content, got %d blocks", len(doc.Content))
	}
}

func TestDocument_AddParagraph(t *testing.T) {
	doc := NewDocument()
	p := NewParagraph("Hello, World!")

	doc.AddParagraph(p)

	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Content))
	}
	if doc.Content[0].Type != BlockTypeParagraph {
		t.Errorf("expected paragraph type, got %s", doc.Content[0].Type)
	}
	if doc.Content[0].Paragraph.Text != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %s", doc.Content[0].Paragraph.Text)
	}
}

func TestDocument_AddTable(t *testing.T) {
	doc := NewDocument()
	table := NewTable(2, 3)
	table.SetCell(0, 0, "Header 1")
	table.SetCell(0, 1, "Header 2")
	table.SetCell(0, 2, "Header 3")

	doc.AddTable(table)

	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Content))
	}
	if doc.Content[0].Type != BlockTypeTable {
		t.Errorf("expected table type, got %s", doc.Content[0].Type)
	}
	if doc.Content[0].Table.Rows != 2 {
		t.Errorf("expected 2 rows, got %d", doc.Content[0].Table.Rows)
	}
}

func TestDocument_AddImage(t *testing.T) {
	doc := NewDocument()
	img := NewImage("img001")
	img.Alt = "Test image"

	doc.AddImage(img)

	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Content))
	}
	if doc.Content[0].Type != BlockTypeImage {
		t.Errorf("expected image type, got %s", doc.Content[0].Type)
	}
	if doc.Content[0].Image.ID != "img001" {
		t.Errorf("expected 'img001', got %s", doc.Content[0].Image.ID)
	}
}

func TestDocument_AddList(t *testing.T) {
	doc := NewDocument()
	list := NewUnorderedList()
	list.AddItem("Item 1")
	list.AddItem("Item 2")

	doc.AddList(list)

	if len(doc.Content) != 1 {
		t.Fatalf("expected 1 block, got %d", len(doc.Content))
	}
	if doc.Content[0].Type != BlockTypeList {
		t.Errorf("expected list type, got %s", doc.Content[0].Type)
	}
	if len(doc.Content[0].List.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(doc.Content[0].List.Items))
	}
}

func TestDocument_JSONSerialization(t *testing.T) {
	doc := NewDocument()
	doc.Metadata.Title = "Test Document"
	doc.Metadata.Author = "Test Author"

	p := NewParagraph("Test paragraph")
	p.SetHeading(1)
	doc.AddParagraph(p)

	// Serialize
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Deserialize
	var restored Document
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify
	if restored.Version != doc.Version {
		t.Errorf("version mismatch: expected %s, got %s", doc.Version, restored.Version)
	}
	if restored.Metadata.Title != doc.Metadata.Title {
		t.Errorf("title mismatch: expected %s, got %s", doc.Metadata.Title, restored.Metadata.Title)
	}
	if len(restored.Content) != len(doc.Content) {
		t.Errorf("content length mismatch: expected %d, got %d", len(doc.Content), len(restored.Content))
	}
}

// Paragraph tests
func TestParagraph_SetHeading(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{6, 6},
		{7, 6},
		{10, 6},
	}

	for _, tc := range tests {
		p := NewParagraph("Test")
		p.SetHeading(tc.input)
		if p.Style.HeadingLevel != tc.expected {
			t.Errorf("SetHeading(%d): expected %d, got %d", tc.input, tc.expected, p.Style.HeadingLevel)
		}
	}
}

func TestParagraph_AddRun(t *testing.T) {
	p := NewParagraph("Base text")
	p.AddRun("bold", TextStyle{Bold: true})
	p.AddRun("italic", TextStyle{Italic: true})

	if len(p.Runs) != 2 {
		t.Errorf("expected 2 runs, got %d", len(p.Runs))
	}

	if !p.Runs[0].Style.Bold {
		t.Error("expected first run to be bold")
	}

	if !p.Runs[1].Style.Italic {
		t.Error("expected second run to be italic")
	}
}

func TestParagraph_IsEmpty(t *testing.T) {
	empty := NewParagraph("")
	if !empty.IsEmpty() {
		t.Error("expected empty paragraph to be empty")
	}

	notEmpty := NewParagraph("Text")
	if notEmpty.IsEmpty() {
		t.Error("expected non-empty paragraph to not be empty")
	}

	withRuns := NewParagraph("")
	withRuns.AddRun("run", TextStyle{})
	if withRuns.IsEmpty() {
		t.Error("expected paragraph with runs to not be empty")
	}
}

// Table tests
func TestTable_GetCell(t *testing.T) {
	table := NewTable(2, 3)
	table.SetCell(0, 0, "A1")
	table.SetCell(1, 2, "B3")

	cell := table.GetCell(0, 0)
	if cell == nil || cell.Text != "A1" {
		t.Error("expected cell (0,0) to be 'A1'")
	}

	cell = table.GetCell(1, 2)
	if cell == nil || cell.Text != "B3" {
		t.Error("expected cell (1,2) to be 'B3'")
	}

	// Out of bounds
	if table.GetCell(-1, 0) != nil {
		t.Error("expected nil for negative row")
	}
	if table.GetCell(0, -1) != nil {
		t.Error("expected nil for negative col")
	}
	if table.GetCell(10, 0) != nil {
		t.Error("expected nil for out of bounds row")
	}
	if table.GetCell(0, 10) != nil {
		t.Error("expected nil for out of bounds col")
	}
}

func TestTable_SetHeaderRow(t *testing.T) {
	table := NewTable(2, 3)
	table.SetHeaderRow()

	if !table.HasHeader {
		t.Error("expected HasHeader to be true")
	}

	for j := 0; j < 3; j++ {
		if !table.Cells[0][j].Style.IsHeader {
			t.Errorf("expected cell (0,%d) to be header", j)
		}
	}
}

func TestNewTableFromRawText(t *testing.T) {
	raw := "A\tB\tC\nD\tE\tF"
	table := NewTableFromRawText(raw, 2, 3)

	if table.RawText != raw {
		t.Error("expected RawText to match")
	}
	if table.Rows != 2 {
		t.Errorf("expected 2 rows, got %d", table.Rows)
	}
	if table.Cols != 3 {
		t.Errorf("expected 3 cols, got %d", table.Cols)
	}
}

// List tests
func TestNewOrderedList(t *testing.T) {
	list := NewOrderedList()
	if !list.Ordered {
		t.Error("expected ordered list")
	}
	if list.Start != 1 {
		t.Errorf("expected start=1, got %d", list.Start)
	}
}

func TestNewUnorderedList(t *testing.T) {
	list := NewUnorderedList()
	if list.Ordered {
		t.Error("expected unordered list")
	}
}

func TestList_AddItemWithLevel(t *testing.T) {
	list := NewUnorderedList()
	list.AddItemWithLevel("Level 0", 0)
	list.AddItemWithLevel("Level 1", 1)
	list.AddItemWithLevel("Level 2", 2)

	if len(list.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(list.Items))
	}

	if list.Items[1].Level != 1 {
		t.Errorf("expected level 1, got %d", list.Items[1].Level)
	}
}

func TestList_IsEmpty(t *testing.T) {
	list := NewUnorderedList()
	if !list.IsEmpty() {
		t.Error("expected empty list")
	}

	list.AddItem("Item")
	if list.IsEmpty() {
		t.Error("expected non-empty list")
	}
}

// Image tests
func TestNewImage(t *testing.T) {
	img := NewImage("test-id")
	if img.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %s", img.ID)
	}
}

func TestImage_SetDimensions(t *testing.T) {
	img := NewImage("img1")
	img.SetDimensions(800, 600)

	if img.Width != 800 {
		t.Errorf("expected width 800, got %d", img.Width)
	}
	if img.Height != 600 {
		t.Errorf("expected height 600, got %d", img.Height)
	}
}

func TestImage_HasData(t *testing.T) {
	img := NewImage("img1")
	if img.HasData() {
		t.Error("expected HasData to be false for new image")
	}

	img.Data = []byte{0x89, 0x50, 0x4E, 0x47}
	if !img.HasData() {
		t.Error("expected HasData to be true after setting data")
	}
}
