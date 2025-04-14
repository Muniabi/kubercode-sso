package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/repository"
	"kubercode-sso/internal/infrastructure/utils"
	sso "kubercode-sso/proto/pb/go"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AbstractTokenService interface {
	RotateToken(ctx context.Context, token *sso.RefreshToken) (*sso.PairTokens, error)
	TokensIssue(ctx context.Context, subjectEmail string, subjectId uuid.UUID,
		deviceId string) (*sso.PairTokens, error)
	TokenRevoke(ctx context.Context, subjectEmail string, deviceId uuid.UUID) (*sso.VerifyTokenResponse, error)
	VerifyToken(ctx context.Context, token string) (jwt.MapClaims, error)
}

type DataToEncode struct {
	TokenType string `json:"token_type"`
	Exp       int    `json:"exp"`
	Iat       int    `json:"iat"`
	Jti       string `json:"jti"`
	DeviceId  string `json:"deviceId"`
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	SubjectId string `json:"subjectId"`
	IsMentor  bool   `json:"isMentor"`
}

type Claims struct {
	jwt.RegisteredClaims
	SubjectEmail string    `json:"subjectEmail"`
	DeviceId     string    `json:"deviceId"`
	SubjectId    uuid.UUID `json:"subjectId"`
	IsMentor     bool      `json:"isMentor"`
}

type ServiceJWT struct {
	cfg             *config.Config
	log             *slog.Logger
	tokenRepository repository.AbstractTokenRepository
	redisClient     *redis.Client
}

func NewJWTService(cfg *config.Config, log *slog.Logger, tokenRepository repository.AbstractTokenRepository,
	client *redis.Client) *ServiceJWT {
	return &ServiceJWT{
		cfg:             cfg,
		log:             log,
		tokenRepository: tokenRepository,
		redisClient:     client,
	}
}

func (j *ServiceJWT) getPublicKey() (*rsa.PublicKey, error) {
	keyPath := j.cfg.PublicKeyPath
	readKey, err := os.ReadFile(keyPath)
	if err != nil{
		return nil, errors.New("failed to read public key")
	}
	block, _ := pem.Decode(readKey)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		j.log.Error("failed to parse public key", err.Error())
		return nil, err
	}
	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("not an RSA public key")
	}
	return publicKey, nil
}

func (j *ServiceJWT) getPrivateKey() (*rsa.PrivateKey, error) {
	keyPath := j.cfg.PrivateKeyPath
	readKey, err := os.ReadFile(keyPath)
	block, _ := pem.Decode(readKey)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the privateKey")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		j.log.Error("failed to parse privateKey", err)
		return nil, err
	}
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		log.Fatalf("not an RSA private key")
	}
	return privateKey, nil
}

func (j *ServiceJWT) generateToken(ctx context.Context, data *DataToEncode) (*sso.Token, error) {
	var (
		token *jwt.Token
	)
	if j == nil {
		return nil, errors.New("jwt service is nil")
	}
	privateKey, err := j.getPrivateKey()

	if err != nil {
		return nil, err
	}
	mappedData, err := utils.StructToMap(data)
	if err != nil {
		return nil, err
	}
	dataToEncode := jwt.MapClaims{}
	for k, v := range mappedData {
		dataToEncode[k] = v
	}
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, dataToEncode)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	ssoToken := &sso.Token{Token: signedToken}
	return ssoToken, nil
}

func (j *ServiceJWT) generateAccessToken(ctx context.Context, subjectEmail string, deviceId string,
	subjectId uuid.UUID, isMentor bool) (*sso.AccessToken, error) {
	data := DataToEncode{
		TokenType: "access",
		Exp:       int(time.Now().Unix()) + (j.cfg.tokens * 60),
		Iat:       int(time.Now().Unix()),
		Jti:       uuid.New().String(),
		DeviceId:  deviceId,
		Issuer:    "auth.service",
		Subject:   subjectEmail,
		SubjectId: subjectId.String(),
		IsMentor:  isMentor,
	}
	token, err := j.generateToken(ctx, &data)
	if err != nil {
		return nil, err
	}
	jti, err := uuid.Parse(data.Jti)
	if err != nil {
		return nil, err
	}
	dbToken := repository.NewToken(jti, token.Token, "access", subjectEmail, deviceId)
	err = j.tokenRepository.SaveTokens(ctx, *dbToken)
	if err != nil {
		return nil, err
	}
	accessToken := sso.AccessToken{AccessToken: token}
	return &accessToken, nil
}

