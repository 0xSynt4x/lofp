package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Service sends transactional emails via Resend.
type Service struct {
	apiKey  string
	fromAddr string
	baseURL  string // frontend URL for links
}

// New creates a new email service. If apiKey is empty, emails are logged but not sent.
func New(apiKey, fromAddr, baseURL string) *Service {
	if fromAddr == "" {
		fromAddr = "Legends of Future Past <noreply@lofp.metavert.io>"
	}
	return &Service{apiKey: apiKey, fromAddr: fromAddr, baseURL: baseURL}
}

// Enabled returns true if the email service is configured.
func (s *Service) Enabled() bool {
	return s.apiKey != ""
}

type resendRequest struct {
	From    string `json:"from"`
	To      []string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

func (s *Service) send(to, subject, html string) error {
	if s.apiKey == "" {
		fmt.Printf("[email-dev] To: %s | Subject: %s\n", to, subject)
		return nil
	}
	body := resendRequest{
		From:    s.fromAddr,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	}
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend API error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend API returned %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// SendVerification sends an email verification link and code.
func (s *Service) SendVerification(to, token, code string) error {
	link := fmt.Sprintf("%s/verify-email?token=%s", s.baseURL, token)
	html := fmt.Sprintf(`
<div style="font-family: monospace; max-width: 600px; margin: 0 auto; background: #0a0a0a; color: #e0e0e0; padding: 32px; border: 1px solid #333;">
  <h1 style="color: #f59e0b; font-size: 18px;">Legends of Future Past</h1>
  <p>Welcome to the Shattered Realms!</p>
  <p>Click the button below to verify your email address:</p>
  <p style="text-align: center; margin: 24px 0;"><a href="%s" style="background: #b45309; color: white; padding: 12px 24px; text-decoration: none; font-size: 16px;">Verify Email Address</a></p>
  <p>Or enter this verification code on the website or in your MUD client:</p>
  <p style="text-align: center; font-size: 24px; letter-spacing: 4px; color: #f59e0b; background: #1a1a1a; padding: 12px; border: 1px solid #444;">%s</p>
  <p style="color: #888; font-size: 12px;">This link and code expire in 24 hours. If you didn't create this account, you can ignore this email.</p>
</div>`, link, code)
	return s.send(to, "Verify your email — Legends of Future Past", html)
}

// SendPasswordReset sends a password reset link.
func (s *Service) SendPasswordReset(to, token string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", s.baseURL, token)
	html := fmt.Sprintf(`
<div style="font-family: monospace; max-width: 600px; margin: 0 auto; background: #0a0a0a; color: #e0e0e0; padding: 32px; border: 1px solid #333;">
  <h1 style="color: #f59e0b; font-size: 18px;">Legends of Future Past</h1>
  <p>A password reset was requested for your account.</p>
  <p>Click the link below to set a new password:</p>
  <p><a href="%s" style="color: #f59e0b; font-size: 16px;">Reset Password</a></p>
  <p style="color: #888; font-size: 12px;">This link expires in 1 hour. If you didn't request this, you can ignore this email.</p>
</div>`, link)
	return s.send(to, "Reset your password — Legends of Future Past", html)
}
