package part

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anden007/dp_clean_core/misc"
	"github.com/anden007/dp_clean_core/pkg"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/spf13/viper"
)

type IJwtService interface {
	// CreateToken 生成Token，默认存入Cookie中
	CreateToken(ctx iris.Context, claims interface{}) (token string, err error)
	// GetClaims 验证Token并返回其中包含的Claims对象指针，需要在Middleware处理后才能获得，如需手动验证，请使用VerifyToken方法，强烈建议使用此方法
	GetClaims(ctx iris.Context) (claims interface{}, err error)
	// VerifyToken 从默认渠道(Cookie)中获取Token信息并验证，验证成功后返回VerifiedToken对象，仅用于手动验证，不建议使用
	VerifyToken(ctx iris.Context) (result *jwt.VerifiedToken, err error)
	// RemoveToken 从Cookie中移除Token信息，一般用于"注销"功能
	RemoveToken(ctx iris.Context)
	// GetMiddleware 获取iris路由中间件，用于验证请求合法性
	GetMiddleware() iris.Handler
	// 获取用户JWT对象
	GetExecutorInfo(ctx iris.Context) (result pkg.BaseUserInfo, err error)
	// 获取用户请求Context
	GetExecutorContext(ctx iris.Context) context.Context
	// 获取用户请求Context
	GetExecutorContextWithInfo(ctx iris.Context) (result context.Context, info *pkg.BaseUserInfo)
	// 从cookie中获取Token
	VerifyCookieToken(token string) (result *jwt.VerifiedToken, err error)
}

type JWTService struct {
	jwtHeader        string
	jwtSecretKey     string
	jwtHeaderType    string
	jwtMaxAgeSeconds int
	jwtSigner        *jwt.Signer
	jwtVerifier      *jwt.Verifier
}

func NewJWT() IJwtService {
	return NewCustomJWT("")
}

func NewCustomJWT(header string) *JWTService {
	instance := new(JWTService)
	if header == "" {
		header = viper.GetString("jwt.header")
	}
	instance.jwtHeader = header
	instance.jwtSecretKey = viper.GetString("jwt.secret_key")
	instance.jwtHeaderType = viper.GetString("jwt.header_type")
	instance.jwtMaxAgeSeconds = viper.GetInt("jwt.max_age_seconds")
	instance.jwtSigner = jwt.NewSigner(jwt.HS256, instance.jwtSecretKey, time.Second*time.Duration(instance.jwtMaxAgeSeconds))
	instance.jwtVerifier = jwt.NewVerifier(jwt.HS256, instance.jwtSecretKey, jwt.Expected{
		Issuer: APP_NAME,
	})
	instance.jwtVerifier.WithDefaultBlocklist()
	instance.jwtVerifier.Extractors = append(instance.jwtVerifier.Extractors, instance.FromHeader, instance.FromCookie)
	return instance
}

// CreateToken 生成Token，默认存入Cookie中
func (m *JWTService) CreateToken(ctx iris.Context, claims interface{}) (token string, err error) {
	/* 这里省去了对用户的验证，在实际使用过程中需要验证用户是否存在，密码是否正确 */
	if tokenData, err := m.jwtSigner.Sign(claims, jwt.Claims{
		ID:     misc.NewGuidString(),
		Issuer: APP_NAME,
	}); err == nil {
		token = string(tokenData)
		ctx.SetCookieKV(m.jwtHeader, fmt.Sprintf("%s %s", m.jwtHeaderType, token), iris.CookieHTTPOnly(false), iris.CookieExpires(time.Second*time.Duration(m.jwtMaxAgeSeconds)))
	}
	return
}

// VerifyToken 从默认渠道(Cookie)中获取Token信息并验证，验证成功后返回VerifiedToken对象，仅用于手动验证，不建议使用
func (m *JWTService) VerifyToken(ctx iris.Context) (result *jwt.VerifiedToken, err error) {
	token_str := ""
	header_token_str := m.FromHeader(ctx)
	cookie_token_str := m.FromCookie(ctx)
	if header_token_str != "" {
		token_str = header_token_str
	} else if cookie_token_str != "" {
		token_str = cookie_token_str
	}
	result, err = m.jwtVerifier.VerifyToken([]byte(token_str))
	return
}

func (m *JWTService) VerifyCookieToken(token string) (result *jwt.VerifiedToken, err error) {
	tokenCookieParts := strings.Split(token, " ")
	if len(tokenCookieParts) != 2 || !strings.EqualFold(tokenCookieParts[0], m.jwtHeaderType) {
		return nil, nil
	}
	tokenStr := tokenCookieParts[1]
	result, err = m.jwtVerifier.VerifyToken([]byte(tokenStr))
	return
}

// RemoveToken 从Cookie中移除Token信息，一般用于"注销"功能
func (m *JWTService) RemoveToken(ctx iris.Context) {
	ctx.Logout()
	ctx.RemoveCookie(m.jwtHeader)
}

// GetClaims 验证Token并返回其中包含的Claims对象指针，需要在Middleware处理后才能获得，如需手动验证，请使用VerifyToken方法，强烈建议使用此方法
func (m *JWTService) GetClaims(ctx iris.Context) (claims interface{}, err error) {
	claims = jwt.Get(ctx)
	return
}

// GetMiddleware 获取iris路由中间件，用于验证请求合法性
func (m *JWTService) GetMiddleware() iris.Handler {
	return m.jwtVerifier.Verify(func() interface{} {
		return new(pkg.BaseUserInfoClaims)
	})
}

func (m *JWTService) FromCookie(ctx iris.Context) string {
	tokenCookie := ctx.GetCookie(m.jwtHeader)
	if tokenCookie == "" {
		return ""
	}

	tokenCookieParts := strings.Split(tokenCookie, " ")
	if len(tokenCookieParts) != 2 || !strings.EqualFold(tokenCookieParts[0], m.jwtHeaderType) {
		return ""
	}

	return tokenCookieParts[1]
}

func (m *JWTService) FromHeader(ctx iris.Context) string {
	return ctx.Request().Header.Get(m.jwtHeader)
}

func (m *JWTService) GetExecutorInfo(ctx iris.Context) (result pkg.BaseUserInfo, err error) {
	if ctx != nil {
		if claims, cErr := m.GetClaims(ctx); claims != nil {
			userInfoClaims := claims.(*pkg.BaseUserInfoClaims)
			result = userInfoClaims.BaseUserInfo
		} else {
			err = cErr
		}
	}
	return
}

func (m *JWTService) GetExecutorContext(ctx iris.Context) context.Context {
	if executorInfo, eErr := m.GetExecutorInfo(ctx); eErr == nil {
		return context.WithValue(context.TODO(), "executor", executorInfo)
	}
	return context.TODO()
}

func (m *JWTService) GetExecutorContextWithInfo(ctx iris.Context) (result context.Context, info *pkg.BaseUserInfo) {
	if executorInfo, eErr := m.GetExecutorInfo(ctx); eErr == nil {
		return context.WithValue(context.TODO(), "executor", executorInfo), &executorInfo
	}
	return context.TODO(), nil
}
