# LiveKit PSRPC库详解（修订版）

## PSRPC概述

PSRPC (Publish-Subscribe RPC)是LiveKit开发的一个分布式RPC框架，它基于发布/订阅模式构建，专为分布式系统设计。与传统RPC不同，PSRPC允许消息通过主题(topic)进行路由，支持一对一、一对多和多对多通信模式。

## 核心概念和组件

根据`roommanager.psrpc.go`的实现和您的确认，PSRPC的核心组件包括：

### 1. MessageBus (消息总线)

消息总线是PSRPC的核心组件，负责处理所有消息的发布和订阅：

```go
// 从bus.go中定义
type MessageBus interface {
    // 发布消息
    Publish(topic string, data []byte) error
    // 订阅主题
    Subscribe(topic string, handler func(data []byte)) (Subscription, error)
    // 关闭连接
    Close() error
    // ...其他方法
}
```

在LiveKit中，MessageBus主要有两种实现：
- **Redis实现**：使用Redis的Pub/Sub功能
- **本地内存实现**：用于单节点或测试场景

### 2. 客户端和服务端接口

从`roommanager.psrpc.go`中可以看到，PSRPC生成了类型安全的客户端和服务端接口：

```go
// 客户端接口
type RoomManagerClient[NodeIdTopicType ~string] interface {
    CreateRoom(ctx context.Context, nodeId NodeIdTopicType, req *livekit6.CreateRoomRequest, opts ...psrpc.RequestOption) (*livekit1.Room, error)
    // ...其他方法
    Close()
}

// 服务端接口
type RoomManagerServer[NodeIdTopicType ~string] interface {
    RegisterCreateRoomTopic(nodeId NodeIdTopicType) error
    DeregisterCreateRoomTopic(nodeId NodeIdTopicType)
    // ...其他注册方法
    Shutdown()
    Kill()
}
```

### 3. 主题路由

PSRPC使用主题进行消息路由，允许按节点ID、房间ID或其他标识符路由请求：

```go
// 从roommanager.proto中
rpc CreateRoom(livekit.CreateRoomRequest) returns (livekit.Room) {
    option (psrpc.options) = {
        topics: true
        topic_params: {
            group: "node"  // 按节点分组
            names: ["node_id"]  // 路由参数
            typed: true  // 类型安全
        };
    };
};
```

## PSRPC的工作流程

基于`roommanager.psrpc.go`和您提供的信息，以下是一个典型的工作流程：

1. **创建服务实现**：
   ```go
   // 实现服务接口
   type myRoomManagerService struct {
       // ...
   }
   
   func (s *myRoomManagerService) CreateRoom(ctx context.Context, req *livekit.CreateRoomRequest) (*livekit.Room, error) {
       // 实现创建房间逻辑
       return &livekit.Room{...}, nil
   }
   ```

2. **初始化服务器**：
   ```go
   // 创建Redis客户端
   redisClient, err := createRedisClient(redisConfig)
   if err != nil {
       log.Fatal(err)
   }
   
   // 创建消息总线
   messageBus := psrpc.NewRedisMessageBus(redisClient)
   
   // 创建服务实现
   service := &myRoomManagerService{}
   
   // 创建服务器
   server, err := rpc.NewTypedRoomManagerServer(service, messageBus, psrpc.WithServerID("server-1"))
   if err != nil {
       log.Fatal(err)
   }
   
   // 注册主题
   err = server.RegisterCreateRoomTopic(livekit.NodeID("node-1"))
   if err != nil {
       log.Fatal(err)
   }
   ```

3. **创建客户端**：
   ```go
   // 创建类型安全的客户端
   client, err := rpc.NewTypedRoomManagerClient(messageBus)
   if err != nil {
       log.Fatal(err)
   }
   
   // 调用RPC方法
   room, err := client.CreateRoom(ctx, livekit.NodeID("node-1"), &livekit.CreateRoomRequest{
       Name: "my-room",
   })
   ```

4. **消息流**：
   - 客户端将请求发布到对应主题的Redis通道
   - 订阅该主题的服务器从Redis接收并处理请求
   - 服务器将响应发布到响应主题的Redis通道
   - 客户端从响应Redis通道接收结果

