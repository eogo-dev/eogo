# Infrastructure 层详解

> `internal/infra/` 包含所有技术实现，是 Domain 层接口的具体实现。

## 目录结构

```
internal/infra/
├── config/        # 配置加载
├── database/      # 数据库连接
├── cache/         # 缓存抽象
├── redis/         # Redis 客户端
├── jwt/           # JWT 服务
├── email/         # 邮件服务
├── queue/         # 消息队列
├── middleware/    # HTTP 中间件
├── router/        # 路由封装
├── tracing/       # OpenTelemetry 追踪
├── metrics/       # Prometheus 指标
├── health/        # 健康检查
├── breaker/       # 熔断器
├── singleflight/  # 防缓存击穿
├── retry/         # 重试机制
├── ratelimit/     # 限流
└── testing/       # 测试工具
```

## config/ - 配置管理

```go
// internal/infra/config/config.go
type Config struct {
    App     AppConfig
    Server  ServerConfig
    DB      DatabaseConfig
    Redis   RedisConfig
    JWT     JWTConfig
    CORS    CORSConfig
    Tracing TracingConfig
}

type AppConfig struct {
    Name  string `env:"APP_NAME" envDefault:"eogo"`
    Env   string `env:"APP_ENV" envDefault:"development"`
    Debug bool   `env:"APP_DEBUG" envDefault:"true"`
}

type ServerConfig struct {
    Host         string `env:"SERVER_HOST" envDefault:""`
    Port         int    `env:"SERVER_PORT" envDefault:"8080"`
    Mode         string `env:"SERVER_MODE" envDefault:"debug"`
    ReadTimeout  int    `env:"SERVER_READ_TIMEOUT" envDefault:"10"`
    WriteTimeout int    `env:"SERVER_WRITE_TIMEOUT" envDefault:"10"`
}

// MustLoad 加载配置（失败则 panic）
func MustLoad() *Config {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        panic(err)
    }
    return cfg
}
```

## database/ - 数据库

```go
// internal/infra/database/database.go
func NewDB(cfg *config.Config) (*gorm.DB, func(), error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.DB.Username,
        cfg.DB.Password,
        cfg.DB.Host,
        cfg.DB.Port,
        cfg.DB.Database,
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, nil, err
    }

    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    cleanup := func() {
        sqlDB.Close()
    }

    return db, cleanup, nil
}
```

## jwt/ - JWT 服务

```go
// internal/infra/jwt/jwt.go
type Service struct {
    secret     string
    expiration time.Duration
}

func NewService(cfg *config.Config) *Service {
    return &Service{
        secret:     cfg.JWT.Secret,
        expiration: cfg.JWT.Expiration,
    }
}

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func (s *Service) Generate(userID uint, username string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secret))
}

func (s *Service) Validate(tokenString string) (uint, string, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.secret), nil
    })
    if err != nil {
        return 0, "", err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims.UserID, claims.Username, nil
    }

    return 0, "", errors.New("invalid token")
}
```

## middleware/ - HTTP 中间件

### JWT 认证中间件

```go
// internal/infra/middleware/jwt.go
var jwtService *jwt.Service

func SetJWTService(s *jwt.Service) {
    jwtService = s
}

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token == "" {
            response.Unauthorized(c, "Missing authorization token")
            c.Abort()
            return
        }

        userID, username, err := jwtService.Validate(token)
        if err != nil {
            response.Unauthorized(c, "Invalid token")
            c.Abort()
            return
        }

        c.Set("user_id", userID)
        c.Set("username", username)
        c.Next()
    }
}

func extractToken(c *gin.Context) string {
    auth := c.GetHeader("Authorization")
    if strings.HasPrefix(auth, "Bearer ") {
        return strings.TrimPrefix(auth, "Bearer ")
    }
    return c.Query("token")
}
```

### 其他中间件

```go
// CORS
func CORS() gin.HandlerFunc

// 请求日志
func Logger() gin.HandlerFunc

// 恢复
func Recovery() gin.HandlerFunc

// 限流
func RateLimit(max int, duration time.Duration) gin.HandlerFunc

// 请求超时
func Timeout(timeout time.Duration) gin.HandlerFunc

// 请求 ID
func RequestID() gin.HandlerFunc
```

## tracing/ - 分布式追踪

```go
// internal/infra/tracing/tracing.go
type TracerProvider struct {
    provider *sdktrace.TracerProvider
    tracer   trace.Tracer
}

func NewTracerProvider(cfg *Config) (*TracerProvider, error) {
    if !cfg.Enabled {
        return &TracerProvider{
            tracer: otel.Tracer(cfg.ServiceName),
        }, nil
    }

    // 创建 OTLP 导出器
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(cfg.Endpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }

    // 创建 TracerProvider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(cfg.ServiceName),
        )),
        sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)),
    )

    otel.SetTracerProvider(tp)
    return &TracerProvider{provider: tp, tracer: tp.Tracer(cfg.ServiceName)}, nil
}

// HTTP 中间件
func Middleware(serviceName string) gin.HandlerFunc {
    tracer := otel.Tracer(serviceName)
    return func(c *gin.Context) {
        ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
        defer span.End()

        span.SetAttributes(
            attribute.String("http.method", c.Request.Method),
            attribute.String("http.url", c.Request.URL.String()),
        )

        c.Request = c.Request.WithContext(ctx)
        c.Next()

        span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
    }
}

// GORM 插件
func WithTracing(db *gorm.DB, serviceName string) error {
    return db.Use(NewGormPlugin(serviceName))
}
```

