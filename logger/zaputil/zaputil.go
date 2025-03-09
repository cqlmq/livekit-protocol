// Copyright 2023 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zaputil

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DateFormat = "01-02 15:04:05"
)

func encoderWithValues(enc zapcore.Encoder, kvs ...any) zapcore.Encoder {
	clone := enc.Clone()
	for i := 1; i < len(kvs); i += 2 {
		if key, ok := kvs[i-1].(string); ok {
			zap.Any(key, kvs[i]).AddTo(clone)
		}
	}
	return clone
}

type Encoder[T any] interface {
	WithValues(kvs ...any) T
	Core(console, json *WriteEnabler) zapcore.Core
}

type DevelopmentEncoder struct {
	console zapcore.Encoder
	json    zapcore.Encoder
}

func NewDevelopmentEncoder() DevelopmentEncoder {
	devConfig := zap.NewDevelopmentEncoderConfig()
	// 修改时间格式
	devConfig.EncodeTime = zapcore.TimeEncoderOfLayout(DateFormat)
	// 设置日志级别显示的颜色
	devConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 或者使用其他颜色编码器:
	// - zapcore.LowercaseColorLevelEncoder: 小写带颜色
	// - zapcore.CapitalColorLevelEncoder: 大写带颜色
	// - zapcore.LowercaseLevelEncoder: 小写不带颜色
	// - zapcore.CapitalLevelEncoder: 大写不带颜色

	prodConfig := zap.NewProductionEncoderConfig()
	prodConfig.EncodeTime = devConfig.EncodeTime
	// 生产环境通常不需要颜色
	prodConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return DevelopmentEncoder{
		console: zapcore.NewConsoleEncoder(devConfig),
		json:    zapcore.NewJSONEncoder(prodConfig),
	}
}

func (e DevelopmentEncoder) WithValues(kvs ...any) DevelopmentEncoder {
	e.console = encoderWithValues(e.console, kvs...)
	e.json = encoderWithValues(e.json, kvs...)
	return e
}

func (e DevelopmentEncoder) Core(console, json *WriteEnabler) zapcore.Core {
	return zapcore.NewTee(NewEncoderCore(e.console, console), NewEncoderCore(e.json, json))
}

type ProductionEncoder struct {
	json zapcore.Encoder
}

func NewProductionEncoder() ProductionEncoder {
	config := zap.NewProductionEncoderConfig()
	// 修改时间格式
	config.EncodeTime = zapcore.TimeEncoderOfLayout(DateFormat)
	// 或者使用 ISO8601 格式
	// config.EncodeTime = zapcore.ISO8601TimeEncoder

	return ProductionEncoder{
		json: zapcore.NewJSONEncoder(config),
	}
}

func (e ProductionEncoder) WithValues(kvs ...any) ProductionEncoder {
	e.json = encoderWithValues(e.json, kvs...)
	return e
}

func (e ProductionEncoder) Core(console, json *WriteEnabler) zapcore.Core {
	return NewEncoderCore(e.json, console, json)
}
