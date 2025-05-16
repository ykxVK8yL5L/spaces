package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"strconv"
	"strings"
	"sync"
	"time"

	pbPublicUser "github.com/city404/v6-public-rpc-proto/go/v6/user"
	"github.com/google/uuid"
	"github.com/halalcloud/golang-sdk/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	appID                string
	appVersion           string
	appSecret            string
	refreshToken         string
	accessToken          string
	accessTokenExpiredAt int64
	grpcConnection       *grpc.ClientConn
	dopts                halalOptions
	configMutex          sync.RWMutex
}

func NewAuthServiceWithOauth(appID, appVersion, appSecret string, options ...HalalOption) (string, error) {
	svc := &AuthService{
		appID:        appID,
		appVersion:   appVersion,
		appSecret:    appSecret,
		refreshToken: "",
		dopts:        defaultOptions(),
	}
	for _, opt := range options {
		opt.apply(&svc.dopts)
	}

	grpcServer := "grpcuserapi.2dland.cn:443"
	// dialContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	grpcOptions := svc.dopts.grpcOptions
	grpcOptions = append(grpcOptions, grpc.WithAuthority("grpcuserapi.2dland.cn"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})), grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctxx := svc.signContext(method, ctx)
		err := invoker(ctxx, method, req, reply, cc, opts...) // invoking RPC method
		return err
	}))

	grpcConnection, err := grpc.NewClient(grpcServer, grpcOptions...)
	if err != nil {
		return "发生错误", err
	}
	defer grpcConnection.Close()
	userClient := pbPublicUser.NewPubUserClient(grpcConnection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stateString := uuid.New().String()
	// queryValues.Add("callback", oauthToken.Callback)
	oauthToken, err := userClient.CreateAuthToken(ctx, &pbPublicUser.LoginRequest{
		ReturnType: 2,
		State:      stateString,
		ReturnUrl:  "",
	})
	if err != nil {
		return "发生错误", err
	}
	if len(oauthToken.State) < 1 {
		oauthToken.State = stateString
	}

	return oauthToken.Url, nil

}

func NewAuthService(appID, appVersion, appSecret, refreshToken string, accessToken string, accessTokenExpiredAt int64, options ...HalalOption) (*AuthService, error) {

	svc := &AuthService{
		appID:        appID,
		appVersion:   appVersion,
		appSecret:    appSecret,
		refreshToken: refreshToken,
		dopts:        defaultOptions(),
	}

	readedAccessToken := accessToken
	if len(readedAccessToken) > 0 {
		// svc.accessToken = readedAccessToken
		accessTokenExpiredAt := accessTokenExpiredAt
		current := time.Now().UnixMilli()
		if accessTokenExpiredAt < current {
			// access token expired
			svc.accessToken = ""
			svc.accessTokenExpiredAt = 0
			readedAccessToken = ""
			println("access token expired:", accessTokenExpiredAt)
		} else {
			svc.accessTokenExpiredAt = accessTokenExpiredAt
			svc.accessToken = readedAccessToken
		}

	}

	for _, opt := range options {
		opt.apply(&svc.dopts)
	}

	grpcServer := "grpcuserapi.2dland.cn:443"
	grpcOptions := svc.dopts.grpcOptions
	grpcOptions = append(grpcOptions, grpc.WithAuthority("grpcuserapi.2dland.cn"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})), grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// <!---- comment start ---->
		// check if accesstoken is expired, if expired, refresh it, this operation is not necessary for every request
		// it's just a demo, you should not do this in production environment
		// instead, you should refresh token when you get error code unauthenticated
		// or you can use a background goroutine/thread to refresh token periodically
		// thus it's not necessary, because the interceptor will refresh token automatically if token is expired
		//// ignoreAutoRefeshMethod := []string{pbPublicUser.PubUser_Login_FullMethodName, pbPublicUser.PubUser_Refresh_FullMethodName, pbPublicUser.PubUser_SendSmsVerifyCode_FullMethodName}
		////ignoreAutoRefesh := false
		////for _, m := range ignoreAutoRefeshMethod {
		////	if m == method {
		////		ignoreAutoRefesh = true
		////		break
		////	}
		////}
		////if !ignoreAutoRefesh && len(svc.accessToken) > 0 && accessTokenExpiredAt+120000 < time.Now().UnixMilli() && len(refreshToken) > 0 {
		////	// refresh token
		////	refreshResponse, err := pbPublicUser.NewPubUserClient(cc).Refresh(ctx, &pbPublicUser.Token{
		////		RefreshToken: refreshToken,
		////	})
		////	if err != nil {
		////		return err
		////	}
		////	if len(refreshResponse.AccessToken) > 0 {
		////		accessToken = refreshResponse.AccessToken
		////		accessTokenExpiredAt = refreshResponse.AccessTokenExpireTs
		////	}
		////}
		// <!---- comment end ---->
		// currentTimeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		ctxx := svc.signContext(method, ctx)
		err := invoker(ctxx, method, req, reply, cc, opts...) // invoking RPC method
		if err != nil {
			grpcStatus, ok := status.FromError(err)
			// if error is grpc error and error code is unauthenticated and error message contains "invalid accesstoken" and refresh token is not empty
			// then refresh access token and retry
			if ok && grpcStatus.Code() == codes.Unauthenticated && strings.Contains(grpcStatus.Err().Error(), "invalid accesstoken") && len(refreshToken) > 0 {
				// refresh token
				refreshResponse, err := pbPublicUser.NewPubUserClient(cc).Refresh(ctx, &pbPublicUser.Token{
					RefreshToken: refreshToken,
				})
				if err != nil {
					return err
				}
				if len(refreshResponse.AccessToken) > 0 {
					svc.accessToken = refreshResponse.AccessToken
					svc.accessTokenExpiredAt = refreshResponse.AccessTokenExpireTs
					//svc.OnAccessTokenRefreshed(refreshResponse.AccessToken, refreshResponse.AccessTokenExpireTs, refreshResponse.RefreshToken, refreshResponse.RefreshTokenExpireTs)
				}
				// retry
				ctxx := svc.signContext(method, ctx)
				err = invoker(ctxx, method, req, reply, cc, opts...) // invoking RPC method
				if err != nil {
					return err
				} else {
					return nil
				}
			}
		}
		// post-processing
		return err
	}))
	grpcConnection, err := grpc.NewClient(grpcServer, grpcOptions...)

	if err != nil {
		return nil, err
	}

	svc.grpcConnection = grpcConnection

	/*
		testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		refreshResponse, err := pbPublicUser.NewPubUserClient(svc.grpcConnection).Refresh(testCtx, &pbPublicUser.Token{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
		// if len(refreshResponse.AccessToken) > 0 {
		svc.OnAccessTokenRefreshed(refreshResponse.AccessToken, refreshResponse.AccessTokenExpireTs, refreshResponse.RefreshToken, refreshResponse.RefreshTokenExpireTs)
	*/

	return svc, err
}

