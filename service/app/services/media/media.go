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

// UploadInput 是上传素材入参；文件内容由 handler 从 multipart 中读出后传入。
type UploadInput struct {
	Kind         string
	OriginalName string
	Mime         string
	Size         int64
	Reader       io.Reader
}

// Uploader 是上传素材的统一入口。
// 不同上传逻辑（本地处理、转码、对象存储直传等）提供不同实现即可。
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

// Lister 是素材列表查询能力。
type Lister interface {
	List(ctx context.Context, input ListInput) ([]models.Media, int64, error)
}

// Manager 是素材详情、修改与删除能力。
type Manager interface {
	Detail(ctx context.Context, id uint64) (*models.Media, error)
	Update(ctx context.Context, id uint64, input UpdateInput) (*models.Media, error)
	Delete(ctx context.Context, id uint64) error
	DeleteMany(ctx context.Context, ids []uint64) error
}

// Service 是控制器使用的素材用例集合。
type Service interface {
	Uploader
	Lister
	Manager
}
