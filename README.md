# hwp2markdown

HWP(한글 워드프로세서) 문서를 Markdown으로 변환하는 CLI 도구

## 개요

이 프로젝트는 다양한 HWP 버전의 문서를 Markdown으로 변환하여 범용적으로 활용할 수 있게 하는 것을 목표로 합니다.

## 설치

### 바이너리 다운로드

[GitHub Releases](https://github.com/roboco-io/hwp2markdown/releases)에서 플랫폼별 바이너리를 다운로드할 수 있습니다.

| 플랫폼 | 아키텍처 | 파일 |
|--------|----------|------|
| Windows | x64 | `hwp2markdown_*_windows_amd64.zip` |
| macOS | Intel | `hwp2markdown_*_darwin_amd64.tar.gz` |
| macOS | Apple Silicon | `hwp2markdown_*_darwin_arm64.tar.gz` |
| Linux | x64 | `hwp2markdown_*_linux_amd64.tar.gz` |

### Go 설치

```bash
go install github.com/roboco-io/hwp2markdown/cmd/hwp2markdown@latest
```

## 사용법

```bash
# 기본 사용
hwp2markdown convert document.hwpx -o output.md

# 이미지 추출과 함께
hwp2markdown convert document.hwp -o output.md --extract-images ./images/

# 도움말
hwp2markdown --help
hwp2markdown convert --help
```

## 지원 포맷

| 포맷 | 상태 | 설명 |
|------|------|------|
| HWPX | 개발 중 | XML 기반 개방형 포맷 |
| HWP 5.x | 계획 | OLE/Compound 바이너리 포맷 |
| HWP 3.x | 미정 | 레거시 바이너리 포맷 |

## 개발

### 요구사항

- Go 1.21 이상

### 빌드

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

### 크로스 컴파일

```bash
# 모든 플랫폼용 바이너리 빌드
make release
```

## 문서

- [HWP 포맷 조사 보고서](docs/hwp-format-research.md)
- [기존 솔루션 조사](docs/existing-solutions-research.md)
- [PRD](docs/PRD.md)
- [기술 스택](docs/tech-stack.md)

## 라이센스

MIT License - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.
