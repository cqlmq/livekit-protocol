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

package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/livekit/protocol/xtls"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/livekit/protocol/logger"
)

var ErrNotConfigured = errors.New("Redis is not configured")

type RedisConfig struct {
	Address           string       `yaml:"address,omitempty"`              // 地址
	Username          string       `yaml:"username,omitempty"`             // 用户名
	Password          string       `yaml:"password,omitempty"`             // 密码
	DB                int          `yaml:"db,omitempty"`                   // 数据库
	TLS               *xtls.Config `yaml:"tls,omitempty"`                  // 配置TLS
	MasterName        string       `yaml:"sentinel_master_name,omitempty"` // 哨兵主节点名称
	SentinelUsername  string       `yaml:"sentinel_username,omitempty"`    // 哨兵用户名
	SentinelPassword  string       `yaml:"sentinel_password,omitempty"`    // 哨兵密码
	SentinelAddresses []string     `yaml:"sentinel_addresses,omitempty"`   // 哨兵地址
	ClusterAddresses  []string     `yaml:"cluster_addresses,omitempty"`    // 集群地址
	DialTimeout       int          `yaml:"dial_timeout,omitempty"`         // 连接超时时间
	ReadTimeout       int          `yaml:"read_timeout,omitempty"`         // 读取超时时间
	WriteTimeout      int          `yaml:"write_timeout,omitempty"`        // 写入超时时间
	// for clustererd mode only, number of redirects to follow, defaults to 2
	MaxRedirects *int          `yaml:"max_redirects,omitempty"` // 集群模式下，重定向次数，默认2
	PoolTimeout  time.Duration `yaml:"pool_timeout,omitempty"`  // 池超时时间
	PoolSize     int           `yaml:"pool_size,omitempty"`     // 池大小

	// Deprecated: use TLS instead of UseTLS
	UseTLS bool `yaml:"use_tls,omitempty"` // 是否使用TLS
}

// IsConfigured 判断是否配置了Redis
func (r *RedisConfig) IsConfigured() bool {
	if r.Address != "" {
		return true
	}
	if len(r.SentinelAddresses) > 0 {
		return true
	}
	if len(r.ClusterAddresses) > 0 {
		return true
	}
	return false
}

// GetMaxRedirects 获取最大重定向次数
func (r *RedisConfig) GetMaxRedirects() int {
	if r.MaxRedirects != nil {
		return *r.MaxRedirects
	}
	return 2
}

// GetRedisClient 获取Redis客户端
func GetRedisClient(conf *RedisConfig) (redis.UniversalClient, error) {
	if conf == nil {
		return nil, nil
	}

	if !conf.IsConfigured() {
		return nil, ErrNotConfigured
	}

	var rcOptions *redis.UniversalOptions
	var rc redis.UniversalClient
	var tlsConfig *tls.Config

	if conf.TLS != nil && conf.TLS.Enabled {
		var err error
		tlsConfig, err = conf.TLS.ClientTLSConfig()
		if err != nil {
			return nil, err
		}
	} else if conf.UseTLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if len(conf.SentinelAddresses) > 0 {
		logger.Infow("connecting to redis", "sentinel", true, "addr", conf.SentinelAddresses, "masterName", conf.MasterName)

		// By default DialTimeout set to 2s
		if conf.DialTimeout == 0 {
			conf.DialTimeout = 2000
		}
		// By default ReadTimeout set to 0.2s
		if conf.ReadTimeout == 0 {
			conf.ReadTimeout = 200
		}
		// By default WriteTimeout set to 0.2s
		if conf.WriteTimeout == 0 {
			conf.WriteTimeout = 200
		}

		rcOptions = &redis.UniversalOptions{
			Addrs:            conf.SentinelAddresses,
			SentinelUsername: conf.SentinelUsername,
			SentinelPassword: conf.SentinelPassword,
			MasterName:       conf.MasterName,
			Username:         conf.Username,
			Password:         conf.Password,
			DB:               conf.DB,
			TLSConfig:        tlsConfig,
			DialTimeout:      time.Duration(conf.DialTimeout) * time.Millisecond,
			ReadTimeout:      time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout:     time.Duration(conf.WriteTimeout) * time.Millisecond,
			PoolTimeout:      conf.PoolTimeout,
			PoolSize:         conf.PoolSize,
		}
	} else if len(conf.ClusterAddresses) > 0 {
		logger.Infow("connecting to redis", "cluster", true, "addr", conf.ClusterAddresses)
		rcOptions = &redis.UniversalOptions{
			Addrs:        conf.ClusterAddresses,
			Username:     conf.Username,
			Password:     conf.Password,
			DB:           conf.DB,
			TLSConfig:    tlsConfig,
			MaxRedirects: conf.GetMaxRedirects(),
			PoolTimeout:  conf.PoolTimeout,
			PoolSize:     conf.PoolSize,
		}
	} else {
		logger.Infow("connecting to redis", "simple", true, "addr", conf.Address)
		rcOptions = &redis.UniversalOptions{
			Addrs:       []string{conf.Address},
			Username:    conf.Username,
			Password:    conf.Password,
			DB:          conf.DB,
			TLSConfig:   tlsConfig,
			PoolTimeout: conf.PoolTimeout,
			PoolSize:    conf.PoolSize,
		}
	}
	rc = redis.NewUniversalClient(rcOptions)

	if err := rc.Ping(context.Background()).Err(); err != nil {
		err = errors.Wrap(err, "unable to connect to redis")
		return nil, err
	}

	return rc, nil
}
