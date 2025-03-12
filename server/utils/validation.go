package utils

import "regexp"
import "log"

// TeamIDPattern - チームIDが2桁の数字であることを確認する正規表現パターン
var TeamIDPattern = regexp.MustCompile(`^\d{2}$`)

// ValidateTeamIDs - チームID一覧を受け取り、正しい形式（0埋めされた2桁の数字）であるものだけを返します。
// バリデーションに失敗したチームIDはログに出力します。
func ValidateTeamIDs(teams []string) []string {
	validTeams := []string{}
	for _, teamID := range teams {
		if TeamIDPattern.MatchString(teamID) {
			validTeams = append(validTeams, teamID)
		} else {
			log.Printf("無効なチームID: %s", teamID)
		}
	}
	return validTeams
}
// ValidTeamIDs - 有効なチームIDを保持するグローバルマップです。
var ValidTeamIDs map[string]bool

// SetValidTeamIDs - 有効なチームID一覧をマップに設定します。
func SetValidTeamIDs(teams []string) {
	ValidTeamIDs = make(map[string]bool)
	for _, team := range teams {
		ValidTeamIDs[team] = true
	}
}

// IsValidTeamID - 指定された teamID が有効なチームIDかどうかをチェックします。
func IsValidTeamID(teamID string) bool {
	// ValidTeamIDs が未設定の場合はチェックをスキップ（true を返す）
	if ValidTeamIDs == nil {
		return true
	}
	return ValidTeamIDs[teamID]
}