## metrics/ - Prometheus 指标

```go
// internal/infra/metrics/metrics.go
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.FullPath()
        method := c.Request.Method

        c.Next()

        status := strconv.Itoa(c.Writer.Status())
        duration := time.Since(start).Seconds()

        httpRequestsTotal.WithLabelValues(method, path, status).Inc()
        httpRequestDuration.WithLabelValues(method, path).Observe(duration)
    }
}

func Handler() gin.HandlerFunc {
    h := promhttp.Handler()
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}
```

## health/ - 健康检查

```go
// internal/infra/health/health.go
type Health struct {
    checkers map[string]Checker
}

type Checker func(ctx context.Context) CheckResult

type CheckResult struct {
    Status  Status
    Message string
    Details map[string]any
}

func (h *Health) Register(name string, checker Checker) {
    h.checkers[name] = checker
}

// 内置检查器
func DatabaseChecker(db *gorm.DB) Checker {
    return func(ctx context.Context) CheckResult {
        sqlDB, _ := db.DB()
        if err := sqlDB.PingContext(ctx); err != nil {
            return CheckResult{Status: StatusDown, Message: "database ping failed"}
        }
        return CheckResult{Status: StatusUp}
    }
}

// HTTP 处理器
func (h *Health) Handler() gin.HandlerFunc
func LivenessHandler() gin.HandlerFunc   // /health/live
func (h *Health) ReadinessHandler() gin.HandlerFunc  // /health/ready
```

## breaker/ - 熔断器

```go
// internal/infra/breaker/breaker.go
type Breaker interface {
    Name() string
    Allow() (Promise, error)
    Do(req func() error) error
    DoWithFallback(req func() error, fallback Fallback) error
    State() State
}

type Config struct {
    Name                string
    Threshold           int           // 失败阈值
    Timeout             time.Duration // 熔断超时
    MaxHalfOpenRequests int           // 半开状态最大请求数
}

func New(cfg Config) Breaker {
    return &circuitBreaker{
        name:      cfg.Name,
        threshold: cfg.Threshold,
        timeout:   cfg.Timeout,
        state:     StateClosed,
    }
}

// 使用示例
breaker := breaker.New(breaker.Config{
    Name:      "payment-service",
    Threshold: 5,
    Timeout:   10 * time.Second,
})

err := breaker.Do(func() error {
    return paymentService.Charge(amount)
})
```

## singleflight/ - 防缓存击穿

```go
// internal/infra/singleflight/singleflight.go
type SingleFlight interface {
    Do(key string, fn func() (any, error)) (any, error)
    DoEx(key string, fn func() (any, error)) (val any, fresh bool, err error)
}

func New() SingleFlight {
    return &flightGroup{calls: make(map[string]*call)}
}

// 使用示例
sf := singleflight.New()

// 100 个并发请求，只有 1 个会真正查询数据库
result, err := sf.Do("user:123", func() (any, error) {
    return db.GetUser(123)
})
```

## retry/ - 重试机制

```go
// internal/infra/retry/retry.go
type Config struct {
    MaxAttempts  int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
    Jitter       float64
}

func Do(ctx context.Context, fn func(ctx context.Context) error, opts ...Option) error

// 使用示例
err := retry.Do(ctx, func(ctx context.Context) error {
    return externalService.Call()
},
    retry.WithMaxAttempts(3),
    retry.WithInitialDelay(100*time.Millisecond),
    retry.WithMultiplier(2.0),
)

// 便捷函数
retry.ExponentialBackoff(ctx, fn, 3)
retry.Times(ctx, 3, fn)
```

## ratelimit/ - 限流

```go
// internal/infra/ratelimit/limiter.go
type Config struct {
    Max      int
    Duration time.Duration
    KeyFunc  func(*gin.Context) string
}

func Middleware(cfg Config) gin.HandlerFunc

// 便捷函数
ratelimit.PerMinute(60)      // 每分钟 60 次
ratelimit.PerHour(1000)      // 每小时 1000 次
ratelimit.ByUser(100, time.Minute, getUserID)  // 按用户限流
ratelimit.ByRoute(30, time.Minute)  // 按路由限流
```

## 最佳实践

### DO ✅

- 实现 Domain 层定义的接口
- 使用配置驱动
- 提供清理函数（cleanup）
- 支持优雅关闭

### DON'T ❌

- 在 Infrastructure 层写业务逻辑
- 硬编码配置
- 忽略错误处理
- 忽略资源清理

## 下一步

- [业务模块开发](./03-business-modules.md)
- [可观测性](../production/01-observability.md)
