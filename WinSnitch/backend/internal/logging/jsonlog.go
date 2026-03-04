package logging

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func New(programDataDir string) (*logrus.Logger, error) {
	if err := os.MkdirAll(programDataDir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(programDataDir, "winsnitch.json.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	l := logrus.New()
	l.SetOutput(f)
	l.SetFormatter(&logrus.JSONFormatter{})
	return l, nil
}