## PSRPC的高级特性

### 1. RPC类型

PSRPC支持多种RPC类型：

- **单一RPC** - 一个请求，一个响应
- **多重RPC** - 一个请求，多个响应
- **流式RPC** - 双向流式通信
- **亲和性RPC** - 服务器按亲和度选择处理请求

### 2. 拦截器

PSRPC提供了拦截器机制用于实现中间件功能：

```go
// 服务器拦截器
func loggingInterceptor(ctx context.Context, req proto.Message, info psrpc.RPCInfo, handler psrpc.ServerRPCHandler) (proto.Message, error) {
    log.Printf("Received: %s", info.Method)
    return handler(ctx, req)
}

// 添加到服务器
server, err := rpc.NewTypedRoomManagerServer(
    service, 
    messageBus, 
    psrpc.WithServerRPCInterceptors(loggingInterceptor),
)
```

### 3. 错误处理

PSRPC定义了自己的错误类型，可以携带错误代码和上下文：

```go
// 创建PSRPC错误
err := psrpc.NewError(psrpc.InvalidArgument, fmt.Errorf("invalid room name"))

// 检查错误类型
if psrpcErr, ok := err.(psrpc.Error); ok {
    code := psrpcErr.Code()
    // 处理特定错误
}
```

### 4. 服务选择

对于亲和性RPC，客户端可以定义服务选择选项：

```go
opts := psrpc.SelectionOpts{
    MinimumAffinity: 0.5,
    AffinityTimeout: time.Second,
}

room, err := client.CreateRoom(ctx, nodeID, req, psrpc.WithSelectionOpts(opts))
```

## 在LiveKit中的应用

在LiveKit中，PSRPC主要用于实现内部服务通信，如：

1. **节点间RPC** - 实现不同节点间的协调
2. **房间管理** - 创建、更新和删除房间
3. **参与者管理** - 处理参与者加入/离开等事件
4. **媒体处理** - 协调SFU（选择性转发单元）节点

从`roommanager.psrpc.go`可以看出，LiveKit使用泛型和类型别名增强了类型安全，确保只有正确类型的主题ID才能用于特定RPC调用。

## Redis作为消息总线的特点

LiveKit中使用Redis作为消息总线有以下特点：

1. **复用现有基础设施** - 利用已有的Redis实例，无需额外组件
2. **简单可靠** - Redis的Pub/Sub机制易于理解和维护
3. **与状态存储集成** - 同一Redis实例可同时用于状态存储和消息传递
4. **适合中小规模部署** - 对于大多数部署场景性能足够

实现方式：
```go
// 从LiveKit服务器代码
func getMessageBus(rc redis2.UniversalClient) psrpc.MessageBus {
    if rc == nil {
        return psrpc.NewLocalMessageBus()
    }
    return psrpc.NewRedisMessageBus(rc)
}
```

## 使用建议

如果您正在定制LiveKit，以下是使用PSRPC的一些建议：

1. **使用类型安全的包装器** - 优先使用`typed_api.go`中定义的类型安全客户端和服务器
2. **遵循现有模式** - 观察LiveKit如何组织RPC定义和实现
3. **理解主题路由** - 确保您理解主题如何影响消息路由
4. **利用错误处理** - 使用PSRPC的错误类型传递有意义的错误信息
5. **监控Redis性能** - 关注Redis的Pub/Sub性能指标，特别是在高负载场景

## 总结

PSRPC是LiveKit的分布式通信基础设施，它提供了：

1. **灵活的通信模式** - 支持各种RPC类型满足不同需求
2. **主题路由** - 允许消息准确路由到正确的处理程序
3. **类型安全** - 使用泛型确保API使用正确
4. **中间件支持** - 通过拦截器实现横切关注点
5. **错误处理** - 丰富的错误上下文和传播机制

LiveKit中的PSRPC实现基于Redis和本地内存，而非NATS。这种设计与LiveKit的其他组件无缝集成，为分布式通信提供了可靠的基础。

通过理解PSRPC的设计和使用模式，您可以更有效地定制和扩展LiveKit，而不需要修改底层PSRPC库本身。
