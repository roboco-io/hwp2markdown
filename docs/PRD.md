# hwp2markdown PRD (Product Requirements Document)

## 1. 개요

### 1.1 프로젝트 명
hwp2markdown

### 1.2 목적
HWP(한글 워드프로세서) 문서를 Markdown으로 변환하는 오픈소스 도구 개발

### 1.3 배경
- HWP는 한국에서 널리 사용되는 문서 포맷이나, 버전 간 호환성 문제와 폐쇄적 생태계로 인해 활용에 제약이 있음
- LLM/AI 시대에 문서를 Markdown으로 변환하여 처리하려는 수요 증가
- 기존 솔루션들의 한계:
  - `unhwp` (Rust): 유망하나 Rust 의존성 필요
  - `hwpjs`: JSON까지만 파싱, Markdown 변환 미지원
  - 상용 서비스: 외부 의존성, 비용, 프라이버시 우려

### 1.4 목표
- HWPX 및 HWP 5.x 포맷의 Markdown 변환 지원
- 텍스트, 표, 이미지, 스타일 등 주요 요소의 정확한 변환
- CLI 및 라이브러리 형태로 제공
- MIT 라이선스 오픈소스

---

## 2. 사용자 및 사용 사례

### 2.1 대상 사용자

| 사용자 유형 | 설명 |
|-------------|------|
| 개발자 | HWP 문서를 프로그래밍 방식으로 처리해야 하는 개발자 |
| 데이터 엔지니어 | HWP 문서에서 텍스트를 추출하여 데이터 파이프라인에 활용 |
| AI/ML 엔지니어 | LLM 학습/추론을 위해 HWP 문서를 Markdown으로 변환 |
| 일반 사용자 | CLI를 통해 HWP 문서를 Markdown으로 변환하려는 사용자 |

### 2.2 주요 사용 사례

#### UC-1: CLI를 통한 단일 파일 변환
```bash
hwp2markdown input.hwp -o output.md
```

#### UC-2: CLI를 통한 배치 변환
```bash
hwp2markdown ./documents/*.hwp -o ./output/
```

#### UC-3: 라이브러리로 프로그래밍 방식 사용
```python
from hwp2markdown import convert

markdown = convert("document.hwp")
print(markdown)
```

#### UC-4: 이미지 추출과 함께 변환
```bash
hwp2markdown input.hwp -o output.md --extract-images ./images/
```

---

## 3. 기능 요구사항

### 3.1 지원 포맷

| 우선순위 | 포맷 | 설명 |
|----------|------|------|
| P0 | HWPX | XML 기반 개방형 포맷 (ZIP + XML) |
| P1 | HWP 5.x | OLE/Compound 바이너리 포맷 |
| P2 | HWP 3.x | 레거시 바이너리 포맷 (향후 검토) |

### 3.2 변환 요소

#### P0 (필수)

| 요소 | Markdown 변환 |
|------|---------------|
| 일반 텍스트 | 그대로 출력 |
| 문단 | 빈 줄로 구분 |
| 제목/개요 | `#`, `##`, `###` 등 |
| 굵게 | `**텍스트**` |
| 기울임 | `*텍스트*` |
| 취소선 | `~~텍스트~~` |
| 순서 없는 목록 | `- 항목` |
| 순서 있는 목록 | `1. 항목` |
| 표 | GFM 테이블 문법 |
| 하이퍼링크 | `[텍스트](URL)` |

#### P1 (중요)

| 요소 | Markdown 변환 |
|------|---------------|
| 이미지 | `![alt](path)` + 파일 추출 |
| 코드 블록 | ``` 문법 |
| 인용문 | `>` |
| 각주 | `[^1]` 문법 |
| 수평선 | `---` |

#### P2 (향후)

| 요소 | Markdown 변환 |
|------|---------------|
| 수식 | LaTeX 문법 (`$...$`) |
| 차트 | 이미지로 추출 또는 무시 |
| 양식 필드 | 텍스트로 변환 또는 무시 |
| 머리글/바닥글 | 별도 섹션 또는 무시 |

### 3.3 CLI 인터페이스

```
hwp2markdown [OPTIONS] <INPUT>...

