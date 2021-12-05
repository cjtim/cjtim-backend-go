package files

import "github.com/cjtim/cjtim-backend-go/internal/app/repository"

type FilesMocks struct {
	// Add - Upload file and save to DB
	M_Add func(fullFileName string, byteData []byte, lineUID string) (*repository.FileScheama, error)

	// AddFromLine - Contents being upload via line chat will parse here
	M_AddFromLine func(messageID string, lineUID string) (*repository.FileScheama, error)

	// Delete - Remove file from storage and rebrandly
	M_Delete func(fullFileName string, lineUID string) error
}

func (f *FilesMocks) Add(fullFileName string, byteData []byte, lineUID string) (*repository.FileScheama, error) {
	return f.M_Add(fullFileName, byteData, lineUID)
}

func (f *FilesMocks) AddFromLine(messageID string, lineUID string) (*repository.FileScheama, error) {
	return f.M_AddFromLine(messageID, lineUID)
}

func (f *FilesMocks) Delete(fullFileName string, lineUID string) error {
	return f.M_Delete(fullFileName, lineUID)
}
