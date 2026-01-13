package hwp5

import (
	"encoding/binary"
	"fmt"
)

// DocInfo는 문서 정보 스트림에서 파싱된 데이터
type DocInfo struct {
	Properties  *DocumentProperties
	IDMappings  *IDMappings
	BinDataList []*BinDataInfo
	FaceNames   []string
	CharShapes  []*CharShape
	ParaShapes  []*ParaShape
	Styles      []*Style
}

// DocumentProperties는 문서 속성 (HWPTAG_DOCUMENT_PROPERTIES)
type DocumentProperties struct {
	SectionCount  uint16 // 구역 개수
	PageStartNum  uint16 // 시작 페이지 번호
	FootnoteStart uint16 // 각주 시작 번호
	EndnoteStart  uint16 // 미주 시작 번호
	PictureStart  uint16 // 그림 시작 번호
	TableStart    uint16 // 표 시작 번호
	EquationStart uint16 // 수식 시작 번호
	ListIDCount   uint32 // 리스트 ID 개수
	ParaIDCount   uint32 // 문단 ID 개수
	CharUnitLoc   uint32 // 글자 단위 위치
}

// IDMappings는 ID 매핑 테이블 크기 (HWPTAG_ID_MAPPINGS)
type IDMappings struct {
	BinDataCount           int32
	FaceNameKorCount       int32
	FaceNameEngCount       int32
	FaceNameHanCount       int32
	FaceNameJpnCount       int32
	FaceNameOtherCount     int32
	FaceNameSymCount       int32
	FaceNameUserCount      int32
	BorderFillCount        int32
	CharShapeCount         int32
	TabDefCount            int32
	NumberingCount         int32
	BulletCount            int32
	ParaShapeCount         int32
	StyleCount             int32
	MemoShapeCount         int32
	TrackChangeCount       int32
	TrackChangeAuthorCount int32
}

// BinDataInfo는 바이너리 데이터 정보 (HWPTAG_BIN_DATA)
type BinDataInfo struct {
	Type      uint16 // 바이너리 데이터 타입
	AbsPath   string // 절대 경로
	RelPath   string // 상대 경로
	BinDataID uint16 // BinData 스토리지 내 ID
	Extension string // 확장자
}

// CharShape는 글자 모양 (HWPTAG_CHAR_SHAPE)
type CharShape struct {
	FaceID       [7]uint16 // 언어별 글꼴 ID
	Ratios       [7]uint8  // 언어별 장평
	Spacings     [7]int8   // 언어별 자간
	RelSizes     [7]uint8  // 언어별 상대 크기
	Offsets      [7]int8   // 언어별 오프셋
	Height       int32     // 기준 크기 (100분의 1pt)
	Attributes   uint32    // 속성 플래그
	ShadowGap1   int8      // 그림자 간격 1
	ShadowGap2   int8      // 그림자 간격 2
	TextColor    uint32    // 글자 색
	UnderColor   uint32    // 밑줄 색
	ShadeColor   uint32    // 음영 색
	ShadowColor  uint32    // 그림자 색
	BorderFillID uint16    // 테두리/배경 ID
	StrikeColor  uint32    // 취소선 색
}

// ParaShape는 문단 모양 (HWPTAG_PARA_SHAPE)
type ParaShape struct {
	Attributes1     uint32 // 속성 1
	LeftMargin      int32  // 왼쪽 여백
	RightMargin     int32  // 오른쪽 여백
	Indent          int32  // 들여쓰기
	ParaSpaceBefore int32  // 문단 위 간격
	ParaSpaceAfter  int32  // 문단 아래 간격
	LineSpacing     int32  // 줄 간격
	TabDefID        uint16 // 탭 정의 ID
	NumberingID     uint16 // 문단 번호 ID
	BorderFillID    uint16 // 테두리/배경 ID
	BorderOffset1   int16  // 테두리 왼쪽 오프셋
	BorderOffset2   int16  // 테두리 오른쪽 오프셋
	BorderOffset3   int16  // 테두리 위 오프셋
	BorderOffset4   int16  // 테두리 아래 오프셋
	Attributes2     uint32 // 속성 2
	Attributes3     uint32 // 속성 3
	LineWrap        uint32 // 줄 나눔 기준
	AutoSpacing     uint32 // 자동 줄 간격
}

