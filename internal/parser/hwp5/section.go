package hwp5

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Section은 본문 섹션 데이터
type Section struct {
	Paragraphs []*Paragraph
	Tables     []*Table
	Images     []*Image
}

// Paragraph는 문단 데이터
type Paragraph struct {
	Text           string
	Controls       []ControlInfo
	CharShapeID    uint16
	ParaShapeID    uint16
	StyleID        uint16
	DivisionType   uint8
	CharShapeCount uint16
	RangeTagCount  uint16
	LineAlignCount uint16
	InstanceID     uint32
}

// Table은 표 데이터
type Table struct {
	Rows         int
	Cols         int
	Cells        [][]*TableCell
	BorderFill   uint16
	CellSpacing  int16
	LeftMargin   int16
	RightMargin  int16
	TopMargin    int16
	BottomMargin int16
}

// TableCell은 표 셀 데이터
type TableCell struct {
	Row        int
	Col        int
	RowSpan    int
	ColSpan    int
	Width      int
	Height     int
	Paragraphs []*Paragraph
}

// Image는 이미지 데이터
type Image struct {
	BinDataID   uint16
	BorderColor uint32
	Width       int32
	Height      int32
	XOffset     int32
	YOffset     int32
}

// SectionParser parses Section streams.
type SectionParser struct {
	records       []*Record
	textExtractor *TextExtractor
	docInfo       *DocInfo
}

// NewSectionParser creates a new section parser.
func NewSectionParser(docInfo *DocInfo) *SectionParser {
	return &SectionParser{
		textExtractor: NewTextExtractor(),
		docInfo:       docInfo,
	}
}

// Parse parses a section stream.
func (sp *SectionParser) Parse(data []byte) (*Section, error) {
	reader := NewRecordReader(data)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read section records: %w", err)
	}

	sp.records = records
	section := &Section{}

	i := 0
	for i < len(records) {
		rec := records[i]

		switch rec.TagID {
		case TagParaHeader:
			// Level 0 문단은 본문 문단
			if rec.Level == 0 {
				para, tables, nextIdx := sp.parseParagraphWithTables(i)
				if para != nil && strings.TrimSpace(para.Text) != "" {
					section.Paragraphs = append(section.Paragraphs, para)
				}
				// 문단 내에서 발견된 테이블 추가
				section.Tables = append(section.Tables, tables...)
				i = nextIdx
				continue
			}

		case TagCtrlHeader:
			// Level 0 또는 1인 CTRL_HEADER만 본문 레벨 테이블로 처리
			// (테이블 내 중첩 테이블은 Level이 더 높음)
			if rec.Level <= 1 && len(rec.Data) >= 4 {
				ctrlID := string(rec.Data[0:4])
				// " lbt" 또는 "tbl " = 테이블
				if ctrlID == " lbt" || ctrlID == "tbl " {
					table, nextIdx := sp.parseTableBlock(i)
					if table != nil && table.Rows > 0 && table.Cols > 0 {
						section.Tables = append(section.Tables, table)
					}
					i = nextIdx
					continue
				}
			}
		}

		i++
	}

	return section, nil
}

// parseParagraphWithTables parses a paragraph starting at index i and collects any tables within it.
// Returns the paragraph, any tables found, and the next index to process.
func (sp *SectionParser) parseParagraphWithTables(startIdx int) (*Paragraph, []*Table, int) {
	if startIdx >= len(sp.records) {
		return nil, nil, startIdx + 1
	}

	rec := sp.records[startIdx]
	if rec.TagID != TagParaHeader {
		return nil, nil, startIdx + 1
	}

	para := sp.parseParaHeader(rec.Data)
	startLevel := rec.Level
	i := startIdx + 1
	var tables []*Table

	// 다음 PARA_TEXT 찾기
	for i < len(sp.records) {
		nextRec := sp.records[i]

		// 같은 레벨 이하의 다른 PARA_HEADER가 나오면 종료
		if nextRec.TagID == TagParaHeader && nextRec.Level <= startLevel {
			break
		}

		// PARA_TEXT 처리
		if nextRec.TagID == TagParaText {
			text, controls := sp.textExtractor.ExtractTextWithControls(nextRec.Data)
			para.Text = text
			para.Controls = controls
		}

		// 테이블 컨트롤 처리 (문단 내 삽입된 테이블)
		if nextRec.TagID == TagCtrlHeader && len(nextRec.Data) >= 4 {
			ctrlID := string(nextRec.Data[0:4])
			if ctrlID == " lbt" || ctrlID == "tbl " {
				table, nextIdx := sp.parseTableBlock(i)
				if table != nil && table.Rows > 0 && table.Cols > 0 {
					tables = append(tables, table)
				}
				i = nextIdx
				continue
			}
		}

		i++
	}

	return para, tables, i
}

