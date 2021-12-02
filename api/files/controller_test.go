package files_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjtim/cjtim-backend-go/api/files"
	"github.com/cjtim/cjtim-backend-go/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initFileRoute() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &linebot.UserProfileResponse{
			UserID:      "aaaaaaaaaabbbbbbbbb",
			DisplayName: "Unit Test",
			PictureURL:  "",
		})
		return c.Next()
	})
	app.Get("/files/list", files.List)
	app.Post("/files/upload", files.Upload)
	app.Post("/files/delete", files.Delete)
	return app
}

func Test_List_Success(t *testing.T) {
	expect := []repository.FileScheama{
		{
			FileName: "test1",
			URL:      repository.URLScheama{ShortURL: "a", Destination: "b"},
			LineUID:  "aaaaaaaaaabbbbbbbbb",
		},
	}
	orig := repository.FileRepo.Find
	restoreOrigiFind := func() {
		repository.FileRepo.Find = orig
	}
	repository.FileRepo.Find = func(data, filter interface{}, opts ...*options.FindOptions) error {
		b, _ := json.Marshal(expect)
		return json.Unmarshal(b, data)
	}
	defer restoreOrigiFind()

	// start test
	app := initFileRoute()

	req := httptest.NewRequest(http.MethodGet, "/files/list", nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer req.Body.Close()

	// compare expect and response body
	response := files.Response{}
	actualBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(actualBytes, &response)
	assert.Nil(t, err)
	assert.Equal(t, expect, response.Files)

}