// Style은 스타일 정의 (HWPTAG_STYLE)
type Style struct {
	Name        string // 스타일 이름
	EngName     string // 영문 스타일 이름
	Type        uint8  // 스타일 타입
	NextStyleID uint8  // 다음 스타일 ID
	LangID      int16  // 언어 ID
	ParaShapeID uint16 // 문단 모양 ID
	CharShapeID uint16 // 글자 모양 ID
	LockForm    bool   // 잠금 여부
}

// ParseDocInfo parses the DocInfo stream.
func ParseDocInfo(data []byte) (*DocInfo, error) {
	reader := NewRecordReader(data)
	info := &DocInfo{}

	for {
		rec, err := reader.Read()
		if err != nil {
			break
		}

		switch rec.TagID {
		case TagDocumentProperties:
			info.Properties = parseDocumentProperties(rec.Data)
		case TagIDMappings:
			info.IDMappings = parseIDMappings(rec.Data)
		case TagBinData:
			binData := parseBinDataInfo(rec.Data)
			if binData != nil {
				info.BinDataList = append(info.BinDataList, binData)
			}
		case TagFaceName:
			name := parseFaceName(rec.Data)
			info.FaceNames = append(info.FaceNames, name)
		case TagCharShape:
			cs := parseCharShape(rec.Data)
			if cs != nil {
				info.CharShapes = append(info.CharShapes, cs)
			}
		case TagParaShape:
			ps := parseParaShape(rec.Data)
			if ps != nil {
				info.ParaShapes = append(info.ParaShapes, ps)
			}
		case TagStyle:
			style := parseStyle(rec.Data)
			if style != nil {
				info.Styles = append(info.Styles, style)
			}
		}
	}

	return info, nil
}

func parseDocumentProperties(data []byte) *DocumentProperties {
	if len(data) < 26 {
		return nil
	}

	return &DocumentProperties{
		SectionCount:  binary.LittleEndian.Uint16(data[0:2]),
		PageStartNum:  binary.LittleEndian.Uint16(data[2:4]),
		FootnoteStart: binary.LittleEndian.Uint16(data[4:6]),
		EndnoteStart:  binary.LittleEndian.Uint16(data[6:8]),
		PictureStart:  binary.LittleEndian.Uint16(data[8:10]),
		TableStart:    binary.LittleEndian.Uint16(data[10:12]),
		EquationStart: binary.LittleEndian.Uint16(data[12:14]),
		ListIDCount:   binary.LittleEndian.Uint32(data[14:18]),
		ParaIDCount:   binary.LittleEndian.Uint32(data[18:22]),
		CharUnitLoc:   binary.LittleEndian.Uint32(data[22:26]),
	}
}

func parseIDMappings(data []byte) *IDMappings {
	if len(data) < 72 {
		return nil
	}

	return &IDMappings{
		BinDataCount:           int32(binary.LittleEndian.Uint32(data[0:4])),
		FaceNameKorCount:       int32(binary.LittleEndian.Uint32(data[4:8])),
		FaceNameEngCount:       int32(binary.LittleEndian.Uint32(data[8:12])),
		FaceNameHanCount:       int32(binary.LittleEndian.Uint32(data[12:16])),
		FaceNameJpnCount:       int32(binary.LittleEndian.Uint32(data[16:20])),
		FaceNameOtherCount:     int32(binary.LittleEndian.Uint32(data[20:24])),
		FaceNameSymCount:       int32(binary.LittleEndian.Uint32(data[24:28])),
		FaceNameUserCount:      int32(binary.LittleEndian.Uint32(data[28:32])),
		BorderFillCount:        int32(binary.LittleEndian.Uint32(data[32:36])),
		CharShapeCount:         int32(binary.LittleEndian.Uint32(data[36:40])),
		TabDefCount:            int32(binary.LittleEndian.Uint32(data[40:44])),
		NumberingCount:         int32(binary.LittleEndian.Uint32(data[44:48])),
		BulletCount:            int32(binary.LittleEndian.Uint32(data[48:52])),
		ParaShapeCount:         int32(binary.LittleEndian.Uint32(data[52:56])),
		StyleCount:             int32(binary.LittleEndian.Uint32(data[56:60])),
		MemoShapeCount:         int32(binary.LittleEndian.Uint32(data[60:64])),
		TrackChangeCount:       int32(binary.LittleEndian.Uint32(data[64:68])),
		TrackChangeAuthorCount: int32(binary.LittleEndian.Uint32(data[68:72])),
	}
}