Arguments:
  <INPUT>...  입력 HWP/HWPX 파일 또는 디렉토리

Options:
  -o, --output <PATH>       출력 파일 또는 디렉토리
  -f, --format <FORMAT>     출력 포맷 [기본값: markdown]
                            가능한 값: markdown, text, json
  --extract-images <DIR>    이미지 추출 디렉토리
  --image-format <FORMAT>   이미지 참조 형식 [기본값: relative]
                            가능한 값: relative, absolute, base64
  --table-format <FORMAT>   테이블 형식 [기본값: gfm]
                            가능한 값: gfm, html
  --heading-style <STYLE>   제목 스타일 [기본값: atx]
                            가능한 값: atx, setext
  -v, --verbose             상세 출력
  -q, --quiet               조용한 모드
  -h, --help                도움말 출력
  -V, --version             버전 출력
```

### 3.4 라이브러리 API

#### Python API

```python
from hwp2markdown import HwpConverter, ConvertOptions

# 간단한 사용
markdown = hwp2markdown.convert("document.hwp")

# 옵션과 함께 사용
options = ConvertOptions(
    extract_images=True,
    image_dir="./images",
    table_format="gfm",
)
converter = HwpConverter(options)
result = converter.convert("document.hwp")

print(result.markdown)
print(result.images)  # 추출된 이미지 목록
print(result.metadata)  # 문서 메타데이터
```

#### 결과 객체

```python
@dataclass
class ConvertResult:
    markdown: str              # 변환된 Markdown
    images: list[ImageInfo]    # 추출된 이미지 정보
    metadata: DocumentMetadata # 문서 메타데이터
    warnings: list[str]        # 변환 중 경고
