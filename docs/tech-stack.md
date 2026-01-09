# hwp2markdown 기술 스택

## 1. 개요

hwp2markdown은 크로스 플랫폼 CLI 도구로, Windows, macOS, Linux에서 모두 동작해야 한다. 이 문서는 기술 스택 선정 과정과 결정 사항을 정리한다.

---

## 2. 언어 선택

### 후보 비교

| 언어 | 크로스 플랫폼 | 바이너리 배포 | 바이너리 크기 | 개발 속도 | 동시성 |
|------|--------------|---------------|---------------|-----------|--------|
| **Go** | O | 네이티브 바이너리 | 5-15MB | 빠름 | 고루틴 내장 |
| **Rust** | O | 네이티브 바이너리 | 2-10MB | 느림 | 비동기 지원 |
| **Python** | O | PyInstaller (30-50MB) | 큼 | 빠름 | GIL 제약 |
| **TypeScript** | O | pkg/nexe (50MB+) | 매우 큼 | 빠름 | 이벤트 루프 |

### 결정: Go

**이유:**
1. **단일 바이너리 배포**: 의존성 없이 단일 실행 파일로 배포 가능
2. **크로스 컴파일**: `GOOS`/`GOARCH` 설정만으로 모든 플랫폼 빌드
3. **빠른 빌드 속도**: 대규모 프로젝트도 수 초 내 빌드
4. **간결한 문법**: 학습 곡선이 낮고 코드 가독성 높음
5. **표준 라이브러리**: `archive/zip`, `encoding/xml` 등 풍부한 내장 기능
6. **동시성**: 고루틴으로 배치 처리 시 성능 최적화 용이

---

## 3. Go 버전

### 결정: Go 1.21+

**이유:**
- `slices`, `maps` 패키지 표준 라이브러리 포함
- 제네릭 기능 안정화
- `slog` 구조화된 로깅
- 향상된 PGO (Profile Guided Optimization)

---

## 4. 핵심 의존성

### 4.1 파일 파싱

| 패키지 | 용도 | 라이선스 |
|--------|------|----------|
| `archive/zip` (표준) | HWPX ZIP 압축 해제 | BSD |
| `encoding/xml` (표준) | HWPX XML 파싱 | BSD |
| `github.com/richardlehane/mscfb` | HWP 5.x OLE/CFBF 파싱 | Apache-2.0 |

#### HWPX 파싱 (표준 라이브러리)

```go
import (
    "archive/zip"
    "encoding/xml"
)

func parseHWPX(path string) (*Document, error) {
    r, err := zip.OpenReader(path)
    if err != nil {
        return nil, err
    }
    defer r.Close()

    for _, f := range r.File {
        if f.Name == "Contents/content.hpf" {
            rc, _ := f.Open()
            defer rc.Close()
            // XML 파싱
        }
    }
    return &Document{}, nil
}
```

#### HWP 5.x OLE 파싱

```go
import "github.com/richardlehane/mscfb"

func parseHWP5(path string) (*Document, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    doc, err := mscfb.New(file)
    if err != nil {
        return nil, err
    }

    for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
        // FileHeader, DocInfo, BodyText/Section0 등 처리
        fmt.Println(entry.Name)
    }
    return &Document{}, nil
}
```

### 4.2 CLI 프레임워크

| 패키지 | 용도 | 선택 이유 |
|--------|------|-----------|
| `github.com/spf13/cobra` | CLI 인터페이스 | 업계 표준, 서브커맨드 지원 |
| `github.com/spf13/pflag` | POSIX 플래그 | cobra 의존성 |

#### cobra 예시

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "hwp2markdown",
    Short: "HWP/HWPX 문서를 Markdown으로 변환",
    Long:  `HWP(한글 워드프로세서) 문서를 Markdown으로 변환하는 CLI 도구입니다.`,
}

var convertCmd = &cobra.Command{
    Use:   "convert [input]",
    Short: "HWP/HWPX 파일을 Markdown으로 변환",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        input := args[0]
        output, _ := cmd.Flags().GetString("output")
        return convert(input, output)
    },
}

func init() {
    convertCmd.Flags().StringP("output", "o", "", "출력 파일")
    convertCmd.Flags().String("extract-images", "", "이미지 추출 디렉토리")
    rootCmd.AddCommand(convertCmd)
}
```

### 4.3 기타 유틸리티

| 패키지 | 용도 | 라이선스 |
|--------|------|----------|
| `golang.org/x/text/encoding/korean` | EUC-KR 인코딩 처리 | BSD |
| `github.com/fatih/color` | 컬러 터미널 출력 (선택적) | MIT |

### 4.4 전체 의존성 (go.mod)

```go
module github.com/roboco-io/hwp2markdown

go 1.21

require (
    github.com/richardlehane/mscfb v1.0.4
    github.com/spf13/cobra v1.8.0
    golang.org/x/text v0.14.0
)
```

---

## 5. 바이너리 배포

### 5.1 크로스 컴파일

Go는 단일 명령으로 모든 플랫폼용 바이너리 생성 가능:

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o hwp2markdown-windows-x64.exe

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o hwp2markdown-macos-x64

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o hwp2markdown-macos-arm64

# Linux
GOOS=linux GOARCH=amd64 go build -o hwp2markdown-linux-x64
```

