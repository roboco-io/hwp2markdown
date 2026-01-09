package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	outputFile    string
	extractImages string
	format        string
	verbose       bool
	quiet         bool
)

var convertCmd = &cobra.Command{
	Use:   "convert [input]",
	Short: "HWP/HWPX 파일을 Markdown으로 변환",
	Long: `HWP 또는 HWPX 파일을 Markdown으로 변환합니다.

예시:
  hwp2markdown convert document.hwpx
  hwp2markdown convert document.hwp -o output.md
  hwp2markdown convert document.hwpx --extract-images ./images`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func init() {
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "출력 파일 (기본값: stdout)")
	convertCmd.Flags().StringVar(&extractImages, "extract-images", "", "이미지 추출 디렉토리")
	convertCmd.Flags().StringVarP(&format, "format", "f", "markdown", "출력 포맷 (markdown, text)")
	convertCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "상세 출력")
	convertCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "조용한 모드")

	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// 파일 존재 확인
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("입력 파일을 찾을 수 없습니다: %s", inputFile)
	}

	// 파일 확장자 확인
	ext := strings.ToLower(filepath.Ext(inputFile))
	if ext != ".hwp" && ext != ".hwpx" {
		return fmt.Errorf("지원하지 않는 파일 형식입니다: %s (지원: .hwp, .hwpx)", ext)
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "입력 파일: %s\n", inputFile)
		fmt.Fprintf(os.Stderr, "출력 포맷: %s\n", format)
		if outputFile != "" {
			fmt.Fprintf(os.Stderr, "출력 파일: %s\n", outputFile)
		}
		if extractImages != "" {
			fmt.Fprintf(os.Stderr, "이미지 추출: %s\n", extractImages)
		}
	}

	// TODO: 실제 변환 로직 구현
	if !quiet {
		fmt.Fprintln(os.Stderr, "\n변환 기능은 아직 구현 중입니다.")
	}

	return nil
}
