// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	objmodel "github.com/solikewind/happyeat/dal/model/object"
	"github.com/spaolacci/murmur3"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传对象（用于菜单图片）
func NewUploadObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadObjectLogic {
	return &UploadObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadObjectLogic) UploadObject(file multipart.File, header *multipart.FileHeader) (*types.UploadObjectReply, error) {
	if l.svcCtx.Cos == nil {
		return nil, errors.New("cos 未配置")
	}
	if header == nil {
		return nil, errors.New("上传文件不能为空")
	}

	hasher := murmur3.New64()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}
	hash := strconv.FormatUint(hasher.Sum64(), 10)
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("文件流不支持重复读取")
	}

	existed, err := l.svcCtx.Object.GetByHash(l.ctx, hash)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}
	if existed != nil {
		out := entObjectToType(existed)
		out.Url = signedURLOrRaw(l.ctx, l.svcCtx, existed.Key, out.Url)
		return &types.UploadObjectReply{Object: out}, nil
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	objectID := strconv.FormatInt(time.Now().UnixMilli(), 10)
	key := fmt.Sprintf("menu/%s/%s%s", time.Now().Format("20060102"), objectID, ext)
	_, err = l.svcCtx.Cos.Object.Put(l.ctx, key, file, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	})
	if err != nil {
		return nil, err
	}

	item, err := l.svcCtx.Object.Create(l.ctx, objmodel.CreateInput{
		Name:        header.Filename,
		Key:         key,
		URL:         strings.TrimRight(l.svcCtx.Cos.BucketURL, "/") + "/" + key,
		ContentType: header.Header.Get("Content-Type"),
		Suffix:      ext,
		Size:        header.Size,
		Hash:        hash,
	})
	if err != nil {
		return nil, err
	}
	out := entObjectToType(item)
	out.Url = signedURLOrRaw(l.ctx, l.svcCtx, item.Key, out.Url)
	return &types.UploadObjectReply{Object: out}, nil
}
