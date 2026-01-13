package hwp5

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// FileHeader는 HWP 5.x 파일 헤더 구조체
// 참조: HWP 5.0 명세서 2.1 파일 인식 정보
type FileHeader struct {
	Signature   [32]byte  // 파일 시그니처 "HWP Document File"
	Version     Version   // 파일 버전
	Flags       uint32    // 속성 플래그
	LicenseInfo [216]byte // 예약 영역
}

// Version은 HWP 파일 버전 (예: 5.0.3.0)
type Version struct {
	Major    uint8
	Minor    uint8
	Build    uint8
	Revision uint8
}

// String returns version string like "5.0.3.0"
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Build, v.Revision)
}

// ParseFileHeader parses the FileHeader from raw bytes.
func ParseFileHeader(data []byte) (*FileHeader, error) {
	if len(data) < FileHeaderSize {
		return nil, fmt.Errorf("file header too small: %d bytes", len(data))
	}

	h := &FileHeader{}

	// 시그니처 (32 bytes)
	copy(h.Signature[:], data[0:32])

	// 시그니처 검증
	sigStr := string(bytes.TrimRight(h.Signature[:], "\x00"))
	if sigStr != Signature {
		return nil, fmt.Errorf("invalid HWP signature: %q", sigStr)
	}

	// 버전 (4 bytes, little-endian)
	// 포맷: [Revision][Build][Minor][Major]
	h.Version.Revision = data[32]
	h.Version.Build = data[33]
	h.Version.Minor = data[34]
	h.Version.Major = data[35]

	// 속성 플래그 (4 bytes, little-endian)
	h.Flags = binary.LittleEndian.Uint32(data[36:40])

	// 나머지는 예약 영역
	copy(h.LicenseInfo[:], data[40:256])

	return h, nil
}

// IsCompressed returns true if the document is compressed.
func (h *FileHeader) IsCompressed() bool {
	return h.Flags&FlagCompressed != 0
}

// IsEncrypted returns true if the document is encrypted.
func (h *FileHeader) IsEncrypted() bool {
	return h.Flags&FlagEncrypted != 0
}

// IsDistributable returns true if this is a distribution document.
func (h *FileHeader) IsDistributable() bool {
	return h.Flags&FlagDistributable != 0
}

// HasScript returns true if the document contains scripts.
func (h *FileHeader) HasScript() bool {
	return h.Flags&FlagScript != 0
}

// HasDRM returns true if the document has DRM protection.
func (h *FileHeader) HasDRM() bool {
	return h.Flags&FlagDRM != 0
}

// HasXMLTemplate returns true if the document contains XML templates.
func (h *FileHeader) HasXMLTemplate() bool {
	return h.Flags&FlagXMLTemplate != 0
}

// HasHistory returns true if the document has revision history.
func (h *FileHeader) HasHistory() bool {
	return h.Flags&FlagHistory != 0
}

// HasSignature returns true if the document has a digital signature.
func (h *FileHeader) HasSignature() bool {
	return h.Flags&FlagSignature != 0
}

// IsCertEncrypted returns true if encrypted with certificate.
func (h *FileHeader) IsCertEncrypted() bool {
	return h.Flags&FlagCertEncrypt != 0
}

// IsCCL returns true if this is a CCL document.
func (h *FileHeader) IsCCL() bool {
	return h.Flags&FlagCCL != 0
}

// IsMobileOptimized returns true if optimized for mobile.
func (h *FileHeader) IsMobileOptimized() bool {
	return h.Flags&FlagMobile != 0
}
