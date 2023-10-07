package logic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sybapp/infoflow/applications/article/api/internal/code"
	"github.com/sybapp/infoflow/applications/article/api/internal/svc"
	"github.com/sybapp/infoflow/applications/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 32 << 20 // 32MB

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *http.Request) (resp *types.UploadCoverResponse, err error) {
	err = req.ParseMultipartForm(maxFileSize)
	if err != nil {
		logx.Errorf("UploadCover maxFileSize: %vMB error: %v", maxFileSize, err)
		return nil, err
	}

	file, handler, err := req.FormFile("cover")
	if err != nil {
		logx.Errorf("UploadCover req: %v error: %v", req, err)
		return nil, err
	}
	defer file.Close()

	bucket, err := l.svcCtx.OssClient.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		logx.Errorf("GetBucket req: %v error: %v", req, err)
		return nil, code.GetBucketErr
	}

	objectKey := l.genFilename(handler.Filename)
	err = bucket.PutObject(objectKey, file)
	if err != nil {
		logx.Errorf("PutObject req: %v error: %v", req, err)
		return nil, code.PutBucketErr
	}
	return &types.UploadCoverResponse{
		CoverUrl: fmt.Sprintf("https://%s/%s", l.svcCtx.Config.Oss.Endpoint, objectKey),
	}, nil
}

func (l *UploadCoverLogic) genFilename(filename string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filename)
}
