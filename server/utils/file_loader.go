package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ProblemIDMapping - 外部ファイルから読み込んだ問題IDマッピングを格納するグローバル変数です。
var ProblemIDMapping map[string]string

// LoadProblemIDMapping - 指定された JSON ファイルから問題IDのマッピングを読み込み、
// ProblemIDMapping に格納します。ファイルを開けない場合、エラーを返します。
//
// パラメータ:
//   - filename: 読み込む JSON ファイルのパス
//
// 戻り値:
//   - エラーが発生した場合は error を返します。それ以外の場合は nil を返します。
func LoadProblemIDMapping(filename string) error {
	// 指定されたファイルを開く
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ファイルを開けませんでした: %v", err)
	}
	defer file.Close() // 関数終了時にファイルを閉じる

	// ファイルの内容を読み込み、JSONとしてパースしてProblemIDMappingに格納
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &ProblemIDMapping)
	return nil
}
