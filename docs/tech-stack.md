# hwp2markdown 기술 스택

## 1. 개요

hwp2markdown은 크로스 플랫폼 CLI 도구로, Windows, macOS, Linux에서 모두 동작해야 한다. 이 문서는 기술 스택 선정 과정과 결정 사항을 정리한다.

---

## 2. 언어 선택

### 후보 비교

| 언어 | 크로스 플랫폼 | 바이너리 배포 | 생태계 | 개발 속도 |
|------|--------------|---------------|--------|-----------|
| **Rust** | O | 네이티브 바이너리 | OLE 라이브러리 부족 | 느림 |
| **Go** | O | 네이티브 바이너리 | OLE 라이브러리 부족 | 보통 |
| **Python** | O | PyInstaller/Nuitka | olefile, lxml 등 풍부 | 빠름 |
| **TypeScript** | O | pkg/nexe | 브라우저 중심 생태계 | 빠름 |

### 결정: Python

**이유:**
1. **OLE 파싱 라이브러리**: `olefile`이 성숙하고 HWP 5.x 파싱에 적합
2. **XML 파싱**: `lxml`이 고성능 XML 파싱 제공
3. **개발 생산성**: 빠른 프로토타이핑과 이터레이션
4. **바이너리 배포**: PyInstaller 또는 Nuitka로 standalone 바이너리 생성 가능
5. **커뮤니티**: Python 사용자층이 넓어 기여 및 유지보수 용이

---

## 3. Python 버전

### 결정: Python 3.10+

**이유:**
- `match` 문법 (구조적 패턴 매칭) 활용
- `dataclasses` 개선사항
- 타입 힌트 기능 강화 (`ParamSpec`, `TypeAlias`)
- 2026년까지 보안 지원

---

## 4. 핵심 의존성

### 4.1 파일 파싱

| 패키지 | 버전 | 용도 | 라이선스 |
|--------|------|------|----------|
| `olefile` | >=0.46 | HWP 5.x OLE/Compound 파일 파싱 | BSD-2-Clause |
| `lxml` | >=4.9 | HWPX XML 파싱 | BSD-3-Clause |

#### olefile

```python
import olefile

ole = olefile.OleFileIO("document.hwp")
streams = ole.listdir()
# ['FileHeader', 'DocInfo', 'BodyText/Section0', ...]
data = ole.openstream("BodyText/Section0").read()
```

#### lxml

```python
from lxml import etree

tree = etree.parse("Contents/section0.xml")
root = tree.getroot()
paragraphs = root.findall(".//hp:p", namespaces={"hp": "..."})
```

### 4.2 CLI 프레임워크

| 패키지 | 용도 | 선택 이유 |
|--------|------|-----------|
| `typer` | CLI 인터페이스 | 타입 힌트 기반, 자동 완성, 도움말 생성 |
| `rich` | 터미널 출력 | 컬러, 프로그레스 바, 테이블 (typer 의존성) |

#### typer 예시

```python
import typer
from pathlib import Path

app = typer.Typer()

@app.command()
def convert(
    input_file: Path = typer.Argument(..., help="입력 HWP/HWPX 파일"),
    output: Path = typer.Option(None, "-o", "--output", help="출력 파일"),
    extract_images: Path = typer.Option(None, "--extract-images", help="이미지 추출 디렉토리"),
):
    """HWP/HWPX 파일을 Markdown으로 변환합니다."""
    ...

if __name__ == "__main__":
    app()
```

### 4.3 이미지 처리 (선택적)

| 패키지 | 버전 | 용도 | 라이선스 |
|--------|------|------|----------|
| `Pillow` | >=9.0 | 이미지 포맷 변환, 메타데이터 | HPND |

### 4.4 전체 의존성

