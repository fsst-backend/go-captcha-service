package server

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wenlng/go-captcha-service/internal/common"
	"github.com/wenlng/go-captcha-service/internal/config"
	"github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config2 "github.com/wenlng/go-captcha-service/internal/pkg/gocaptcha/config"

	"github.com/wenlng/go-captcha-service/internal/cache"
	"github.com/wenlng/go-captcha-service/proto"
)

func TestCacheServer(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	ttl := time.Duration(10) * time.Second
	cleanInt := time.Duration(30) * time.Second
	cacheClient, _ := cache.NewCacheManager(
		&cache.CacheMgrParams{
			Type:      cache.CacheTypeMemory,
			KeyPrefix: "TEST_CAPTCHA_DATA:",
			Ttl:       ttl,
			CleanInt:  cleanInt,
		},
	)
	defer cacheClient.Close()

	dc := &config.DynamicConfig{Config: config.DefaultConfig()}
	captDCfg := &config2.DynamicCaptchaConfig{Config: config2.DefaultConfig()}

	logger, err := zap.NewProduction()
	assert.NoError(t, err)

	captcha, err := gocaptcha.Setup(captDCfg)
	assert.NoError(t, err)

	svcCtx := &common.SvcContext{
		CacheMgr:      cacheClient,
		DynamicConfig: dc,
		Logger:        logger,
		Captcha:       captcha,
	}
	server := NewGoCaptchaServer(svcCtx)

	t.Run("GetData", func(t *testing.T) {
		req := &proto.GetDataRequest{
			Id: "slide-default",
		}
		resp, err := server.GetData(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("GetData_Miss", func(t *testing.T) {
		req := &proto.GetDataRequest{
			Id: "slide-default",
		}
		_, err := server.GetData(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		//assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}
