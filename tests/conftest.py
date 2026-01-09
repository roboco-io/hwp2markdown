"""pytest 설정 및 공통 fixture."""

from pathlib import Path

import pytest


@pytest.fixture
def fixtures_dir() -> Path:
    """테스트 fixture 디렉토리 경로."""
    return Path(__file__).parent / "fixtures"
