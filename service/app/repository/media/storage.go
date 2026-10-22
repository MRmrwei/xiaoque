package media

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"suki/xqyun/app/utils/ulid"
)

// Storage 定义素材文件存储能力。
// 当前提供本地磁盘实现，后续可增加对象存储（OSS / COS / S3）等实现，
// service 层只依赖这个接口，切换存储方式不需要改业务逻辑。
type Storage interface {
	// Save 写入文件，返回相对文件名与可直接访问的 URL。
	Save(reader io.Reader, ext string, at time.Time) (filename string, url string, err error)
	// Remove 删除已保存的文件，文件不存在时视为成功。
	Remove(filename string) error
}

// NewStorage 按配置创建存储实现；新增实现时在这里注册 provider。
func NewStorage(provider, dir, urlPrefix, publicBaseURL string) (Storage, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", "local":
		return NewLocalStorage(dir, urlPrefix, publicBaseURL), nil
	default:
		return nil, fmt.Errorf("不支持的素材存储方式：%s", provider)
	}
}

// LocalStorage 把素材保存到本地磁盘，并返回静态访问地址。
type LocalStorage struct {
	dir           string
	urlPrefix     string
	publicBaseURL string
}

var _ Storage = (*LocalStorage)(nil)

// NewLocalStorage 创建本地素材存储。
func NewLocalStorage(dir string, urlPrefix string, publicBaseURL string) *LocalStorage {
	return &LocalStorage{
		dir:           dir,
		urlPrefix:     strings.TrimRight(urlPrefix, "/"),
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}
}

// Save 把文件写入 dir/yyyy/mm/，返回相对路径（写入数据库的 filename）与完整访问 URL。
func (s *LocalStorage) Save(reader io.Reader, ext string, at time.Time) (string, string, error) {
	id, err := ulid.New()
	if err != nil {
		return "", "", fmt.Errorf("生成素材文件名失败：%w", err)
	}

	filename := strings.ToLower(id) + "." + strings.TrimPrefix(strings.ToLower(ext), ".")
	// 数据库与 URL 统一使用 /，磁盘路径在写入时再转换为系统分隔符。
	relative := path.Join(at.Format("2006"), at.Format("01"), filename)
	fullPath := filepath.Join(s.dir, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", "", fmt.Errorf("创建素材目录失败：%w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", "", fmt.Errorf("创建素材文件失败：%w", err)
	}
	if _, err := io.Copy(file, reader); err != nil {
		_ = file.Close()
		_ = os.Remove(fullPath)
		return "", "", fmt.Errorf("写入素材文件失败：%w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(fullPath)
		return "", "", fmt.Errorf("保存素材文件失败：%w", err)
	}

	urlPath := s.urlPrefix + "/" + relative
	if s.publicBaseURL != "" {
		urlPath = s.publicBaseURL + urlPath
	}
	return relative, urlPath, nil
}

// Remove 删除已保存的文件；文件不存在时视为成功。
func (s *LocalStorage) Remove(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return nil
	}
	if err := os.Remove(filepath.Join(s.dir, filepath.FromSlash(filename))); err != nil &&
		!os.IsNotExist(err) {
		return fmt.Errorf("删除素材文件失败：%w", err)
	}
	return nil
}
