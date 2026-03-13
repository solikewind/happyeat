package timeutil

import "time"

// 移除手动定义的 ISO8601Layout，优先使用 time.RFC3339
const (
	// DateOnlyLayout 保持不变，用于仅需日期的场景
	DateOnlyLayout = "2006-01-02"
)

// TimeToString 将 time.Time 转换为标准的 UTC RFC3339 字符串 (带 Z)
// 示例：2024-03-14T15:04:05Z
func TimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// 使用内置的 time.RFC3339，它比手动拼接字符串更健壮
	return t.UTC().Format(time.RFC3339)
}

// TimeToDateString 仅转换日期部分（基于 UTC 日期）
// 示例：2024-03-14
func TimeToDateString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(DateOnlyLayout)
}

// ParseTime 将前端传回的 ISO 字符串转回 time.Time
// 这是一个实用的补充，因为转换通常是双向的
func ParseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}