```

---

## 4. 비기능 요구사항

### 4.1 성능

| 항목 | 목표 |
|------|------|
| 변환 속도 | 10MB 문서 기준 10초 이내 |
| 메모리 사용 | 입력 파일 크기의 10배 이내 |
| 동시 처리 | 배치 변환 시 멀티프로세싱 지원 |

### 4.2 품질

| 항목 | 목표 |
|------|------|
| 텍스트 정확도 | 99% 이상 (글자 손실 없음) |
| 구조 보존율 | 95% 이상 (제목, 목록, 표 구조) |
| 테스트 커버리지 | 80% 이상 |

### 4.3 호환성

| 항목 | 요구사항 |
|------|----------|
| Python | 3.10 이상 |
| OS | Linux, macOS, Windows |
| 인코딩 | UTF-8 출력 |

### 4.4 배포

| 항목 | 요구사항 |
|------|----------|
| 패키지 | PyPI 배포 (`pip install hwp2markdown`) |
| 바이너리 | 주요 OS용 standalone 바이너리 제공 (향후) |
| Docker | Docker 이미지 제공 (향후) |

---

## 5. 기술 설계

### 5.1 아키텍처

```
┌─────────────────────────────────────────────────────────┐
│                      hwp2markdown                       │
├─────────────────────────────────────────────────────────┤
│  CLI Layer                                              │
│  ┌─────────────────────────────────────────────────┐   │
│  │  click/typer 기반 CLI                            │   │
│  └─────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────┤
│  Core API                                               │
│  ┌─────────────────────────────────────────────────┐   │
│  │  HwpConverter                                    │   │
│  │  - convert(path) -> ConvertResult               │   │
│  └─────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────┤
│  Parser Layer                                           │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  HwpxParser      │  │  Hwp5Parser      │            │
│  │  (ZIP + XML)     │  │  (OLE/CFBF)      │            │
│  └──────────────────┘  └──────────────────┘            │
├─────────────────────────────────────────────────────────┤
│  Document Model (AST)                                   │
│  ┌─────────────────────────────────────────────────┐   │
│  │  Document, Section, Paragraph, Run, Table, ...  │   │
│  └─────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────┤
│  Renderer Layer                                         │
│  ┌──────────────────┐  ┌──────────────────┐            │
│  │  MarkdownRenderer│  │  TextRenderer    │            │
│  └──────────────────┘  └──────────────────┘            │
└─────────────────────────────────────────────────────────┘
```

### 5.2 핵심 모듈

| 모듈 | 역할 |
|------|------|
| `hwp2markdown.parser.hwpx` | HWPX 파일 파싱 (ZIP + XML) |
| `hwp2markdown.parser.hwp5` | HWP 5.x 파일 파싱 (OLE/CFBF) |
| `hwp2markdown.model` | 문서 AST 정의 |
| `hwp2markdown.renderer.markdown` | AST → Markdown 변환 |
| `hwp2markdown.renderer.text` | AST → Plain Text 변환 |
| `hwp2markdown.cli` | CLI 인터페이스 |

### 5.3 의존성

| 패키지 | 용도 |
|--------|------|
| `olefile` | OLE/Compound 파일 파싱 (HWP 5.x) |
| `lxml` | XML 파싱 (HWPX) |
| `click` 또는 `typer` | CLI 인터페이스 |
| `Pillow` | 이미지 처리 (선택적) |

---

## 6. 구현 단계

### Phase 1: HWPX 기본 지원

**목표**: HWPX 파일의 텍스트, 문단, 기본 스타일 변환

| 작업 | 설명 |
|------|------|
| 프로젝트 구조 설정 | Python 패키지 구조, 테스트 환경 |
| HWPX 파서 구현 | ZIP 압축 해제, XML 파싱 |
| 문서 모델 정의 | Document, Section, Paragraph, Run |
| 기본 Markdown 렌더러 | 텍스트, 제목, 굵게, 기울임, 목록 |
| CLI 기본 구현 | 단일 파일 변환 |

### Phase 2: HWPX 고급 기능

**목표**: 표, 이미지, 각주 등 고급 요소 지원

| 작업 | 설명 |
|------|------|
| 표 파싱 및 변환 | GFM 테이블 문법 |
| 이미지 추출 | BinData에서 이미지 추출 |
| 각주/미주 처리 | Markdown 각주 문법 |
| 메타데이터 추출 | 제목, 작성자, 날짜 등 |

### Phase 3: HWP 5.x 지원

**목표**: HWP 5.x 바이너리 포맷 파싱 및 변환

| 작업 | 설명 |
|------|------|
| OLE 파일 파싱 | olefile 라이브러리 활용 |
| 바이너리 레코드 파싱 | HWP 5.x 레코드 구조 해석 |
| 기존 렌더러 재사용 | 동일 AST → Markdown 변환 |

### Phase 4: 배포 및 안정화

**목표**: PyPI 배포, 문서화, 테스트 강화

| 작업 | 설명 |
|------|------|
| PyPI 배포 | 패키지 빌드 및 배포 |
| 문서화 | README, API 문서, 예제 |
| 테스트 강화 | 다양한 HWP 샘플로 테스트 |
| CI/CD 설정 | GitHub Actions |

---

## 7. 성공 지표

| 지표 | 목표 |
|------|------|
| 지원 포맷 | HWPX, HWP 5.x |
| 변환 정확도 | 텍스트 99%, 구조 95% |
| PyPI 다운로드 | 출시 후 3개월 내 1,000회 |
| GitHub Stars | 출시 후 6개월 내 100개 |
| 이슈 대응 | 평균 7일 이내 응답 |

---

## 8. 위험 및 대응

| 위험 | 영향 | 대응 |
|------|------|------|
| HWP 5.x 바이너리 구조 복잡성 | 개발 지연 | HWPX 우선 지원, 기존 파서 참고 |
| 다양한 HWP 버전 호환성 | 변환 실패 | 점진적 버전 지원, 사용자 피드백 수집 |
| 복잡한 레이아웃 변환 한계 | 품질 저하 | Markdown 한계 명시, 대안 포맷 제공 |
| 한글과컴퓨터 라이선스 이슈 | 법적 문제 | 리버스 엔지니어링 법적 검토, 공개 스펙 활용 |

---

## 9. 참고 자료

- [HWP 포맷 조사 보고서](hwp-format-research.md)
- [기존 솔루션 조사](existing-solutions-research.md)
- [unhwp (Rust)](https://lib.rs/crates/unhwp)
- [hwpjs (JavaScript)](https://github.com/ohah/hwpjs)
- [olefile (Python OLE parser)](https://olefile.readthedocs.io/)
