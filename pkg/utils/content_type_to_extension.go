package utils

import "mime"

func ContentTypeToExtension(contentType string) ([]string, error) {
	ext, err := mime.ExtensionsByType(contentType)
	if ext == nil {
		ext = []string{".bak"}
	}
	if err != nil {
		return nil, err
	}
	return ext, nil
}
