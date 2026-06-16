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
	"assistant/pkg/llm"
)

var leetCodeSystemPrompt = `你是顶级算法工程师，正在参加技术面试。解决用户给出的算法题，用 Golang 实现。

要求：

1. 给出最优解
2. 代码末尾用块注释给出最优解分析思路(分步)、算法思想(分步)、时间复杂度、空间复杂度
3. 要给出完整可运行的 Go 代码，包含 main() 和题目中的测试用例, 如题目中没有测试用例，那么给出 2 个测试用例
4. 不要输出代码块以外的任何内容
5. 代码要简洁，不要冗余注释和变量

输出格式：

package main

import "fmt"

func solve(...) { ... }

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

func (s *Service) SolveLeetCode() error {
	text, err := s.readClipboard()
	if err != nil {
		return fmt.Errorf("read clipboard: %w", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("clipboard is empty")
	}

	s.notify("!...")

	client := psl.GetLLMClient()
	if client == nil {
		return fmt.Errorf("LLM client not initialized")
	}

	cfg := psl.GetConfig().LLM
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := llm.Complete(ctx, client, text,
		llm.WithSystemPrompt(leetCodeSystemPrompt),
		llm.WithTemperature(0.3),
		llm.WithMaxTokens(cfg.MaxTokens),
	)
	if err != nil {
		return fmt.Errorf("LLM request: %w", err)
	}

	result := strings.TrimSpace(resp.Content)
	result = stripThinkTags(result)
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

	s.notify("!!!")
	return nil
}

func (s *Service) SolveLeetCodeScreenshot() error {
	s.notify("!...")
	imgBase64, err := s.screenshot()
	if err != nil {
		return fmt.Errorf("screenshot: %w", err)
	}

	client := psl.GetLLMClient()
	if client == nil {
		return fmt.Errorf("LLM client not initialized")
	}

	cfg := psl.GetConfig().LLM
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	prompt := "请识别截图中显示的算法题，用 Golang 实现最优解。"

	resp, err := llm.Complete(ctx, client, prompt,
		llm.WithSystemPrompt(leetCodeSystemPrompt),
		llm.WithTemperature(0.3),
		llm.WithMaxTokens(cfg.MaxTokens),
		llm.WithImageBase64(imgBase64),
	)
	if err != nil {
		return fmt.Errorf("LLM request: %w", err)
	}

	result := strings.TrimSpace(resp.Content)
	result = stripThinkTags(result)
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

	s.notify("!!!")
	return nil
}

func (s *Service) screenshot() (string, error) {
	tool := "import"
	if _, err := exec.LookPath(tool); err != nil {
		return "", fmt.Errorf("screenshot tool not found: %s", tool)
	}
	tmp := filepath.Join(os.TempDir(), "assistant-screenshot.png")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, tool, "-window", "root", tmp)
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("screenshot failed: %w", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		return "", fmt.Errorf("read screenshot: %w", err)
	}
	os.Remove(tmp)
	return base64.StdEncoding.EncodeToString(data), nil
}
