package convert

import (
	"os/exec"
)

func MarkdownToHTML(file_path string) (string, error) {
	cmd := exec.Command("sh", "-c", "/usr/bin/pandoc -f markdown_github -t html5 -s " + file_path)
	res, err := cmd.Output()

	return string(res), err
}
