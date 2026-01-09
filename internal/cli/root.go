package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "hwp2markdown",
	Short: "HWP/HWPX 문서를 Markdown으로 변환",
	Long: `hwp2markdown은 HWP(한글 워드프로세서) 문서를 Markdown으로 변환하는 CLI 도구입니다.

지원 포맷:
  - HWPX (XML 기반 개방형 포맷)
  - HWP 5.x (OLE/Compound 바이너리 포맷)

사용 예시:
  hwp2markdown convert document.hwpx -o output.md
  hwp2markdown convert document.hwp --extract-images ./images`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "버전 정보 출력",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hwp2markdown %s\n", version)
	},
}

func SetVersion(v string) {
	version = v
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