func parseBinDataInfo(data []byte) *BinDataInfo {
	if len(data) < 2 {
		return nil
	}

	info := &BinDataInfo{
		Type: binary.LittleEndian.Uint16(data[0:2]),
	}

	offset := 2
	dataType := info.Type & 0x0F

	// 타입에 따른 추가 파싱
	switch dataType {
	case 0: // LINK - 절대 경로 + 상대 경로
		if offset+2 <= len(data) {
			absPathLen := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
			if offset+absPathLen*2 <= len(data) {
				info.AbsPath = DecodeUTF16LE(data[offset : offset+absPathLen*2])
				offset += absPathLen * 2
			}
		}
		if offset+2 <= len(data) {
			relPathLen := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
			if offset+relPathLen*2 <= len(data) {
				info.RelPath = DecodeUTF16LE(data[offset : offset+relPathLen*2])
			}
		}
	case 1: // EMBEDDING - BinData ID
		if offset+2 <= len(data) {
			info.BinDataID = binary.LittleEndian.Uint16(data[offset : offset+2])
			offset += 2
		}
		if offset+2 <= len(data) {
			extLen := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
			if offset+extLen*2 <= len(data) {
				info.Extension = DecodeUTF16LE(data[offset : offset+extLen*2])
			}
		}
	case 2: // STORAGE - BinData ID
		if offset+2 <= len(data) {
			info.BinDataID = binary.LittleEndian.Uint16(data[offset : offset+2])
		}
	}

	return info
}

func parseFaceName(data []byte) string {
	if len(data) < 3 {
		return ""
	}

	// 속성 1바이트 건너뛰기
	offset := 1

	// 이름 길이 (2바이트)
	if offset+2 > len(data) {
		return ""
	}
	nameLen := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
	offset += 2

	// 이름 (UTF-16LE)
	if offset+nameLen*2 > len(data) {
		return ""
	}
	return DecodeUTF16LE(data[offset : offset+nameLen*2])
}

func parseCharShape(data []byte) *CharShape {
	if len(data) < 72 {
		return nil
	}

	cs := &CharShape{}

	// 언어별 글꼴 ID (7 * 2 = 14 bytes)
	for i := range 7 {
		cs.FaceID[i] = binary.LittleEndian.Uint16(data[i*2 : i*2+2])
	}

	// 언어별 장평 (7 bytes)
	copy(cs.Ratios[:], data[14:21])

	// 언어별 자간 (7 bytes)
	for i := range 7 {
		cs.Spacings[i] = int8(data[21+i])
	}

	// 언어별 상대 크기 (7 bytes)
	copy(cs.RelSizes[:], data[28:35])

	// 언어별 오프셋 (7 bytes)
	for i := range 7 {
		cs.Offsets[i] = int8(data[35+i])
	}

	// 기준 크기 (4 bytes)
	cs.Height = int32(binary.LittleEndian.Uint32(data[42:46]))

	// 속성 (4 bytes)
	cs.Attributes = binary.LittleEndian.Uint32(data[46:50])

	// 그림자 간격 (2 bytes)
	cs.ShadowGap1 = int8(data[50])
	cs.ShadowGap2 = int8(data[51])

	// 색상들 (4 bytes each)
	cs.TextColor = binary.LittleEndian.Uint32(data[52:56])
	cs.UnderColor = binary.LittleEndian.Uint32(data[56:60])
	cs.ShadeColor = binary.LittleEndian.Uint32(data[60:64])
	cs.ShadowColor = binary.LittleEndian.Uint32(data[64:68])

	// 테두리/배경 ID (2 bytes, 옵션)
	if len(data) >= 70 {
		cs.BorderFillID = binary.LittleEndian.Uint16(data[68:70])
	}

	// 취소선 색 (4 bytes, 옵션)
	if len(data) >= 74 {
		cs.StrikeColor = binary.LittleEndian.Uint32(data[70:74])
	}

	return cs
}

