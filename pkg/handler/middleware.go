package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

type AccessDetails struct {
	AccessUuid string
	UserId     int
}
type token struct {
	jwt.StandardClaims
	AdminId int `json:"admin_id"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func (h *Handler) adminIdentify(c *gin.Context) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	tokenAuth, err := ExtractAdminTokenMetadata(tokenString)
	// fmt.Println("userIdentify, tokenauth", tokenAuth)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := h.FetchAuth(tokenAuth)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "token is wrong")
		// fmt.Println("return from middleware")
		return
	}
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func (h *Handler) userIdentify(c *gin.Context) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	tokenAuth, err := ExtractUserTokenMetadata(tokenString)
	// fmt.Println("userIdentify, tokenauth", tokenAuth)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := h.FetchAuth(tokenAuth)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "token is wrong")
		// fmt.Println("return from middleware")
		return
	}
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func ExtractToken(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
		// newErrorResponse(c, http.StatusUnauthorized, "")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
		// newErrorResponse(c, http.StatusUnauthorized, "")
	}

	if len(headerParts[1]) == 0 {
		// newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

func (h *Handler) DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := h.redis.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func (h *Handler) FetchAuth(authD *AccessDetails) (int, error) {
	// fmt.Println("FetchAuth authD", authD)
	userid, err := h.redis.Get(authD.AccessUuid).Result()
	// fmt.Println("FetchAuth user_id", userid)
	if err != nil {
		// fmt.Println("FetchAuth error", err.Error())
		return 0, err
	}
	userID, _ := strconv.Atoi(userid)
	return userID, nil
}

func (h *Handler) CreateAuth(userid int, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := h.redis.Set(td.AccessUuid, strconv.Itoa(userid), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := h.redis.Set(td.RefreshUuid, strconv.Itoa(userid), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
