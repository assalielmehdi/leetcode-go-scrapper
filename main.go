package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
)

const editorSelector string = "html>body>div>div:nth-child(3)>div>div>div>div>div:nth-child(3)>div>div>div>div>div>div:nth-child(3)>div>div>div>div>div:nth-child(2)>div>div:nth-child(4)"

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Format: leetcode-scrapper <problem-url>")
	}

	url := os.Args[1]

	lines := parseCode(url)

	writeCode(url, lines)
}

func parseCode(url string) []string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var code []string
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible("body"),
		chromedp.Evaluate(`
			localStorage.setItem('2487_0_lang', 'golang');
			localStorage.setItem('global_lang', 'golang');
		`, nil),
		chromedp.Reload(),
		chromedp.WaitVisible(editorSelector),
		chromedp.Evaluate(`
		[...document.querySelector('html>body>div>div:nth-child(3)>div>div>div>div>div:nth-child(3)>div>div>div>div>div>div:nth-child(3)>div>div>div>div>div:nth-child(2)>div>div:nth-child(4)').children].map(line => [...line.children[0].children].map(el => el.innerText).join(''))
		`, &code),
	})
	if err != nil {
		log.Fatal(err)
	}

	return code
}

func writeCode(url string, lines []string) {
	slug := extractSlug(url)
	file := fmt.Sprintf("./leetcode/%s.go", slug)

	code := fmt.Sprintf("package leetcode\n\n%s", strings.Join(lines, "\n"))
	code = strings.ReplaceAll(code, "\u00A0", " ")

	if err := os.WriteFile(file, []byte(code), 0o644); err != nil {
		log.Fatal(err)
	}

}

func extractSlug(url string) string {
	url, _ = strings.CutPrefix(url, "https://leetcode.com/problems/")
	url = strings.TrimSuffix(url, "/description/")
	url = strings.TrimSuffix(url, "/")
	url = strings.ReplaceAll(url, "-", "_")

	return url
}
