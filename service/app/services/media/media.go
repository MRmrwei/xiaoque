// Package media 定义媒体素材的业务接口与上传入参。
//
// 所有上传都必须通过 Uploader 接口消费：控制器只依赖接口，
// 具体实现（校验、存储、落库）由 impl 包提供；以后要更换上传逻辑，
// 新增一个实现并在装配处替换即可，控制器与存储实现都不用改。
package media

import (
	"context"
	"errors"
	"io"

	"suki/xqyun/app/models"
)

// 素材类型。
const (
	KindVideo = "video"
	KindImage = "image"
)

var (
	// ErrEmptyFile 表示没有选择文件。
	ErrEmptyFile = errors.New("请选择要上传的文件")
	// ErrUnsupportedFile 表示文件类型不支持。
	ErrUnsupportedFile = errors.New("仅支持 MP4 / WebM / MOV 视频，或 JPG / PNG / WebP / GIF 图片")
	// ErrUnsupportedKind 表示素材类型不支持。
	ErrUnsupportedKind = errors.New("暂不支持该素材类型")
	// ErrFileTooLarge 表示文件超过配置的大小上限。
	ErrFileTooLarge = errors.New("文件超过大小限制")
	// ErrOriginalNameTooLong 表示原始文件名过长。
	ErrOriginalNameTooLong = errors.New("文件名不能超过 255 个字符")
	// ErrOriginalNameRequired 表示素材名称为空。
	ErrOriginalNameRequired = errors.New("请填写素材名称")
	// ErrMediaNotFound 表示素材不存在或已删除。
	ErrMediaNotFound = errors.New("素材不存在")
)

// StoredFile 是已写入存储的文件信息，不包含数据库记录。
type StoredFile struct {
	Kind         string
	OriginalName string
	Filename     string
	URL          string
	Mime         string
	Size         int64
	CreatedAt    int64
}

// FileUploader 是纯上传能力：校验并写入存储，不落库。
// 落库（写素材记录）由 Service 负责，两者可以分开装配。
type FileUploader interface {
	Upload(ctx context.Context, input UploadInput) (StoredFile, error)
	// Remove 删除已写入的文件，用于落库失败时回滚。
	Remove(filename string) error
}

// Limits 是上传文件的大小上限（字节）。
type Limits struct {
	MaxVideoBytes int64
	MaxImageBytes int64
}

// UploadInput 是上传素材入参；文件内容由 handler 从 multipart 中读出后传入。
type UploadInput struct {
	Kind         string
	OriginalName string
	Mime         string
	Size         int64
	Reader       io.Reader
}

// Uploader 是「上传并落库」的完整入口：写入存储后还要写一条素材记录。
// 只需要文件上传、不需要素材库记录的场景请用 FileUploader。
type Uploader interface {
	Upload(ctx context.Context, input UploadInput, operatorID *uint64) (*models.Media, error)
}

// ListInput 是素材列表查询条件；Kind 为空表示全部类型。
type ListInput struct {
	Kind     string
	Keyword  string
	Page     int
	PageSize int
}

// UpdateInput 是修改素材信息的入参。
type UpdateInput struct {
	OriginalName string
}

// 列表、详情、重命名与删除属于直接业务逻辑，不抽接口：
// 需要时由 impl.Service 直接提供，控制器调用具体实现即可。