// parseTableBlock parses a table block starting at the CTRL_HEADER for the table.
// Returns the table and the next index to process.
func (sp *SectionParser) parseTableBlock(startIdx int) (*Table, int) {
	if startIdx >= len(sp.records) {
		return nil, startIdx + 1
	}

	ctrlRec := sp.records[startIdx]
	if ctrlRec.TagID != TagCtrlHeader {
		return nil, startIdx + 1
	}

	tableLevel := ctrlRec.Level
	i := startIdx + 1

	var table *Table
	var cells []*TableCell

	// 테이블 내 레코드 처리
	for i < len(sp.records) {
		rec := sp.records[i]

		// 테이블 레벨보다 낮은 레벨의 PARA_HEADER가 나오면 테이블 종료
		if rec.TagID == TagParaHeader && rec.Level <= tableLevel {
			break
		}

		// 같은 레벨의 다른 CTRL_HEADER가 나오면 테이블 종료
		if rec.TagID == TagCtrlHeader && rec.Level <= tableLevel {
			break
		}

		switch rec.TagID {
		case TagTable:
			table = sp.parseTableRecord(rec.Data)

		case TagListHeader:
			// LIST_HEADER는 셀을 나타냄
			// 이 셀의 문단들을 수집
			cell, nextIdx := sp.parseCellBlock(i, tableLevel)
			if cell != nil {
				cells = append(cells, cell)
			}
			i = nextIdx
			continue
		}

		i++
	}

	// 셀을 테이블에 배치
	if table != nil && len(cells) > 0 {
		sp.arrangeCellsInTable(table, cells)
	}

	return table, i
}

// parseCellBlock parses a cell starting at LIST_HEADER.
func (sp *SectionParser) parseCellBlock(startIdx int, tableLevel uint16) (*TableCell, int) {
	if startIdx >= len(sp.records) {
		return nil, startIdx + 1
	}

	listRec := sp.records[startIdx]
	if listRec.TagID != TagListHeader {
		return nil, startIdx + 1
	}

	cellLevel := listRec.Level
	cell := &TableCell{
		ColSpan: 1,
		RowSpan: 1,
	}

	i := startIdx + 1
	var paragraphs []*Paragraph

	// 셀 내 레코드 처리
	for i < len(sp.records) {
		rec := sp.records[i]

		// 셀 레벨보다 낮거나 같은 LIST_HEADER가 나오면 다음 셀
		if rec.TagID == TagListHeader && rec.Level <= cellLevel {
			break
		}

		// 테이블 레벨 이하의 레코드가 나오면 테이블 종료
		if rec.Level <= tableLevel && (rec.TagID == TagParaHeader || rec.TagID == TagCtrlHeader) {
			break
		}

		switch rec.TagID {
		case TagParaHeader:
			// 셀 내 문단
			para, nextIdx := sp.parseCellParagraph(i, cellLevel)
			if para != nil {
				paragraphs = append(paragraphs, para)
			}
			i = nextIdx
			continue
		}

		i++
	}

	cell.Paragraphs = paragraphs
	return cell, i
}

// parseCellParagraph parses a paragraph within a cell.
func (sp *SectionParser) parseCellParagraph(startIdx int, cellLevel uint16) (*Paragraph, int) {
	if startIdx >= len(sp.records) {
		return nil, startIdx + 1
	}

	rec := sp.records[startIdx]
	if rec.TagID != TagParaHeader {
		return nil, startIdx + 1
	}

	para := sp.parseParaHeader(rec.Data)
	paraLevel := rec.Level
	i := startIdx + 1

	// 문단 관련 레코드 처리
	for i < len(sp.records) {
		nextRec := sp.records[i]

		// 같은 레벨의 다른 PARA_HEADER가 나오면 종료
		if nextRec.TagID == TagParaHeader && nextRec.Level <= paraLevel {
			break
		}

		// 셀 레벨 이하의 LIST_HEADER가 나오면 종료
		if nextRec.TagID == TagListHeader && nextRec.Level <= cellLevel {
			break
		}

		// PARA_TEXT 처리
		if nextRec.TagID == TagParaText {
			text, controls := sp.textExtractor.ExtractTextWithControls(nextRec.Data)
			para.Text = text
			para.Controls = controls
		}

		i++
	}

	return para, i
}

