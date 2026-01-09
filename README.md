# hwp2markdown

HWP(한글 워드프로세서) 문서를 Markdown으로 변환하는 도구

## 개요

이 프로젝트는 다양한 HWP 버전의 문서를 Markdown으로 변환하여 범용적으로 활용할 수 있게 하는 것을 목표로 합니다.

## 설치

```bash
pip install hwp2markdown
```

또는 [GitHub Releases](https://github.com/roboco-io/hwp2markdown/releases)에서 플랫폼별 바이너리를 다운로드할 수 있습니다.

## 사용법

```bash
# 기본 사용
hwp2markdown input.hwpx -o output.md

# 이미지 추출과 함께
hwp2markdown input.hwp -o output.md --extract-images ./images/

# 도움말
hwp2markdown --help
```

## 지원 포맷

| 포맷 | 상태 | 설명 |
|------|------|------|
| HWPX | 계획 | XML 기반 개방형 포맷 |
| HWP 5.x | 계획 | OLE/Compound 바이너리 포맷 |
| HWP 3.x | 미정 | 레거시 바이너리 포맷 |

## 개발

```bash
# 저장소 클론
git clone https://github.com/roboco-io/hwp2markdown.git
cd hwp2markdown

# 개발 의존성 설치
pip install -e ".[dev]"

# Git hooks 설정
./scripts/setup-hooks.sh

# 테스트 실행
pytest

# 린트 체크
ruff check src tests
mypy src
```

## 문서

- [HWP 포맷 조사 보고서](docs/hwp-format-research.md)
- [기존 솔루션 조사](docs/existing-solutions-research.md)
- [PRD](docs/PRD.md)
- [기술 스택](docs/tech-stack.md)

## 라이센스

MIT License - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.