func parseParaShape(data []byte) *ParaShape {
	if len(data) < 54 {
		return nil
	}

	ps := &ParaShape{
		Attributes1:     binary.LittleEndian.Uint32(data[0:4]),
		LeftMargin:      int32(binary.LittleEndian.Uint32(data[4:8])),
		RightMargin:     int32(binary.LittleEndian.Uint32(data[8:12])),
		Indent:          int32(binary.LittleEndian.Uint32(data[12:16])),
		ParaSpaceBefore: int32(binary.LittleEndian.Uint32(data[16:20])),
		ParaSpaceAfter:  int32(binary.LittleEndian.Uint32(data[20:24])),
		LineSpacing:     int32(binary.LittleEndian.Uint32(data[24:28])),
		TabDefID:        binary.LittleEndian.Uint16(data[28:30]),
		NumberingID:     binary.LittleEndian.Uint16(data[30:32]),
		BorderFillID:    binary.LittleEndian.Uint16(data[32:34]),
		BorderOffset1:   int16(binary.LittleEndian.Uint16(data[34:36])),
		BorderOffset2:   int16(binary.LittleEndian.Uint16(data[36:38])),
		BorderOffset3:   int16(binary.LittleEndian.Uint16(data[38:40])),
		BorderOffset4:   int16(binary.LittleEndian.Uint16(data[40:42])),
		Attributes2:     binary.LittleEndian.Uint32(data[42:46]),
		Attributes3:     binary.LittleEndian.Uint32(data[46:50]),
		LineWrap:        binary.LittleEndian.Uint32(data[50:54]),
	}

	if len(data) >= 58 {
		ps.AutoSpacing = binary.LittleEndian.Uint32(data[54:58])
	}

	return ps
}

func parseStyle(data []byte) *Style {
	if len(data) < 2 {
		return nil
	}

	style := &Style{}
	offset := 0

	// 한글 이름 길이
	nameLen := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
	offset += 2

	// 한글 이름
	if offset+nameLen*2 <= len(data) {
		style.Name = DecodeUTF16LE(data[offset : offset+nameLen*2])
		offset += nameLen * 2
	}

	// 영문 이름 길이
	if offset+2 > len(data) {
		return style
	}
	engNameLen := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
	offset += 2

	// 영문 이름
	if offset+engNameLen*2 <= len(data) {
		style.EngName = DecodeUTF16LE(data[offset : offset+engNameLen*2])
		offset += engNameLen * 2
	}

	// 스타일 타입 (1 byte)
	if offset < len(data) {
		style.Type = data[offset]
		offset++
	}

	// 다음 스타일 ID (1 byte)
	if offset < len(data) {
		style.NextStyleID = data[offset]
		offset++
	}

	// 언어 ID (2 bytes)
	if offset+2 <= len(data) {
		style.LangID = int16(binary.LittleEndian.Uint16(data[offset : offset+2]))
		offset += 2
	}

	// 문단 모양 ID (2 bytes)
	if offset+2 <= len(data) {
		style.ParaShapeID = binary.LittleEndian.Uint16(data[offset : offset+2])
		offset += 2
	}

	// 글자 모양 ID (2 bytes)
	if offset+2 <= len(data) {
		style.CharShapeID = binary.LittleEndian.Uint16(data[offset : offset+2])
	}

	return style
}

// GetBinDataPath returns the path to binary data in the BinData storage.
func (info *BinDataInfo) GetBinDataPath() string {
	if info.BinDataID == 0 {
		return ""
	}
	return fmt.Sprintf("BIN%04X.%s", info.BinDataID, info.Extension)
}

// IsBold returns true if the character shape is bold.
func (cs *CharShape) IsBold() bool {
	return cs.Attributes&0x01 != 0
}

// IsItalic returns true if the character shape is italic.
func (cs *CharShape) IsItalic() bool {
	return cs.Attributes&0x02 != 0
}

// IsUnderline returns true if the character shape has underline.
func (cs *CharShape) IsUnderline() bool {
	return (cs.Attributes>>3)&0x07 != 0 // bits 3-5
}

// IsStrikeout returns true if the character shape has strikeout.
func (cs *CharShape) IsStrikeout() bool {
	return (cs.Attributes>>6)&0x07 != 0 // bits 6-8
}

// GetFontSizePt returns the font size in points.
func (cs *CharShape) GetFontSizePt() float64 {
	return float64(cs.Height) / 100.0
}
