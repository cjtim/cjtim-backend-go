package binance_test

import (
	"testing"
	"time"

	"github.com/cjtim/cjtim-backend-go/repository"
)

func Test_Cron(t *testing.T) {
	user := repository.BinanceScheama{
		LineNotifyTime: int64(time.Now().Minute()),
	}
	userTime := (user.LineNotifyTime) % 60
	currentMinute := time.Now().Minute()
	needNotify := (currentMinute % int(userTime)) == 0
	t.Log(currentMinute)
	if needNotify {
		t.Log(currentMinute)
	} else {
		t.Fatal("Not right time")
	}
}