func (s *AuthService) OnAccessTokenRefreshed(accessToken string, accessTokenExpiredAt int64, refreshToken string, refreshTokenExpiredAt int64) {
	// s.accessToken = accessToken
	// s.accessTokenExpiredAt = accessTokenExpiredAt
	s.configMutex.Lock()
	defer s.configMutex.Unlock()
	s.refreshToken = refreshToken
	if s.dopts.onTokenRefreshed != nil {
		s.dopts.onTokenRefreshed(accessToken, accessTokenExpiredAt, refreshToken, refreshTokenExpiredAt)
	}

}

func (s *AuthService) GetGrpcConnection() *grpc.ClientConn {
	return s.grpcConnection
}

type AuthInfo struct {
	AppID                string `json:"app_id"`
	AppVersion           string `json:"app_version"`
	AppSecret            string `json:"app_secret"`
	RefreshToken         string `json:"refresh_token"`
	AccessToken          string `json:"access_token"`
	AccessTokenExpiredAt int64  `json:"expire_at"`
}

func (s *AuthService) GetAuth() *AuthInfo {
	sai := &AuthInfo{
		AppID:                s.appID,
		AppVersion:           s.appVersion,
		AppSecret:            s.appSecret,
		RefreshToken:         s.refreshToken,
		AccessToken:          s.accessToken,
		AccessTokenExpiredAt: s.accessTokenExpiredAt,
	}
	return sai
}

func (s *AuthService) Close() {
	s.grpcConnection.Close()
}

func (s *AuthService) signContext(method string, ctx context.Context) context.Context {
	kvString := []string{}
	currentTimeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	bufferedString := bytes.NewBufferString(method)
	kvString = append(kvString, "timestamp", currentTimeStamp)
	bufferedString.WriteString(currentTimeStamp)
	kvString = append(kvString, "appid", s.appID)
	bufferedString.WriteString(s.appID)
	kvString = append(kvString, "appversion", s.appVersion)
	bufferedString.WriteString(s.appVersion)
	if len(s.accessToken) > 0 {
		authorization := "Bearer " + s.accessToken
		kvString = append(kvString, "authorization", authorization)
		bufferedString.WriteString(authorization)
	}
	bufferedString.WriteString(s.appSecret)
	sign := utils.GetMD5Hash(bufferedString.String())
	kvString = append(kvString, "sign", sign)
	return metadata.AppendToOutgoingContext(ctx, kvString...)
}
