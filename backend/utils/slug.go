package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var (
	slugNonWord = regexp.MustCompile(`[^a-z0-9]+`) // 非英数字符
	slugDash    = regexp.MustCompile(`-+`)          // 连续连字符
)

// GenerateSlug 由标题生成 slug：
// 纯英文/数字标题转为小写连字符形式；含中文（或其他非 ASCII 字符）时降级为 article-<Unix秒>
func GenerateSlug(title string) string {
	if hasNonASCII(title) {
		return fmt.Sprintf("article-%d", time.Now().Unix())
	}

	slug := strings.ToLower(strings.TrimSpace(title))
	slug = slugNonWord.ReplaceAllString(slug, "-")
	slug = slugDash.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return fmt.Sprintf("article-%d", time.Now().Unix())
	}
	return slug
}

// EnsureUnique 回调式去重：若 base 已存在，则依次追加 -1、-2 … 直到不存在
func EnsureUnique(base string, exists func(string) bool) string {
	if !exists(base) {
		return base
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !exists(candidate) {
			return candidate
		}
	}
}

// hasNonASCII 判断字符串是否包含非 ASCII 字符（如中文）
func hasNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}
