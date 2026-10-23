package impl

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	mediarepository "suki/xqyun/app/repository/media"
	"suki/xqyun/app/services/media"
)

// allowedVideoExtensions 是允许上传的视频扩展名。
var allowedVideoExtensions = map[string]struct{}{
	".mp4":  {},
	".webm": {},
	".ogv":  {},
	".ogg":  {},
	".mov":  {},
	".m4v":  {},
}

// allowedImageExtensions 是允许上传的图片扩展名。
// 暂不接收 SVG：同源直开会执行脚本，后续需要时再做清洗。
var allowedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
	".gif":  {},
	".avif": {},
}

// FileUploader 是 media.FileUploader 的默认实现：
// 只依赖存储与大小限制，不访问数据库。
type FileUploader struct {
	storage mediarepository.Storage
	limits  media.Limits
	now     func() time.Time
}

var _ media.FileUploader = (*FileUploader)(nil)

// NewFileUploader 创建纯上传实现；storage 通过接口注入，可切换本地磁盘或对象存储。
func NewFileUploader(storage mediarepository.Storage, limits media.Limits) *FileUploader {
	return &FileUploader{storage: storage, limits: limits, now: time.Now}
}

// Upload 校验素材类型与大小，写入存储并返回文件信息；不落库。
// Kind 为空时按扩展名推断；显式传入时必须是 video / image 且与扩展名匹配。
func (u *FileUploader) Upload(_ context.Context, input media.UploadInput) (media.StoredFile, error) {
	originalName := filepath.Base(strings.TrimSpace(input.OriginalName))
	if originalName == "" || originalName == "." || input.Size <= 0 {
		return media.StoredFile{}, media.ErrEmptyFile
	}
	if len(originalName) > 255 {
		return media.StoredFile{}, media.ErrOriginalNameTooLong
	}

	ext := strings.ToLower(filepath.Ext(originalName))
	kind, ok := detectKind(ext)
	if !ok {
		return media.StoredFile{}, media.ErrUnsupportedFile
	}
	if requested := strings.TrimSpace(input.Kind); requested != "" && requested != kind {
		return media.StoredFile{}, media.ErrUnsupportedFile
	}

	maxBytes := u.limits.MaxImageBytes
	expectedMimePrefix := "image/"
	if kind == media.KindVideo {
		maxBytes = u.limits.MaxVideoBytes
		expectedMimePrefix = "video/"
	}
	if maxBytes > 0 && input.Size > maxBytes {
		return media.StoredFile{}, media.ErrFileTooLarge
	}

	mime := strings.ToLower(strings.TrimSpace(input.Mime))
	if mime != "" && mime != "application/octet-stream" &&
		!strings.HasPrefix(mime, expectedMimePrefix) {
		return media.StoredFile{}, media.ErrUnsupportedFile
	}

	now := u.now()
	filename, url, err := u.storage.Save(input.Reader, ext, now)
	if err != nil {
		return media.StoredFile{}, err
	}

	return media.StoredFile{
		Kind:         kind,
		OriginalName: originalName,
		Filename:     filename,
		URL:          url,
		Mime:         mime,
		Size:         input.Size,
		CreatedAt:    now.Unix(),
	}, nil
}

// Remove 删除已写入的文件。
func (u *FileUploader) Remove(filename string) error {
	return u.storage.Remove(filename)
}

// detectKind 按扩展名判断素材类型。
func detectKind(ext string) (string, bool) {
	if _, ok := allowedVideoExtensions[ext]; ok {
		return media.KindVideo, true
	}
	if _, ok := allowedImageExtensions[ext]; ok {
		return media.KindImage, true
	}
	return "", false
}