### 5.2 빌드 최적화

```bash
# 릴리스 빌드 (심볼 제거, 크기 최적화)
go build -ldflags="-s -w" -o hwp2markdown

# 버전 정보 삽입
go build -ldflags="-s -w -X main.version=1.0.0" -o hwp2markdown
```

### 5.3 배포 대상

| 플랫폼 | 아키텍처 | 파일명 | 예상 크기 |
|--------|----------|--------|-----------|
| Windows | x64 | `hwp2markdown-windows-x64.exe` | ~8MB |
| macOS | x64 | `hwp2markdown-macos-x64` | ~8MB |
| macOS | arm64 | `hwp2markdown-macos-arm64` | ~8MB |
| Linux | x64 | `hwp2markdown-linux-x64` | ~8MB |

---

## 6. 프로젝트 구조

```
hwp2markdown/
├── cmd/
│   └── hwp2markdown/
│       └── main.go              # 진입점
├── internal/
│   ├── cli/
│   │   ├── root.go              # cobra 루트 커맨드
│   │   └── convert.go           # convert 서브커맨드
│   ├── parser/
│   │   ├── parser.go            # 파서 인터페이스
│   │   ├── hwpx.go              # HWPX 파서
│   │   ├── hwp5.go              # HWP 5.x 파서
│   │   └── detector.go          # 포맷 감지
│   ├── model/
│   │   ├── document.go          # 문서 AST
│   │   ├── paragraph.go         # 문단
│   │   ├── table.go             # 표
│   │   └── image.go             # 이미지
│   └── renderer/
│       ├── renderer.go          # 렌더러 인터페이스
│       ├── markdown.go          # Markdown 렌더러
│       └── text.go              # Plain Text 렌더러
├── pkg/
│   └── hwp2markdown/
│       └── convert.go           # 공개 API
├── testdata/
│   ├── sample.hwpx
│   └── sample.hwp
├── docs/
│   ├── hwp-format-research.md
│   ├── existing-solutions-research.md
│   ├── PRD.md
│   └── tech-stack.md
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── LICENSE
└── .github/
    └── workflows/
        ├── test.yml
        └── release.yml
```

---

## 7. 빌드 및 개발

### 7.1 Makefile

```makefile
.PHONY: build test lint clean release

VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/hwp2markdown ./cmd/hwp2markdown

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/ dist/

release:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/hwp2markdown-windows-x64.exe ./cmd/hwp2markdown
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/hwp2markdown-macos-x64 ./cmd/hwp2markdown
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/hwp2markdown-macos-arm64 ./cmd/hwp2markdown
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/hwp2markdown-linux-x64 ./cmd/hwp2markdown
```

### 7.2 개발 환경 설정

```bash
# 저장소 클론
git clone https://github.com/roboco-io/hwp2markdown.git
cd hwp2markdown

# 의존성 다운로드
go mod download

# 빌드
make build

# 테스트
make test

# 린트 (golangci-lint 필요)
make lint
```

---

## 8. CI/CD

### 8.1 테스트 워크플로

```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ["1.21", "1.22"]

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Install dependencies
        run: go mod download

      - name: Lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

      - name: Test
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out
```

### 8.2 릴리스 워크플로

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 8.3 GoReleaser 설정

```yaml
# .goreleaser.yml
project_name: hwp2markdown

builds:
  - id: hwp2markdown
    main: ./cmd/hwp2markdown
    binary: hwp2markdown
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}

archives:
  - id: default
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"

release:
  github:
    owner: roboco-io
    name: hwp2markdown
```

---

## 9. 코드 품질

### 9.1 golangci-lint 설정

```yaml
# .golangci.yml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - revive

linters-settings:
  revive:
    rules:
      - name: exported
        disabled: false
```

### 9.2 pre-commit 설정

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.55.2
    hooks:
      - id: golangci-lint

  - repo: local
    hooks:
      - id: go-test
        name: go test
        entry: go test ./...
        language: system
        pass_filenames: false
```

---

## 10. 요약

| 항목 | 선택 |
|------|------|
| 언어 | Go 1.21+ |
| CLI 프레임워크 | cobra |
| HWPX 파싱 | archive/zip + encoding/xml (표준) |
| HWP 5.x 파싱 | mscfb |
| 한글 인코딩 | golang.org/x/text/encoding/korean |
| 린터 | golangci-lint |
| 테스트 | go test (표준) |
| 릴리스 | GoReleaser |
| CI/CD | GitHub Actions |

### Python 대비 Go의 장점

| 항목 | Python | Go |
|------|--------|-----|
| 바이너리 크기 | 30-50MB (PyInstaller) | 5-15MB |
| 의존성 | 런타임 필요 또는 번들링 | 없음 (단일 바이너리) |
| 시작 시간 | 느림 | 즉시 |
| 크로스 컴파일 | 복잡 | `GOOS`/`GOARCH`로 간단 |
| 동시성 | GIL 제약 | 고루틴 |
