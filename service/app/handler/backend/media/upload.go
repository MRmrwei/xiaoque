package media

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
	mediaservice "suki/xqyun/app/services/media"
)

// Upload 上传素材（视频 / 图片）。
// POST /api/v1/admin/media/upload，multipart/form-data，字段名 file，可选 kind。
func (m *Controller) Upload(ctx context.Context, c *app.RequestContext) {
	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		response.Fail(c, response.CodeBadRequest, "请选择要上传的文件")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, response.CodeBadRequest, "读取上传文件失败，请重试")
		return
	}
	defer func() { _ = file.Close() }()

	service, err := newService(m.App())
	if err != nil {
		fail(ctx, c, err)
		return
	}

	item, err := service.Upload(
		ctx,
		mediaservice.UploadInput{
			Kind:         c.PostForm("kind"),
			OriginalName: fileHeader.Filename,
			Mime:         string(fileHeader.Header.Get("Content-Type")),
			Size:         fileHeader.Size,
			Reader:       file,
		},
		currentUserID(c),
	)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, toMediaItem(item))
}
