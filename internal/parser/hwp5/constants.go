// Package hwp5 provides a parser for HWP 5.x binary documents.
package hwp5

// HWP 5.x 파일 포맷 상수 정의
// 참조: https://cdn.hancom.com/link/docs/한글문서파일형식_5.0_revision1.3.pdf

const (
	// FileHeader 시그니처
	Signature = "HWP Document File"

	// FileHeader 크기 (고정)
	FileHeaderSize = 256

	// 속성 플래그 비트
	FlagCompressed      uint32 = 1 << 0  // 압축 여부
	FlagEncrypted       uint32 = 1 << 1  // 암호화 여부
	FlagDistributable   uint32 = 1 << 2  // 배포용 문서
	FlagScript          uint32 = 1 << 3  // 스크립트 저장
	FlagDRM             uint32 = 1 << 4  // DRM 보안
	FlagXMLTemplate     uint32 = 1 << 5  // XMLTemplate 저장
	FlagHistory         uint32 = 1 << 6  // 문서 이력 관리
	FlagSignature       uint32 = 1 << 7  // 전자 서명
	FlagCertEncrypt     uint32 = 1 << 8  // 공인 인증서 암호화
	FlagSignatureReserv uint32 = 1 << 9  // 전자 서명 예비
	FlagCertDRM         uint32 = 1 << 10 // 공인 인증서 DRM
	FlagCCL             uint32 = 1 << 11 // CCL 문서
	FlagMobile          uint32 = 1 << 12 // 모바일 최적화
)

// 스트림 이름
const (
	StreamFileHeader  = "FileHeader"
	StreamDocInfo     = "DocInfo"
	StreamBodyText    = "BodyText"
	StreamViewText    = "ViewText"
	StreamSummaryInfo = "\x05HwpSummaryInformation"
	StreamBinData     = "BinData"
	StreamPrvText     = "PrvText"
	StreamPrvImage    = "PrvImage"
	StreamDocOptions  = "DocOptions"
	StreamScripts     = "Scripts"
	StreamXMLTemplate = "XMLTemplate"
	StreamDocHistory  = "DocHistory"
)

