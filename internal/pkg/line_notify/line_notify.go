package line_notify

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
	"go.uber.org/zap"
)

func TriggerLineNotify(users *[]repository.BinanceScheama) ([]string, []string) {
	successUser := make(chan string, len(*users))
	errUser := make(chan string, len(*users))
	var wg sync.WaitGroup
	for _, user := range *users {
		userTime := user.LineNotifyTime % 60
		currentMinute := time.Now().Minute()
		needNotify := currentMinute == int(userTime)
		if userTime > 0 {
			needNotify = (currentMinute % int(userTime)) == 0
		}
		if needNotify {
			wg.Add(1)
			go func(u repository.BinanceScheama) {
				defer wg.Done()
				err := triggerUser(&u)
				if err != nil {
					errUser <- u.LineUID
					return
				}
				successUser <- u.LineUID
			}(user)
		}
	}

	go func() {
		wg.Wait()
		defer close(successUser)
		defer close(errUser)
	}()

	successUserSlice := []string{}
	errUserSlice := []string{}

	for sUser := range successUser {
		successUserSlice = append(successUserSlice, sUser)
	}
	for eUser := range errUser {
		errUserSlice = append(errUserSlice, eUser)
	}

	return successUserSlice, errUserSlice

}

func triggerUser(u *repository.BinanceScheama) error {
	resp, _, err := utils.Http(&utils.HttpReq{
		Method: http.MethodPost,
		URL:    configs.Config.LineNotifyURL,
		Headers: map[string]string{
			configs.AuthorizationHeader: configs.Config.SecretPassphrase,
		},
		Body: u,
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		zap.L().Error("Error trigger binance line notify", zap.String("lineuid", u.LineUID))
		return fmt.Errorf("error trigger binance line notify %s", u.LineUID)
	}
	zap.L().Info("Successfully trigger binance line notify", zap.String("lineuid", u.LineUID))
	return nil
}
