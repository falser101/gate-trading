package copytrading

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CookieManager 管理 Gate.io Cookie 的加密/解密和验证
type CookieManager struct {
	cryptoKey []byte
}

// GateCookie Gate.io Cookie 结构
type GateCookie struct {
	Token     string
	CsrfToken string
	Uid       string
	ExpiresAt time.Time
}

// NewCookieManager 创建 Cookie 管理器
func NewCookieManager(encryptionKey string) *CookieManager {
	// 确保 key 是 32 字节（AES-256）
	key := make([]byte, 32)
	copy(key, []byte(encryptionKey))
	return &CookieManager{cryptoKey: key}
}

// Encrypt 使用 AES-256-GCM 加密字符串
func (m *CookieManager) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(m.cryptoKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字符串
func (m *CookieManager) Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(m.cryptoKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// ParseToken 解析 JWT Token 获取过期时间
func (m *CookieManager) ParseToken(token string) (*time.Time, error) {
	// 不验证签名，只解析 claims
	tokenParsed, _, _ := jwt.NewParser().ParseUnverified(token, &jwt.RegisteredClaims{})
	if tokenParsed == nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := tokenParsed.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if claims.ExpiresAt == nil {
		return nil, errors.New("token has no expiration")
	}

	expiry := claims.ExpiresAt.Time
	return &expiry, nil
}

// NewGateCookie 创建 GateCookie 并解析过期时间
func (m *CookieManager) NewGateCookie(token, csrfToken, uid string) (*GateCookie, error) {
	expiresAt, err := m.ParseToken(token)
	if err != nil {
		// 如果解析失败，默认 30 天过期
		defaultExpiry := time.Now().Add(30 * 24 * time.Hour)
		expiresAt = &defaultExpiry
	}

	return &GateCookie{
		Token:     token,
		CsrfToken: csrfToken,
		Uid:       uid,
		ExpiresAt: *expiresAt,
	}, nil
}

// IsExpired 检查 Cookie 是否过期
func (m *CookieManager) IsExpired(cookie *GateCookie) bool {
	return time.Now().After(cookie.ExpiresAt)
}

// IsExpiringSoon 检查 Cookie 是否即将过期（withinDays 天内）
func (m *CookieManager) IsExpiringSoon(cookie *GateCookie, withinDays int) bool {
	return time.Now().Before(cookie.ExpiresAt) && time.Now().Add(time.Duration(withinDays)*24*time.Hour).After(cookie.ExpiresAt)
}

// EncryptGateCookie 加密 GateCookie 用于存储
func (m *CookieManager) EncryptGateCookie(cookie *GateCookie) (encryptedToken string, encryptedCsrf string, err error) {
	encryptedToken, err = m.Encrypt(cookie.Token)
	if err != nil {
		return "", "", err
	}
	encryptedCsrf, err = m.Encrypt(cookie.CsrfToken)
	if err != nil {
		return "", "", err
	}
	return encryptedToken, encryptedCsrf, nil
}

// DecryptGateCookie 解密 GateCookie
func (m *CookieManager) DecryptGateCookie(encryptedToken, encryptedCsrf, uid string, expiresAt *time.Time) (*GateCookie, error) {
	token, err := m.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	csrfToken, err := m.Decrypt(encryptedCsrf)
	if err != nil {
		return nil, err
	}

	return &GateCookie{
		Token:     token,
		CsrfToken: csrfToken,
		Uid:       uid,
		ExpiresAt: *expiresAt,
	}, nil
}

// FormatCookieHeader 格式化 Cookie 请求头
func (m *CookieManager) FormatCookieHeader(cookie *GateCookie) string {
	var sb strings.Builder
	sb.WriteString("token=")
	sb.WriteString(cookie.Token)
	sb.WriteString("; csrftoken=")
	sb.WriteString(cookie.CsrfToken)
	sb.WriteString("; uid=")
	sb.WriteString(cookie.Uid)
	return sb.String()
}