```toml
# pyproject.toml
[project]
dependencies = [
    "olefile>=0.46",
    "lxml>=4.9",
    "typer>=0.9",
    "rich>=13.0",
]

[project.optional-dependencies]
images = ["Pillow>=9.0"]
dev = [
    "pytest>=7.0",
    "pytest-cov>=4.0",
    "ruff>=0.1",
    "mypy>=1.0",
]
```

---

## 5. 바이너리 배포

### 5.1 도구 비교

| 도구 | 바이너리 크기 | 시작 시간 | 난이도 |
|------|--------------|-----------|--------|
| **PyInstaller** | 30-50MB | 빠름 | 낮음 |
| **Nuitka** | 10-30MB | 매우 빠름 | 중간 |
| **cx_Freeze** | 30-50MB | 빠름 | 낮음 |

### 5.2 결정: PyInstaller (1차) → Nuitka (향후)

**1차 배포: PyInstaller**
- 설정 간단, 빠른 빌드
- 크로스 플랫폼 지원 우수

**향후: Nuitka**
- 더 작은 바이너리
- 더 빠른 실행 속도
- Python 코드를 C로 컴파일

### 5.3 빌드 설정

```python
# hwp2markdown.spec (PyInstaller)
a = Analysis(
    ['src/hwp2markdown/__main__.py'],
    pathex=[],
    binaries=[],
    datas=[],
    hiddenimports=['lxml._elementpath'],
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=[],
    noarchive=False,
)

pyz = PYZ(a.pure)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.datas,
    [],
    name='hwp2markdown',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    console=True,
)
```

### 5.4 배포 대상

| 플랫폼 | 아키텍처 | 파일명 |
|--------|----------|--------|
| Windows | x64 | `hwp2markdown-windows-x64.exe` |
| macOS | x64 | `hwp2markdown-macos-x64` |
| macOS | arm64 | `hwp2markdown-macos-arm64` |
| Linux | x64 | `hwp2markdown-linux-x64` |

---

## 6. 프로젝트 구조

```
hwp2markdown/
├── src/
│   └── hwp2markdown/
│       ├── __init__.py
│       ├── __main__.py          # CLI 진입점
│       ├── cli.py               # typer CLI 정의
│       ├── converter.py         # 메인 변환 로직
│       ├── parser/
│       │   ├── __init__.py
│       │   ├── base.py          # 파서 추상 클래스
│       │   ├── hwpx.py          # HWPX 파서
│       │   └── hwp5.py          # HWP 5.x 파서
│       ├── model/
│       │   ├── __init__.py
│       │   └── document.py      # 문서 AST 모델
│       └── renderer/
│           ├── __init__.py
│           ├── base.py          # 렌더러 추상 클래스
│           ├── markdown.py      # Markdown 렌더러
│           └── text.py          # Plain Text 렌더러
├── tests/
│   ├── __init__.py
│   ├── conftest.py
│   ├── fixtures/                # 테스트용 HWP/HWPX 파일
│   ├── test_parser_hwpx.py
│   ├── test_parser_hwp5.py
│   ├── test_renderer_markdown.py
│   └── test_cli.py
├── docs/
│   ├── hwp-format-research.md
│   ├── existing-solutions-research.md
│   ├── PRD.md
│   └── tech-stack.md
├── pyproject.toml
├── README.md
├── LICENSE
└── .github/
    └── workflows/
        ├── test.yml             # 테스트 CI
        └── release.yml          # 바이너리 빌드 및 릴리스
```

---

## 7. 빌드 및 패키징

### 7.1 pyproject.toml

