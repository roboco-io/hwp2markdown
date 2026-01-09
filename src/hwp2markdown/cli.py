"""CLI 인터페이스."""

from pathlib import Path
from typing import Annotated, Optional

import typer
from rich.console import Console

from hwp2markdown import __version__

app = typer.Typer(
    name="hwp2markdown",
    help="HWP(한글 워드프로세서) 문서를 Markdown으로 변환합니다.",
    add_completion=False,
)
console = Console()


def version_callback(value: bool) -> None:
    """버전 정보 출력."""
    if value:
        console.print(f"hwp2markdown {__version__}")
        raise typer.Exit()


@app.command()
def convert(
    input_file: Annotated[
        Path,
        typer.Argument(
            help="입력 HWP/HWPX 파일",
            exists=True,
            readable=True,
        ),
    ],
    output: Annotated[
        Optional[Path],
        typer.Option(
            "-o",
            "--output",
            help="출력 파일 (기본값: stdout)",
        ),
    ] = None,
    format_: Annotated[
        str,
        typer.Option(
            "-f",
            "--format",
            help="출력 포맷",
        ),
    ] = "markdown",
    extract_images: Annotated[
        Optional[Path],
        typer.Option(
            "--extract-images",
            help="이미지 추출 디렉토리",
        ),
    ] = None,
    verbose: Annotated[
        bool,
        typer.Option(
            "-v",
            "--verbose",
            help="상세 출력",
        ),
    ] = False,
    quiet: Annotated[
        bool,
        typer.Option(
            "-q",
            "--quiet",
            help="조용한 모드",
        ),
    ] = False,
    version: Annotated[
        bool,
        typer.Option(
            "-V",
            "--version",
            help="버전 출력",
            callback=version_callback,
            is_eager=True,
        ),
    ] = False,
) -> None:
    """HWP/HWPX 파일을 Markdown으로 변환합니다."""
    _ = version  # unused but needed for callback
    # TODO: 변환 로직 구현
    if not quiet:
        console.print(f"[bold]입력 파일:[/bold] {input_file}")
        console.print(f"[bold]출력 포맷:[/bold] {format_}")

    if verbose:
        console.print(f"[dim]출력 파일: {output or 'stdout'}[/dim]")
        if extract_images:
            console.print(f"[dim]이미지 추출: {extract_images}[/dim]")

    # 임시: 미구현 메시지
    console.print("[yellow]변환 기능은 아직 구현되지 않았습니다.[/yellow]")


if __name__ == "__main__":
    app()
