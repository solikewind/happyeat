package cos

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/solikewind/happyeat/app/internal/config"
	cossdk "github.com/tencentyun/cos-go-sdk-v5"
)

type Client struct {
	BucketURL string
	secretID  string
	secretKey string
	*cossdk.Client
}

func NewClient(conf config.Cos) *Client {
	bucket := ""
	if conf.BucketUrl != nil {
		bucket = *conf.BucketUrl
	}
	if bucket == "" {
		return nil
	}
	u, err := url.Parse(bucket)
	if err != nil {
		return nil
	}
	secretID := ""
	if conf.SecretId != nil {
		secretID = *conf.SecretId
	}
	secretKey := ""
	if conf.SecretKey != nil {
		secretKey = *conf.SecretKey
	}
	b := &cossdk.BaseURL{BucketURL: u}
	c := cossdk.NewClient(b, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cossdk.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return &Client{
		BucketURL: bucket,
		secretID:  secretID,
		secretKey: secretKey,
		Client:    c,
	}
}

func (c *Client) PresignedGetURL(ctx context.Context, key string, expired time.Duration) (string, error) {
	if c == nil {
		return "", nil
	}
	if expired <= 0 {
		expired = 10 * time.Minute
	}
	u, err := c.Object.GetPresignedURL(ctx, http.MethodGet, key, c.secretID, c.secretKey, expired, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