func (j *ServiceJWT) generateRefreshToken(ctx context.Context, subjectEmail string, deviceId string,
	subjectId uuid.UUID, isMentor bool) (*sso.RefreshToken, error) {
	data := DataToEncode{
		TokenType: "refresh",
		Exp:       int(time.Now().AddDate(time.Now().Year(), j.cfg.RefreshTokenDurationDays, time.Now().Day()).Unix()),
		Iat:       int(time.Now().Unix()),
		Jti:       uuid.New().String(),
		DeviceId:  deviceId,
		Issuer:    "auth.service",
		Subject:   subjectEmail,
		SubjectId: subjectId.String(),
		IsMentor:  isMentor,
	}
	token, err := j.generateToken(ctx, &data)
	if err != nil {
		return nil, err
	}
	jti, err := uuid.Parse(data.Jti)
	if err != nil {
		return nil, err
	}
	dbToken := repository.NewToken(jti, token.Token, "refresh", subjectEmail, deviceId)
	err = j.tokenRepository.SaveTokens(ctx, *dbToken)
	if err != nil {
		return nil, err
	}

	refreshToken := sso.RefreshToken{RefreshToken: token}
	return &refreshToken, nil
}

func (j *ServiceJWT) TokensIssue(ctx context.Context, subjectEmail string, subjectId uuid.UUID,
	deviceId string, isMentor bool) (*sso.PairTokens, error) {
	accessToken, errAccess := j.generateAccessToken(ctx, subjectEmail, deviceId, subjectId, isMentor)
	if errAccess != nil {
		return nil, errAccess
	}
	refreshToken, errRefresh := j.generateRefreshToken(ctx, subjectEmail, deviceId, subjectId, isMentor)
	if errRefresh != nil {
		return nil, errRefresh
	}
	tokenPair := sso.PairTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return &tokenPair, nil
}

func (j *ServiceJWT) VerifyToken(ctx context.Context, token string) (jwt.MapClaims, error) {
	j.log.Info("ServiceJWT.VerifyToken")
	publicKey, err := j.getPublicKey()
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, errors.New("invalid token expiration time")
		}
		j.log.Info("claims", claims)
		if int64(exp) < time.Now().Unix() {
			return nil, errors.New("token is expired")
		}
		jti, err := uuid.Parse(claims["jti"].(string))
		if err != nil {
			return nil, err
		}
		_, err = j.tokenRepository.GetToken(ctx, jti)
		if err != nil {
			return nil, fmt.Errorf("token isn`t correct")
		}

		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (j *ServiceJWT) TokenRevoke(ctx context.Context, subjectEmail string, deviceId string, token *sso.AccessToken) (*sso.VerifyTokenResponse, error) {
	j.log.Info("ServiceJWT.TokenRevoke")
	err := j.tokenRepository.RevokeTokens(ctx, subjectEmail, deviceId)
	if err != nil {
		return nil, err
	}
	res := j.redisClient.Set(ctx, token.AccessToken.Token, subjectEmail, time.Hour*24*time.Duration(j.cfg.RefreshTokenDurationDays))
	if res.Err() != nil {
		return nil, err
	}
	return &sso.VerifyTokenResponse{
		Status: true,
	}, nil
}

func (j *ServiceJWT) RevokeAllTokens(ctx context.Context, subjectEmail string) error {
	j.log.Info("ServiceJWT.RevokeAllTokens")
	tokens, err := j.tokenRepository.RevokeAllTokens(ctx, subjectEmail)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		res := j.redisClient.Set(ctx, token.Token, subjectEmail, time.Hour*24*time.Duration(j.cfg.RefreshTokenDurationDays))
		if res.Err() != nil {
			return err
		}
	}
	return nil
}

func (j *ServiceJWT) RotateToken(ctx context.Context, token *sso.RefreshToken) (*sso.PairTokens, error) {
	j.log.Info("ServiceJWT.RotateToken")
	refreshTokenClaims, err := j.VerifyToken(ctx, token.RefreshToken.GetToken())
	if err != nil {
		return nil, err
	}
	j.log.Info("refreshTokenClaims", refreshTokenClaims)
	jti := refreshTokenClaims["jti"].(string)
	parsedJti, err := uuid.Parse(jti)
	if err != nil {
		return nil, err
	}
	refreshTokenDb, err := j.tokenRepository.GetToken(ctx, parsedJti)
	if err != nil {
		return nil, err
	}
	subjectId, err := uuid.Parse(refreshTokenClaims["subjectId"].(string))
	isMentor := refreshTokenClaims["isMentor"].(bool)
	if err != nil {
		return nil, err
	}
	accessToken, err := j.generateAccessToken(ctx, refreshTokenDb.UserEmail, refreshTokenDb.DeviceId,
		subjectId, isMentor)
	if err != nil {
		return nil, err
	}
	return &sso.PairTokens{AccessToken: accessToken, RefreshToken: token}, nil
}
