package utils

import "regexp"

// TeamIDPattern - チームIDが2桁の数字であることを確認する正規表現パターン
var TeamIDPattern = regexp.MustCompile(`^\d{2}$`)
