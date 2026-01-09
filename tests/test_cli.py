"""CLI 테스트."""

from typer.testing import CliRunner

from hwp2markdown import __version__
from hwp2markdown.cli import app

runner = CliRunner()


def test_version() -> None:
    """버전 출력 테스트."""
    result = runner.invoke(app, ["--version"])
    assert result.exit_code == 0
    assert __version__ in result.stdout


def test_help() -> None:
    """도움말 출력 테스트."""
    result = runner.invoke(app, ["--help"])
    assert result.exit_code == 0
    assert "HWP" in result.stdout
    assert "Markdown" in result.stdout
