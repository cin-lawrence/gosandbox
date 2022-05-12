package services

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cin-lawrence/gosandbox/pkg/config"
	"github.com/cin-lawrence/gosandbox/pkg/db"
	"github.com/cin-lawrence/gosandbox/pkg/models"

	"golang.org/x/crypto/bcrypt"
	jwt "github.com/golang-jwt/jwt/v4"
        redis "github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
)


type AuthService struct{
        Redis *redis.Client
}

func NewAuthService() *AuthService {
        srv := &AuthService{
                Redis: db.RedisClient,
        }
        return srv
}

func (srv *AuthService) GetSecretIfValid(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.Config.AccessTokenSecret), nil
}

func (srv *AuthService) CreateToken(userID uint) (*models.TokenMeta, error) {
        tm := &models.TokenMeta{}
        tm.AccessExpires = time.Now().Add(time.Second * time.Duration(config.Config.AccessTokenExpirationTime)).Unix()
        tm.AccessUUID = uuid.Must(uuid.NewV4())
        tm.RefreshExpires = time.Now().Add(time.Second * time.Duration(config.Config.RefreshTokenExpirationTime)).Unix()
        tm.RefreshUUID = uuid.Must(uuid.NewV4())

        var err error
        accessTokenClaims := jwt.MapClaims{}
        accessTokenClaims["authorized"] = true
        accessTokenClaims["access_uuid"] = tm.AccessUUID
        accessTokenClaims["user_id"] = userID
        accessTokenClaims["exp"] = tm.AccessExpires

        accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
        tm.AccessToken, err = accessToken.SignedString([]byte(config.Config.AccessTokenSecret))
        if err != nil {
                return nil, err
        }

        refreshTokenClaims := jwt.MapClaims{}
        refreshTokenClaims["refresh_uuid"] = tm.RefreshUUID
        refreshTokenClaims["user_id"] = userID
        refreshTokenClaims["exp"] = tm.RefreshExpires
        refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
        tm.RefreshToken, err = refreshToken.SignedString([]byte(config.Config.RefreshTokenSecret))
        if err != nil {
                return nil, err
        }
        return tm, nil
}

func (srv *AuthService) ExtractToken(r *http.Request) string {
        bearerToken := r.Header.Get("Authorization")
        compartments := strings.Split(bearerToken, " ")
        if len(compartments) != 2 {
                return ""
        }
        return compartments[1]
}

func (srv *AuthService) VerifyToken(r *http.Request) (*jwt.Token, error) {
        tokenString := srv.ExtractToken(r)
        token, err := jwt.Parse(tokenString, srv.GetSecretIfValid)
        if err != nil {
                return nil, err
        }
        return token, nil
}

func (srv *AuthService) ExtractAccessMeta(r *http.Request) (*models.AccessMeta, error) {
        token, err := srv.VerifyToken(r)
        if err != nil {
                return nil, err
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !(ok && token.Valid) {
                return nil, err
        }
        accessUUID, ok := claims["access_uuid"].(string)
        if !ok {
                return nil, err
        }
        userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
        if err != nil {
                return nil, err
        }
        return &models.AccessMeta{
                AccessUUID: uuid.Must(uuid.FromString(accessUUID)),
                UserID: userID,
        }, nil
}

func (srv *AuthService) CreateAuth(userID uint, tm *models.TokenMeta) error {
        accessTokenExpires := time.Unix(tm.AccessExpires, 0)
        refreshTokenExpires := time.Unix(tm.RefreshExpires, 0)
        now := time.Now()

        err := srv.Redis.Set(
                db.RedisContext,
                tm.AccessUUID.String(),
                strconv.Itoa(int(userID)),
                accessTokenExpires.Sub(now),
        ).Err()
        if err != nil {
                return err
        }

        err = srv.Redis.Set(
                db.RedisContext,
                tm.RefreshUUID.String(),
                strconv.Itoa(int(userID)),
                refreshTokenExpires.Sub(now),
        ).Err()
        if err != nil {
                return err
        }
        return nil
}

func (srv *AuthService) FetchAuth(am *models.AccessMeta) (int64, error) {
        userID, err := srv.Redis.Get(db.RedisContext, am.AccessUUID.String()).Result()
        if err != nil {
                return 0, err
        }

        userIDString, _ := strconv.ParseInt(userID, 10, 64)
        return userIDString, nil
}

func (srv *AuthService) DeleteAuth(uuid string) (int64, error) {
        deleted, err := srv.Redis.Del(db.RedisContext, uuid).Result()
        if err != nil {
                return 0, err
        }
        return deleted, nil
}

func (srv *AuthService) Login(userLogin models.UserLogin, user models.User) (tokens models.Tokens, err error) {
        bytePassword := []byte(userLogin.Password)
        byteHashedPassword := []byte(user.HashedPassword)

        err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
        if err != nil {
                return tokens, err
        }

        tokenMeta, err := srv.CreateToken(user.ID)
        if err != nil {
                return tokens, err
        }

        err = srv.CreateAuth(user.ID, tokenMeta)
        if err != nil {
                return tokens, err
        }

        tokens.AccessToken = tokenMeta.AccessToken
        tokens.RefreshToken = tokenMeta.RefreshToken
        return tokens, nil
}
