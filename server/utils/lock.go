package utils

import (
    "log"
    "regexp"
    "sync"
    "time"

    "github.com/ictsc/sstate/models"
)

var (
    // RedeployStatus - チームごとの再展開の状態を管理するスレッドセーフなマップ
    RedeployStatus = sync.Map{}

    // TeamLocks - チームごとのロック情報を保持するスレッドセーフなマップ
    TeamLocks = sync.Map{}

    // RedeployQueue - 再展開リクエストを処理するためのキュー
    RedeployQueue = make(chan models.RedeployRequest, 100)

    // InQueue - キューに存在するチームのリストを保持するスレッドセーフなマップ
    InQueue = sync.Map{}

    // TeamIDPattern - チームIDが2桁の数字であることを確認する正規表現パターン
    TeamIDPattern = regexp.MustCompile(`^\d{2}$`)
)

// チームごとのロックを取得し、ロックに成功した場合はロックオブジェクトとtrueを返す。
// タイムアウトが発生した場合はnilとfalseを返す。
func TryTeamLock(teamID string) (*sync.Mutex, bool) {
    // 指定されたチームIDのロックを取得
    teamLock := GetTeamLock(teamID)
    locked := make(chan bool, 1)

    // ロックを非同期で試行し、成功したらチャネルに通知
    go func() {
        teamLock.Lock()
        select {
        case locked <- true:
            log.Printf("ロック取得: チームID=%s", teamID)
        default:
            teamLock.Unlock() // タイムアウト時のアンロックを実行
        }
    }()

    // 100ミリ秒のタイムアウトでロックの取得を試みる
    select {
    case success := <-locked:
        return teamLock, success
    case <-time.After(100 * time.Millisecond):
        return nil, false
    }
}

// 指定されたチームIDに対するロックを取得または作成し、ロックオブジェクトを返す。
func GetTeamLock(teamID string) *sync.Mutex {
    lock, _ := TeamLocks.LoadOrStore(teamID, &sync.Mutex{})
    return lock.(*sync.Mutex)
}
