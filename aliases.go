package startergin

import (
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 下面的类型别名把 Gin 的常用类型重新暴露到 startergin 包。
// 业务代码只需要依赖 startergin，就能在需要时继续使用 Gin 的原生能力。
type Accounts = gin.Accounts
type ContextKeyType = gin.ContextKeyType
type CORSHandlerConfig = cors.Config
type Engine = gin.Engine
type Context = gin.Context
type Error = gin.Error
type ErrorType = gin.ErrorType
type HandlersChain = gin.HandlersChain
type HandlerFunc = gin.HandlerFunc
type H = gin.H
type IRouter = gin.IRouter
type IRoutes = gin.IRoutes
type LogFormatter = gin.LogFormatter
type LogFormatterParams = gin.LogFormatterParams
type LoggerConfig = gin.LoggerConfig
type Negotiate = gin.Negotiate
type OnlyFilesFS = gin.OnlyFilesFS
type OptionFunc = gin.OptionFunc
type Param = gin.Param
type Params = gin.Params
type RecoveryFunc = gin.RecoveryFunc
type ResponseWriter = gin.ResponseWriter
type RouteInfo = gin.RouteInfo
type RouterGroup = gin.RouterGroup
type RoutesInfo = gin.RoutesInfo
type Skipper = gin.Skipper

const (
	// 常用 MIME 类型常量，保持与 gin 包中的定义完全一致。
	MIMEJSON              = gin.MIMEJSON
	MIMEHTML              = gin.MIMEHTML
	MIMEXML               = gin.MIMEXML
	MIMEXML2              = gin.MIMEXML2
	MIMEPlain             = gin.MIMEPlain
	MIMEPOSTForm          = gin.MIMEPOSTForm
	MIMEMultipartPOSTForm = gin.MIMEMultipartPOSTForm
	MIMEYAML              = gin.MIMEYAML
	MIMEYAML2             = gin.MIMEYAML2
	MIMETOML              = gin.MIMETOML
	MIMEPROTOBUF          = gin.MIMEPROTOBUF
	MIMEBSON              = gin.MIMEBSON

	// Gin 在可信平台代理场景下使用的平台标识。
	PlatformGoogleAppEngine = gin.PlatformGoogleAppEngine
	PlatformCloudflare      = gin.PlatformCloudflare
	PlatformFlyIO           = gin.PlatformFlyIO

	// Gin 运行模式常量。
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode

	// Gin 在上下文、绑定和运行时模式中使用的标准 key。
	AuthProxyUserKey = gin.AuthProxyUserKey
	AuthUserKey      = gin.AuthUserKey
	BindKey          = gin.BindKey
	BodyBytesKey     = gin.BodyBytesKey
	ContextKey       = gin.ContextKey
	EnvGinMode       = gin.EnvGinMode
	Version          = gin.Version

	ContextRequestKey = gin.ContextRequestKey

	// Gin 错误分类常量，用于按错误类型过滤日志或响应。
	ErrorTypeBind    = gin.ErrorTypeBind
	ErrorTypeRender  = gin.ErrorTypeRender
	ErrorTypePrivate = gin.ErrorTypePrivate
	ErrorTypePublic  = gin.ErrorTypePublic
	ErrorTypeAny     = gin.ErrorTypeAny
)

// BasicAuth 返回 Gin 的基础认证中间件。
func BasicAuth(accounts Accounts) HandlerFunc {
	return gin.BasicAuth(accounts)
}

// BasicAuthForProxy 返回用于代理认证场景的基础认证中间件。
func BasicAuthForProxy(accounts Accounts, realm string) HandlerFunc {
	return gin.BasicAuthForProxy(accounts, realm)
}

// BasicAuthForRealm 返回带自定义 realm 的基础认证中间件。
func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc {
	return gin.BasicAuthForRealm(accounts, realm)
}

// Bind 把指定值绑定到请求上下文，通常用于复用 Gin 的绑定中间件。
func Bind(val any) HandlerFunc {
	return gin.Bind(val)
}

// CreateTestContext 创建测试用的 Gin Context 和 Engine。
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	return gin.CreateTestContext(w)
}

// CreateTestContextOnly 基于已有 Engine 创建测试用 Context。
func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context) {
	return gin.CreateTestContextOnly(w, r)
}

// CORS 根据 gin-contrib/cors 的配置创建跨域中间件。
func CORS(config CORSHandlerConfig) HandlerFunc {
	return cors.New(config)
}

// CustomRecovery 创建自定义 panic 恢复中间件。
func CustomRecovery(handle RecoveryFunc) HandlerFunc {
	return gin.CustomRecovery(handle)
}

