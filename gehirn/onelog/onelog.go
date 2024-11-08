/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/11/1 -- 18:28
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: onelog.go wrapper of zerolog
*/

package onelog

import (
	"context"
	"github.com/qiguanzhu/infra/nerv/magi/xenv"
	"github.com/rs/zerolog"
	"os"
	"sync"
)

var (
	_once         sync.Once
	defaultLogger zerolog.Logger

	Often          = zerolog.RandomSampler(10)
	Sometimes      = zerolog.RandomSampler(100)
	Rarely         = zerolog.RandomSampler(1000)
	VeryRarely     = zerolog.RandomSampler(10000)
	VeryVeryRarely = zerolog.RandomSampler(100000)
)

func Init(cfg xenv.Enver) error {
	_once.Do(func() {
		defaultLogger = zerolog.New(os.Stderr).With().Timestamp().Str("app", cfg.GetApp()).Str("env", cfg.GetEnv()).Logger()
		if cfg.IsProd() {
			defaultLogger = defaultLogger.Level(zerolog.InfoLevel)
		} else {
			defaultLogger = defaultLogger.Level(zerolog.DebugLevel)
		}
	})

	return nil
}

func Get() zerolog.Logger {
	return defaultLogger
}

func WithContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}

func Level(ctx context.Context, lvl zerolog.Level) context.Context {
	return zerolog.Ctx(ctx).Level(lvl).WithContext(ctx)
}

func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func Str(ctx context.Context, key string, val string) context.Context {
	return zerolog.Ctx(ctx).With().Str(key, val).Logger().WithContext(ctx)
}

func Strs(ctx context.Context, key string, val []string) context.Context {
	return zerolog.Ctx(ctx).With().Strs(key, val).Logger().WithContext(ctx)
}

func Int64(ctx context.Context, key string, val int64) context.Context {
	zerolog.Ctx(ctx).Info()
	return zerolog.Ctx(ctx).With().Int64(key, val).Logger().WithContext(ctx)
}

func Int64s(ctx context.Context, key string, val []int64) context.Context {
	return zerolog.Ctx(ctx).With().Ints64(key, val).Logger().WithContext(ctx)
}

func Bytes(ctx context.Context, key string, val []byte) context.Context {
	return zerolog.Ctx(ctx).With().Bytes(key, val).Logger().WithContext(ctx)
}

func Any(ctx context.Context, key string, val interface{}) context.Context {
	return zerolog.Ctx(ctx).With().Any(key, val).Logger().WithContext(ctx)
}

func Err(ctx context.Context, key string, val error) context.Context {
	return zerolog.Ctx(ctx).With().AnErr(key, val).Logger().WithContext(ctx)
}
