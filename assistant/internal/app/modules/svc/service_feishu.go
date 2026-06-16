package svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) SendToFeishu() error {
	text, err := s.readClipboard()
	if err != nil {
		s.logger.WithError(err).Warn("feishu: read clipboard failed")
		return fmt.Errorf("read clipboard: %w", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("clipboard is empty")
	}

	cfg := psl.GetConfig().Channels
	appID := cfg.Feishu.AppID
	appSecret := cfg.Feishu.AppSecret
	chatID := cfg.Feishu.ChatID
	if appID == "" || appSecret == "" || chatID == "" {
		return fmt.Errorf("missing feishu config: app_id, app_secret, chat_id")
	}

	return s.sendToFeishuText(appID, appSecret, chatID, text)
}

func (s *Service) sendToFeishuText(appID, appSecret, chatID, text string) error {
	token, err := s.getFeishuToken(appID, appSecret)
	if err != nil {
		s.logger.WithError(err).Warn("feishu: get token failed")
		return fmt.Errorf("get feishu token: %w", err)
	}

	escaped, err := json.Marshal(text)
	if err != nil {
		return fmt.Errorf("escape text: %w", err)
	}
	contentStr := fmt.Sprintf(`{"text":%s}`, string(escaped))

	payload := struct {
		ReceiveID string `json:"receive_id"`
		MsgType   string `json:"msg_type"`
		Content   string `json:"content"`
	}{
		ReceiveID: chatID,
		MsgType:   "text",
		Content:   contentStr,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=chat_id", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "assistant/cmd")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.WithError(err).Warn("feishu: HTTP request failed")
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if err != nil {
		return fmt.Errorf("read feishu response: %w", err)
	}
	if resp.StatusCode >= http.StatusMultipleChoices {
		s.logger.WithFields(map[string]interface{}{
			"status": resp.StatusCode,
			"body":   strings.TrimSpace(string(respBody)),
		}).Warn("feishu: API returned non-2xx")
		return fmt.Errorf("feishu API returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(respBody, &result); err == nil && result.Code != 0 {
		s.logger.WithFields(map[string]interface{}{
			"code": result.Code,
			"msg":  result.Msg,
		}).Warn("feishu: business error")
		return fmt.Errorf("feishu API error code=%d: %s", result.Code, result.Msg)
	}
	return nil
}

func (s *Service) getFeishuToken(appID, appSecret string) (string, error) {
	payload := map[string]string{"app_id": appID, "app_secret": appSecret}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost,
		"https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal",
		bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "assistant/cmd")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("get feishu token failed: %s %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	var result struct {
		Code              int    `json:"code"`
		Msg               string `json:"msg"`
		TenantAccessToken string `json:"tenant_access_token"`
		Expire            int    `json:"expire"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if result.Code != 0 {
		return "", fmt.Errorf("get feishu token failed: %s", result.Msg)
	}
	if result.TenantAccessToken == "" {
		return "", fmt.Errorf("get feishu token failed: empty token")
	}
	return result.TenantAccessToken, nil
}
