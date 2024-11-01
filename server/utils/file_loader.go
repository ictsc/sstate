package utils

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
)

// ProblemIDMapping - 外部ファイルから読み込んだproblem_idマッピングを格納するグローバル変数
var ProblemIDMapping map[string]string

// 指定されたJSONファイルからproblem_idのマッピングを読み込み、ProblemIDMappingに格納する
// ファイルを開けない場合はエラーを返す
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