// 레코드 태그 ID (HWPTAG_*)
// 참조: HWP 5.0 명세서 4장 레코드 구조
const (
	// DocInfo 레코드 태그 (0x0010 ~ 0x003F)
	TagDocumentProperties uint16 = 0x0010 // 문서 속성
	TagIDMappings         uint16 = 0x0011 // ID 매핑 테이블 크기
	TagBinData            uint16 = 0x0012 // 바이너리 데이터
	TagFaceName           uint16 = 0x0013 // 글꼴
	TagBorderFill         uint16 = 0x0014 // 테두리/배경
	TagCharShape          uint16 = 0x0015 // 글자 모양
	TagTabDef             uint16 = 0x0016 // 탭 정의
	TagNumbering          uint16 = 0x0017 // 문단 번호
	TagBullet             uint16 = 0x0018 // 글머리표
	TagParaShape          uint16 = 0x0019 // 문단 모양
	TagStyle              uint16 = 0x001A // 스타일
	TagDocData            uint16 = 0x001B // 문서 데이터
	TagDistributeDocData  uint16 = 0x001C // 배포용 문서 데이터
	TagCompatibleDocument uint16 = 0x001E // 호환 문서
	TagLayoutCompatible   uint16 = 0x001F // 레이아웃 호환
	TagTrackChange        uint16 = 0x0020 // 변경 내역 추적
	TagMemoShape          uint16 = 0x0022 // 메모 모양
	TagForbiddenChar      uint16 = 0x0023 // 금칙 문자
	TagTrackChange2       uint16 = 0x0024 // 변경 내역 2
	TagTrackChangeAuthor  uint16 = 0x0025 // 변경 내역 작성자

	// Section/Body 레코드 태그 (0x0040 ~ 0x007F)
	TagParaHeader     uint16 = 0x0042 // 문단 헤더
	TagParaText       uint16 = 0x0043 // 문단 텍스트
	TagParaCharShape  uint16 = 0x0044 // 문단 글자 모양
	TagParaLineSeg    uint16 = 0x0045 // 문단 레이아웃
	TagParaRangeTag   uint16 = 0x0046 // 문단 범위 태그
	TagCtrlHeader     uint16 = 0x0047 // 컨트롤 헤더
	TagListHeader     uint16 = 0x0048 // 리스트 헤더
	TagPageDef        uint16 = 0x0049 // 페이지 정의
	TagFootnoteShape  uint16 = 0x004A // 각주 모양
	TagPageBorderFill uint16 = 0x004B // 페이지 테두리/배경
	TagShapeComponent uint16 = 0x004C // 그리기 개체
	TagTable          uint16 = 0x004D // 표
	TagShapeLine      uint16 = 0x004E // 선
	TagShapeRectangle uint16 = 0x004F // 사각형
	TagShapeEllipse   uint16 = 0x0050 // 타원
	TagShapeArc       uint16 = 0x0051 // 호
	TagShapePolygon   uint16 = 0x0052 // 다각형
	TagShapeCurve     uint16 = 0x0053 // 곡선
	TagShapeOLE       uint16 = 0x0054 // OLE
	TagShapePicture   uint16 = 0x0055 // 그림
	TagShapeContainer uint16 = 0x0056 // 컨테이너
	TagCtrlData       uint16 = 0x0057 // 컨트롤 데이터
	TagEqEdit         uint16 = 0x0058 // 수식
	TagCtrlFormField  uint16 = 0x005B // 양식 컨트롤
	TagMemoList       uint16 = 0x005C // 메모 리스트
	TagCellListHeader uint16 = 0x005E // 셀 리스트 헤더
	TagChartData      uint16 = 0x0060 // 차트 데이터
	TagVideoData      uint16 = 0x0062 // 비디오 데이터
)

// 컨트롤 타입 ID (4바이트 문자열)
const (
	CtrlSection       = "secd" // 구역 정의
	CtrlColumn        = "cold" // 단 정의
	CtrlHeader        = "head" // 머리말
	CtrlFooter        = "foot" // 꼬리말
	CtrlFootnote      = "fn  " // 각주
	CtrlEndnote       = "en  " // 미주
	CtrlAutoNumber    = "atno" // 자동 번호
	CtrlNewNumber     = "nwno" // 새 번호
	CtrlPageHide      = "pghd" // 페이지 숨김
	CtrlPageOddEven   = "pgct" // 홀/짝수 페이지
	CtrlPageNumber    = "pgno" // 페이지 번호
	CtrlIndexMark     = "idxm" // 찾아보기 표식
	CtrlBookmark      = "bokm" // 책갈피
	CtrlOverlapping   = "tcps" // 글자 겹침
	CtrlHiddenComment = "tdut" // 숨은 설명
	CtrlTable         = "tbl " // 표
	CtrlGSO           = "gso " // 그리기 개체
	CtrlEquation      = "eqed" // 수식
	CtrlFieldBegin    = "%beg" // 필드 시작
	CtrlFieldEnd      = "%end" // 필드 끝
)

// 특수 문자 코드
const (
	CharLine           = 0x0000 // 줄 나눔
	CharPara           = 0x000D // 문단 나눔
	CharTab            = 0x0009 // 탭
	CharDrawingObj     = 0x000B // 그리기 개체/표
	CharInlineStart    = 0x000C // 인라인 시작
	CharFieldStart     = 0x0003 // 필드 시작
	CharFieldEnd       = 0x0004 // 필드 끝
	CharBookmark       = 0x0005 // 책갈피
	CharTitleMark      = 0x0006 // 제목 표시
	CharHyphen         = 0x001E // 하이픈
	CharNBSP           = 0x001F // 줄바꿈 방지 공백
	CharFixedWidthNBSP = 0x0018 // 고정폭 빈칸
	CharExtChar        = 0x0014 // 확장 문자
)
