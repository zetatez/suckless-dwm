package svc

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"assistant/internal/bootstrap/psl"
	"assistant/pkg/dwmblocknotify"
	"assistant/pkg/llm"
)

var leetCodeSystemPrompt = `你是顶级算法工程师，正在参加技术面试。解决用户给出的算法题，用 Golang 实现。

要求：

1. 给出最优解
2. 代码末尾注释给出最优解分析思路(分步)、算法思想(分步)、时间复杂度、空间复杂度
3. 要给出完整可运行的 Go 代码，包含 main() 和题目中的测试用例, 如题目中没有测试用例，那么给出 2 个测试用例
4. 代码要简洁，不要冗余注释和变量
5. 不要输出除代码和注释外的任何内容

输出格式：

package main

import "fmt"

func FuncName(...) { ... }

func main() { ... }

/*
问题分析:
1. ...
2. ...

算法思想:(DP|Greedy|DFS|BFS|回溯|二分法|双指针|滑动窗口|...)
1. ...
2. ...

时间复杂度: O(...), why?
空间复杂度: O(...), why?
*/
`

func stripThinkTags(s string) string {
	for {
		start := strings.Index(s, "<think>")
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], "</think>")
		if end < 0 {
			s = strings.TrimSpace(s[:start])
			break
		}
		s = strings.TrimSpace(s[:start] + s[start+end+len("</think>"):])
	}
	return s
}

func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	for _, prefix := range []string{"```go\n", "```python\n", "```rust\n", "```sql\n", "```"} {
		if strings.HasPrefix(s, prefix) {
			s = strings.TrimSpace(s[len(prefix):])
			break
		}
	}
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSpace(s[:len(s)-3])
	}
	return s
}

func (s *Service) SolveLeetCode() error {
	text, err := s.readClipboard()
	if err != nil {
		return fmt.Errorf("read clipboard: %w", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("clipboard is empty")
	}
	dwmblocknotify.PUT("!...", 3*time.Second)

	client := psl.GetLLMClient()
	if client == nil {
		return fmt.Errorf("LLM client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := llm.Complete(ctx, client, text, llm.WithSystemPrompt(leetCodeSystemPrompt), llm.WithTemperature(0.3))
	if err != nil {
		return fmt.Errorf("LLM request: %w", err)
	}

	result := strings.TrimSpace(resp.Content)
	result = stripThinkTags(result)
	result = stripCodeFence(result)
	if result == "" {
		return fmt.Errorf("LLM returned empty response")
	}

	if err := s.writeClipboard(result); err != nil {
		return fmt.Errorf("write clipboard: %w", err)
	}

	leetCodeOutputFile := "/home/shiyi/git/test/go/x.go"
	if err := os.MkdirAll(filepath.Dir(leetCodeOutputFile), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	if err := os.WriteFile(leetCodeOutputFile, []byte(result), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	dwmblocknotify.PUT("!!!", 5*time.Second)
	return nil
}

func (s *Service) SolveLeetCodeScreenshot() error {
	dwmblocknotify.PUT("!...", 3*time.Second)
	imgBase64, err := s.screenshot()
	if err != nil {
		return fmt.Errorf("screenshot: %w", err)
	}

	client := psl.GetLLMClient()
	if client == nil {
		return fmt.Errorf("LLM client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	prompt := "请识别截图中显示的算法题，用 Golang 实现最优解。"

	resp, err := llm.Complete(ctx, client, prompt,
		llm.WithSystemPrompt(leetCodeSystemPrompt),
		llm.WithTemperature(0.3),
		llm.WithImageBase64(imgBase64),
	)
	if err != nil {
		return fmt.Errorf("LLM request: %w", err)
	}

	result := strings.TrimSpace(resp.Content)
	result = stripThinkTags(result)
	result = stripCodeFence(result)
	if result == "" {
		return fmt.Errorf("LLM returned empty response")
	}

	if err := s.writeClipboard(result); err != nil {
		return fmt.Errorf("write clipboard: %w", err)
	}

	leetCodeOutputFile := "/home/shiyi/git/test/go/x.go"
	if err := os.MkdirAll(filepath.Dir(leetCodeOutputFile), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	if err := os.WriteFile(leetCodeOutputFile, []byte(result), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	dwmblocknotify.PUT("!!!", 5*time.Second)
	return nil
}

func (s *Service) screenshot() (string, error) {
	tool := "flameshot"
	if _, err := exec.LookPath(tool); err != nil {
		return "", fmt.Errorf("screenshot tool not found: %s", tool)
	}
	cmd := exec.Command(tool, "full", "--raw")
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("screenshot failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(out), nil
}