// arrangeCellsInTable arranges cells into the table's 2D grid.
func (sp *SectionParser) arrangeCellsInTable(table *Table, cells []*TableCell) {
	if table.Rows <= 0 || table.Cols <= 0 {
		return
	}

	// 셀 배열 초기화
	table.Cells = make([][]*TableCell, table.Rows)
	for i := range table.Cells {
		table.Cells[i] = make([]*TableCell, table.Cols)
	}

	// 셀을 행 우선으로 배치
	cellIdx := 0
	for row := 0; row < table.Rows && cellIdx < len(cells); row++ {
		for col := 0; col < table.Cols && cellIdx < len(cells); col++ {
			// 이미 병합된 셀이 차지하고 있으면 건너뛰기
			if table.Cells[row][col] != nil {
				continue
			}

			cell := cells[cellIdx]
			cell.Row = row
			cell.Col = col
			table.Cells[row][col] = cell

			// 병합된 셀 처리 (rowspan, colspan)
			for r := row; r < row+cell.RowSpan && r < table.Rows; r++ {
				for c := col; c < col+cell.ColSpan && c < table.Cols; c++ {
					if r == row && c == col {
						continue
					}
					// 병합된 영역에 빈 셀 표시
					table.Cells[r][c] = &TableCell{
						Row:     r,
						Col:     c,
						RowSpan: 0, // 병합된 셀임을 표시
						ColSpan: 0,
					}
				}
			}

			cellIdx++
		}
	}
}

func (sp *SectionParser) parseParaHeader(data []byte) *Paragraph {
	if len(data) < 22 {
		return &Paragraph{}
	}

	para := &Paragraph{
		ParaShapeID:    binary.LittleEndian.Uint16(data[4:6]),
		StyleID:        uint16(data[6]),
		DivisionType:   data[7],
		CharShapeCount: binary.LittleEndian.Uint16(data[8:10]),
		RangeTagCount:  binary.LittleEndian.Uint16(data[10:12]),
		LineAlignCount: binary.LittleEndian.Uint16(data[12:14]),
		InstanceID:     binary.LittleEndian.Uint32(data[14:18]),
	}

	return para
}

func (sp *SectionParser) parseTableRecord(data []byte) *Table {
	if len(data) < 18 {
		return nil
	}

	// 표 속성 구조
	// [0:4] - 속성
	// [4:6] - 행 개수
	// [6:8] - 열 개수
	// [8:10] - 셀 간격
	// [10:12] - 왼쪽 여백
	// [12:14] - 오른쪽 여백
	// [14:16] - 위 여백
	// [16:18] - 아래 여백

	rows := int(binary.LittleEndian.Uint16(data[4:6]))
	cols := int(binary.LittleEndian.Uint16(data[6:8]))

	table := &Table{
		Rows:         rows,
		Cols:         cols,
		CellSpacing:  int16(binary.LittleEndian.Uint16(data[8:10])),
		LeftMargin:   int16(binary.LittleEndian.Uint16(data[10:12])),
		RightMargin:  int16(binary.LittleEndian.Uint16(data[12:14])),
		TopMargin:    int16(binary.LittleEndian.Uint16(data[14:16])),
		BottomMargin: int16(binary.LittleEndian.Uint16(data[16:18])),
	}

	return table
}

// GetCellText returns the text content of a cell.
func (c *TableCell) GetCellText() string {
	if c == nil || len(c.Paragraphs) == 0 {
		return ""
	}

	var texts []string
	for _, p := range c.Paragraphs {
		if text := strings.TrimSpace(p.Text); text != "" {
			texts = append(texts, text)
		}
	}

	return strings.Join(texts, "\n")
}
