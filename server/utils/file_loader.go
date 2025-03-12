package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
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
	if err := json.Unmarshal(byteValue, &ProblemIDMapping); err != nil {
		return fmt.Errorf("JSONのパースに失敗しました: %v", err)
	}
	return nil
}

// TeamConfig は YAML の構造に合わせた構造体です。
type TeamConfig struct {
	Teams []string `yaml:"teams"`
}

// LoadTeamList - 指定された YAML ファイルからチーム一覧を読み込みます。
// ファイルを開けない場合、エラーを返します。
//
// パラメータ:
//   - filename: 読み込む YAML ファイルのパス
//
// 戻り値:
//   - チーム一覧の文字列スライスとエラーを返します。
func LoadTeamList(filename string) ([]string, error) {
	// 指定されたファイルを開く
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ファイルを開けませんでした: %v", err)
	}
	defer file.Close()

	// ファイルの内容を読み込み、YAMLとしてパースする
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %v", err)
	}

	var config TeamConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("YAMLのパースに失敗しました: %v", err)
	}

	return config.Teams, nil
}