```toml
[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[project]
name = "hwp2markdown"
version = "0.1.0"
description = "HWP(한글 워드프로세서) 문서를 Markdown으로 변환하는 도구"
readme = "README.md"
license = "MIT"
requires-python = ">=3.10"
authors = [
    { name = "roboco-io" }
]
keywords = ["hwp", "hwpx", "markdown", "converter", "hangul"]
classifiers = [
    "Development Status :: 3 - Alpha",
    "Environment :: Console",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Operating System :: OS Independent",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
    "Programming Language :: Python :: 3.12",
    "Topic :: Text Processing :: Markup",
]
dependencies = [
    "olefile>=0.46",
    "lxml>=4.9",
    "typer>=0.9",
    "rich>=13.0",
]

[project.optional-dependencies]
images = ["Pillow>=9.0"]
dev = [
    "pytest>=7.0",
    "pytest-cov>=4.0",
    "ruff>=0.1",
    "mypy>=1.0",
    "pyinstaller>=6.0",
]

[project.scripts]
hwp2markdown = "hwp2markdown.cli:app"

[project.urls]
Homepage = "https://github.com/roboco-io/hwp2markdown"
Repository = "https://github.com/roboco-io/hwp2markdown"
Issues = "https://github.com/roboco-io/hwp2markdown/issues"

[tool.hatch.build.targets.wheel]
packages = ["src/hwp2markdown"]

[tool.ruff]
line-length = 100
target-version = "py310"

[tool.ruff.lint]
select = ["E", "F", "I", "N", "W", "UP"]

[tool.mypy]
python_version = "3.10"
strict = true

[tool.pytest.ini_options]
testpaths = ["tests"]
addopts = "-v --cov=hwp2markdown --cov-report=term-missing"
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
        python-version: ["3.10", "3.11", "3.12"]

    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: ${{ matrix.python-version }}

      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install -e ".[dev]"

      - name: Lint with ruff
        run: ruff check src tests

      - name: Type check with mypy
        run: mypy src

      - name: Test with pytest
        run: pytest
```

### 8.2 릴리스 워크플로

```yaml
# .github/workflows/release.yml
name: Release

on:
  release:
    types: [created]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            asset_name: hwp2markdown-linux-x64
          - os: macos-latest
            asset_name: hwp2markdown-macos-x64
          - os: macos-14
            asset_name: hwp2markdown-macos-arm64
          - os: windows-latest
            asset_name: hwp2markdown-windows-x64.exe

    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: "3.11"

      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install -e ".[dev]"

      - name: Build with PyInstaller
        run: pyinstaller --onefile --name ${{ matrix.asset_name }} src/hwp2markdown/__main__.py

      - name: Upload Release Asset
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ github.event.release.upload_url }}
          asset_path: ./dist/${{ matrix.asset_name }}
          asset_name: ${{ matrix.asset_name }}
          asset_content_type: application/octet-stream

  publish-pypi:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: "3.11"

      - name: Install build tools
        run: pip install build twine

      - name: Build package
        run: python -m build

      - name: Publish to PyPI
        env:
          TWINE_USERNAME: __token__
          TWINE_PASSWORD: ${{ secrets.PYPI_TOKEN }}
        run: twine upload dist/*
```

---

## 9. 개발 환경 설정

### 9.1 초기 설정

```bash
# 저장소 클론
git clone https://github.com/roboco-io/hwp2markdown.git
cd hwp2markdown

# 가상환경 생성
python -m venv .venv
source .venv/bin/activate  # Windows: .venv\Scripts\activate

# 개발 의존성 설치
pip install -e ".[dev,images]"

# 린트 및 테스트
ruff check src tests
mypy src
pytest
```

### 9.2 로컬 바이너리 빌드

```bash
# PyInstaller로 빌드
pyinstaller --onefile --name hwp2markdown src/hwp2markdown/__main__.py

# 실행
./dist/hwp2markdown --help
```

---

## 10. 요약

| 항목 | 선택 |
|------|------|
| 언어 | Python 3.10+ |
| CLI 프레임워크 | typer + rich |
| OLE 파싱 | olefile |
| XML 파싱 | lxml |
| 이미지 처리 | Pillow (선택적) |
| 패키징 | hatchling |
| 바이너리 빌드 | PyInstaller |
| 린터 | ruff |
| 타입 체커 | mypy |
| 테스트 | pytest |
| CI/CD | GitHub Actions |
