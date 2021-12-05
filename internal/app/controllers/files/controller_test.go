package files_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	controller_files "github.com/cjtim/cjtim-backend-go/internal/app/controllers/files"
	"github.com/cjtim/cjtim-backend-go/internal/app/repository"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/files"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	mustOpen := func(f string) *os.File {
		r, err := os.Open(f)
		if err != nil {
			pwd, _ := os.Getwd()
			fmt.Println("PWD: ", pwd)
			panic(err)
		}
		return r
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	file := mustOpen(fileName)
	fw, err := w.CreateFormFile(fieldName, file.Name())
	if err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

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
	app.Get("/files/list", controller_files.List)
	app.Post("/files/upload", controller_files.Upload)
	app.Post("/files/delete", controller_files.Delete)
	return app
}

func Test_Upload_Success(t *testing.T) {
	app := initFileRoute()

	b, w := createMultipartFormData(t, "file", "controller.go")

	req := httptest.NewRequest(http.MethodPost, "/files/upload", &b)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", w.FormDataContentType())

	mock := &files.FilesMocks{}
	mock.M_Add = func(fullFileName string, byteData []byte, lineUID string) (*repository.FileScheama, error) {
		return &repository.FileScheama{
			ID:       primitive.NewObjectID(),
			FileName: "file",
			URL: repository.URLScheama{
				ShortURL:    "h",
				Destination: "a",
			},
			LineUID: "",
		}, nil
	}
	files.Client = mock
	defer files.RestoreMocks()

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func Test_Upload_FailOnAdd(t *testing.T) {
	app := initFileRoute()

	b, w := createMultipartFormData(t, "file", "controller.go")

	req := httptest.NewRequest(http.MethodPost, "/files/upload", &b)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", w.FormDataContentType())

	mock := &files.FilesMocks{}
	files.Client = mock
	defer files.RestoreMocks()
	mock.M_Add = func(fullFileName string, byteData []byte, lineUID string) (*repository.FileScheama, error) {
		return nil, errors.New("Fake fail")
	}

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Upload_NoFile(t *testing.T) {
	app := initFileRoute()

	req := httptest.NewRequest(http.MethodPost, "/files/upload", nil)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func Test_Upload_File_Invalid(t *testing.T) {
	app := initFileRoute()

	b, w := createMultipartFormData(t, "file", "controller.go")
	failData := fmt.Sprintf("1%s", b.String())

	req := httptest.NewRequest(http.MethodPost, "/files/upload", bytes.NewBuffer([]byte(failData)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

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
	response := controller_files.Response{}
	actualBytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(actualBytes, &response)
	assert.Nil(t, err)
	assert.Equal(t, expect, response.Files)

}

func Test_List_Fail(t *testing.T) {

	repository.FileRepo.Find = func(data, filter interface{}, opts ...*options.FindOptions) error {
		return errors.New("Fail")
	}
	defer repository.RestoreRepoMock()

	// start test
	app := initFileRoute()

	req := httptest.NewRequest(http.MethodGet, "/files/list", nil)
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func Test_Delete_Success(t *testing.T) {
	// start test
	app := initFileRoute()

	mock := &files.FilesMocks{}
	files.Client = mock
	defer files.RestoreMocks()
	mock.M_Delete = func(fullFileName, lineUID string) error {
		return nil
	}

	body := []byte("{\"fileName\":\"a\"}")
	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method:    http.MethodPost,
		URL:       "/files/delete",
		BodyBytes: body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func Test_Delete_Fail(t *testing.T) {
	// start test
	app := initFileRoute()

	mock := &files.FilesMocks{}
	files.Client = mock
	defer files.RestoreMocks()
	mock.M_Delete = func(fullFileName, lineUID string) error {
		return errors.New("Fail")
	}

	body := []byte("{\"fileName\":\"a\"}")
	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method:    http.MethodPost,
		URL:       "/files/delete",
		BodyBytes: body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func Test_Delete_InvaidBody(t *testing.T) {
	// start test
	app := initFileRoute()

	body := []byte("{\"fileName\":")
	req, _ := utils.HttpPrepareRequest(&utils.HttpReq{
		Method:    http.MethodPost,
		URL:       "/files/delete",
		BodyBytes: body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	resp, err := app.Test(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}
