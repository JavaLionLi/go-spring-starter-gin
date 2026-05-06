package startergin

import (
	"io"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	PlatformGoogleAppEngine = gin.PlatformGoogleAppEngine
	PlatformCloudflare      = gin.PlatformCloudflare
	PlatformFlyIO           = gin.PlatformFlyIO

	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode

	AuthProxyUserKey = gin.AuthProxyUserKey
	AuthUserKey      = gin.AuthUserKey
	BindKey          = gin.BindKey
	BodyBytesKey     = gin.BodyBytesKey
	ContextKey       = gin.ContextKey
	EnvGinMode       = gin.EnvGinMode
	Version          = gin.Version

	ContextRequestKey = gin.ContextRequestKey

	ErrorTypeBind    = gin.ErrorTypeBind
	ErrorTypeRender  = gin.ErrorTypeRender
	ErrorTypePrivate = gin.ErrorTypePrivate
	ErrorTypePublic  = gin.ErrorTypePublic
	ErrorTypeAny     = gin.ErrorTypeAny
)

func BasicAuth(accounts Accounts) HandlerFunc {
	return gin.BasicAuth(accounts)
}

func BasicAuthForProxy(accounts Accounts, realm string) HandlerFunc {
	return gin.BasicAuthForProxy(accounts, realm)
}

func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc {
	return gin.BasicAuthForRealm(accounts, realm)
}

func Bind(val any) HandlerFunc {
	return gin.Bind(val)
}

func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	return gin.CreateTestContext(w)
}

func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context) {
	return gin.CreateTestContextOnly(w, r)
}

func CORS(config CORSHandlerConfig) HandlerFunc {
	return cors.New(config)
}

func CustomRecovery(handle RecoveryFunc) HandlerFunc {
	return gin.CustomRecovery(handle)
}

func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) HandlerFunc {
	return gin.CustomRecoveryWithWriter(out, handle)
}

func Default(opts ...OptionFunc) *Engine {
	return gin.Default(opts...)
}

func Dir(root string, listDirectory bool) http.FileSystem {
	return gin.Dir(root, listDirectory)
}

func DisableBindValidation() {
	gin.DisableBindValidation()
}

func DisableConsoleColor() {
	gin.DisableConsoleColor()
}

func EnableJsonDecoderDisallowUnknownFields() {
	gin.EnableJsonDecoderDisallowUnknownFields()
}

func EnableJsonDecoderUseNumber() {
	gin.EnableJsonDecoderUseNumber()
}

func ErrorLogger() HandlerFunc {
	return gin.ErrorLogger()
}

func ErrorLoggerT(typ ErrorType) HandlerFunc {
	return gin.ErrorLoggerT(typ)
}

func ForceConsoleColor() {
	gin.ForceConsoleColor()
}

func IsDebugging() bool {
	return gin.IsDebugging()
}

func Logger() HandlerFunc {
	return gin.Logger()
}

func LoggerWithConfig(conf LoggerConfig) HandlerFunc {
	return gin.LoggerWithConfig(conf)
}

func LoggerWithFormatter(f LogFormatter) HandlerFunc {
	return gin.LoggerWithFormatter(f)
}

func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc {
	return gin.LoggerWithWriter(out, notlogged...)
}

func Mode() string {
	return gin.Mode()
}

func New(opts ...OptionFunc) *Engine {
	return gin.New(opts...)
}

func Recovery() HandlerFunc {
	return gin.Recovery()
}

func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) HandlerFunc {
	return gin.RecoveryWithWriter(out, recovery...)
}

func SetDebugPrintFunc(fn func(format string, values ...any)) {
	gin.DebugPrintFunc = fn
}

func SetDebugPrintRouteFunc(fn func(httpMethod, absolutePath, handlerName string, nuHandlers int)) {
	gin.DebugPrintRouteFunc = fn
}

func SetDefaultErrorWriter(writer io.Writer) {
	gin.DefaultErrorWriter = writer
}

func SetDefaultWriter(writer io.Writer) {
	gin.DefaultWriter = writer
}

func SetMode(value string) {
	gin.SetMode(value)
}

func WrapF(f http.HandlerFunc) HandlerFunc {
	return gin.WrapF(f)
}

func WrapH(h http.Handler) HandlerFunc {
	return gin.WrapH(h)
}
