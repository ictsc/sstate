package utils

import (
	"sync"

	"github.com/ictsc/sstate/models"
)

// RedeployStatus - チームと問題ごとの再展開状態を管理するスレッドセーフなマップ
var RedeployStatus = sync.Map{}

// RedeployQueue - 再展開リクエストを処理するためのキュー
var RedeployQueue = make(chan models.RedeployRequest, 100)

// InQueue - キューに存在するチーム+問題のリストを保持するスレッドセーフなマップ
var InQueue = sync.Map{}
