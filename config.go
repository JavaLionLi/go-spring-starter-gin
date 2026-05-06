package startergin

import "time"

// Config controls the Gin engine auto-configured by this starter.
//
// It is bound from the "spring.gin" configuration prefix.
type Config struct {
	Mode           string       `value:"${mode:=release}"`
	TrustedProxies []string     `value:"${trusted-proxies:=}"`
	Logger         bool         `value:"${logger:=true}"`
	Recovery       bool         `value:"${recovery:=true}"`
	Health         HealthConfig `value:"${health}"`
	CORS           CORSConfig   `value:"${cors}"`
}

type HealthConfig struct {
	Enabled bool   `value:"${enabled:=true}"`
	Healthz string `value:"${healthz:=/healthz}"`
	Ping    string `value:"${ping:=/ping}"`
}

type CORSConfig struct {
	Enabled          bool          `value:"${enabled:=true}"`
	AllowOrigins     []string      `value:"${allow-origins:=*}"`
	AllowMethods     []string      `value:"${allow-methods:=GET,POST,PUT,PATCH,DELETE,OPTIONS}"`
	AllowHeaders     []string      `value:"${allow-headers:=Origin,Content-Type,Accept,Authorization}"`
	ExposeHeaders    []string      `value:"${expose-headers:=Content-Length}"`
	AllowCredentials bool          `value:"${allow-credentials:=false}"`
	MaxAge           time.Duration `value:"${max-age:=12h}"`
}