// CustomRecoveryWithWriter 创建自定义 panic 恢复中间件，并把恢复日志写入 out。
func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) HandlerFunc {
	return gin.CustomRecoveryWithWriter(out, handle)
}

// Default 创建带 Logger 和 Recovery 的 Gin Engine。
func Default(opts ...OptionFunc) *Engine {
	return gin.Default(opts...)
}

// Dir 返回 Gin 使用的静态文件系统。
func Dir(root string, listDirectory bool) http.FileSystem {
	return gin.Dir(root, listDirectory)
}

// DisableBindValidation 关闭 Gin 默认的结构体验证器。
func DisableBindValidation() {
	gin.DisableBindValidation()
}

// DisableConsoleColor 关闭控制台彩色日志输出。
func DisableConsoleColor() {
	gin.DisableConsoleColor()
}

// EnableJsonDecoderDisallowUnknownFields 让 JSON 解码器拒绝未知字段。
func EnableJsonDecoderDisallowUnknownFields() {
	gin.EnableJsonDecoderDisallowUnknownFields()
}

// EnableJsonDecoderUseNumber 让 JSON 解码器使用 json.Number 保留数字精度。
func EnableJsonDecoderUseNumber() {
	gin.EnableJsonDecoderUseNumber()
}

// ErrorLogger 返回记录 Gin 错误的中间件。
func ErrorLogger() HandlerFunc {
	return gin.ErrorLogger()
}

// ErrorLoggerT 返回只记录指定错误类型的中间件。
func ErrorLoggerT(typ ErrorType) HandlerFunc {
	return gin.ErrorLoggerT(typ)
}

// ForceConsoleColor 强制开启控制台彩色日志输出。
func ForceConsoleColor() {
	gin.ForceConsoleColor()
}

// IsDebugging 返回当前 Gin 是否处于 debug 模式。
func IsDebugging() bool {
	return gin.IsDebugging()
}

// Logger 返回 Gin 默认访问日志中间件。
func Logger() HandlerFunc {
	return gin.Logger()
}

// LoggerWithConfig 使用自定义配置创建访问日志中间件。
func LoggerWithConfig(conf LoggerConfig) HandlerFunc {
	return gin.LoggerWithConfig(conf)
}

// LoggerWithFormatter 使用自定义格式化函数创建访问日志中间件。
func LoggerWithFormatter(f LogFormatter) HandlerFunc {
	return gin.LoggerWithFormatter(f)
}

// LoggerWithWriter 使用自定义输出目标创建访问日志中间件。
func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc {
	return gin.LoggerWithWriter(out, notlogged...)
}

// Mode 返回当前 Gin 运行模式。
func Mode() string {
	return gin.Mode()
}

// New 创建一个未默认挂载中间件的 Gin Engine。
func New(opts ...OptionFunc) *Engine {
	return gin.New(opts...)
}

// Recovery 返回 Gin 默认 panic 恢复中间件。
func Recovery() HandlerFunc {
	return gin.Recovery()
}

// RecoveryWithWriter 使用自定义输出目标创建 panic 恢复中间件。
func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) HandlerFunc {
	return gin.RecoveryWithWriter(out, recovery...)
}

// SetDebugPrintFunc 设置 Gin debug 日志输出函数。
func SetDebugPrintFunc(fn func(format string, values ...any)) {
	gin.DebugPrintFunc = fn
}

// SetDebugPrintRouteFunc 设置 Gin 路由注册日志输出函数。
func SetDebugPrintRouteFunc(fn func(httpMethod, absolutePath, handlerName string, nuHandlers int)) {
	gin.DebugPrintRouteFunc = fn
}

// SetDefaultErrorWriter 设置 Gin 默认错误日志输出目标。
func SetDefaultErrorWriter(writer io.Writer) {
	gin.DefaultErrorWriter = writer
}

// SetDefaultWriter 设置 Gin 默认日志输出目标。
func SetDefaultWriter(writer io.Writer) {
	gin.DefaultWriter = writer
}

// SetMode 设置 Gin 运行模式，例如 debug、release 或 test。
func SetMode(value string) {
	gin.SetMode(value)
}

// WrapF 把标准库 http.HandlerFunc 包装为 Gin 中间件。
func WrapF(f http.HandlerFunc) HandlerFunc {
	return gin.WrapF(f)
}

// WrapH 把标准库 http.Handler 包装为 Gin 中间件。
func WrapH(h http.Handler) HandlerFunc {
	return gin.WrapH(h)
}
