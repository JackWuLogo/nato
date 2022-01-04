package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"micro-libs/app"
	"micro-libs/utils/log"
	"micro-libs/utils/tool"
	"mime/multipart"
	"os"
	"path/filepath"
)

// 上传文件返回的结果
type Upload struct {
	SaveFile string `json:"save_file"` // 保存文件
	FileName string `json:"file_name"` // 文件名称
	FileSize int64  `json:"file_size"` // 文件大小
	OldName  string `json:"old_name"`  // 原始名称
}

// 上传文件
func FileUpload(file *multipart.FileHeader, savePath string) (*Upload, error) {
	rootPath := app.Opts.StoreRoot
	fullPath := filepath.Join(rootPath, savePath)

	if err := tool.InitFolder(fullPath, 0755); err != nil {
		return nil, err
	}

	// 读取上传文件
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 获取文件扩展名
	oldName := filepath.Base(file.Filename)
	fileExt := filepath.Ext(oldName)
	// 生产文件名称
	filename := fmt.Sprintf("%s%s", tool.UUID(), fileExt)

	// Destination
	saveFile := filepath.Join(fullPath, filename)
	dst, err := os.Create(saveFile)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	res := &Upload{
		SaveFile: filepath.Join(savePath, filename),
		FileName: filename,
		FileSize: file.Size,
		OldName:  oldName,
	}
	log.Debug("[UPLOAD FILE] %+v", res)

	return res, nil
}

// FileAssets 文件资源访问
func FileAssets(ctx echo.Context) error {
	filename := ctx.Param("*")
	basename := filepath.Base(filename)
	ext := filepath.Ext(basename)

	if ext == "" {
		return echo.ErrNotFound
	}

	fullPath := filepath.Join(app.Opts.StoreRoot, filename)
	if s, err := os.Stat(fullPath); os.IsNotExist(err) {
		return echo.ErrNotFound
	} else if s.IsDir() {
		return echo.ErrNotFound
	}

	return ctx.File(fullPath)
}
