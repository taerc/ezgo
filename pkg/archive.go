package ezgo

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ZipArchive

func ZipArchive(files []string, destPath string) error {

	fp, _ := os.Create(destPath)
	defer func() {
		_ = fp.Close()
	}()

	zw := zip.NewWriter(fp)
	defer func() {
		_ = zw.Close()
	}()

	for _, file := range files {
		if fw, e := zw.Create(filepath.Base(file)); e == nil {
			if fr, e1 := os.Open(file); e1 == nil {
				if nw, e2 := io.Copy(fw, fr); e2 == nil {
					Info(nil, M, fmt.Sprintf(" %d bytes write", nw))
				} else {
					Error(nil, M, "write failed")
				}
				fr.Close()
			} else {
				Error(nil, M, fmt.Sprintf("read [%s] failed", file))
			}
		} else {
			Error(nil, M, fmt.Sprintf("create [%s] failed", filepath.Base(file)))
		}
	}

	return nil
}
