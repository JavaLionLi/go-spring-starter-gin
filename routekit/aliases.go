package routekit

import "github.com/gin-gonic/gin"

type Accounts = gin.Accounts
type ContextKeyType = gin.ContextKeyType
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

func New(opts ...OptionFunc) *Engine {
	return gin.New(opts...)
}

func SetMode(value string) {
	gin.SetMode(value)
}
