package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
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

func (h *Handler) userIdentify(c *gin.Context) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	tokenAuth, err := ExtractTokenMetadata(tokenString)
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

func (h *Handler) userLogout(c *gin.Context) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	au, err := ExtractTokenMetadata(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	deleted, delErr := h.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
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

func verifyToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]))
		}
		return []byte(os.Getenv("ADMIN_SIGNING_KEY")), nil
	})

	if err != nil {
		// fmt.Println("error in verifyToken")
		return nil, err
	}

	return token, nil
}

func ExtractTokenMetadata(tokenString string) (*AccessDetails, error) {
	// fmt.Println(tokenString)
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	// fmt.Println("verify token success")
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.Atoi(fmt.Sprintf("%.f", claims["user_id"]))
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}

	// fmt.Println("token is no valid")
	return nil, err
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

func (h *Handler) CreateToken(adminId int) (*TokenDetails, error) {
	var err error
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 3).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	//Creating Access Token
	// os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = adminId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ADMIN_SIGNING_KEY")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	// os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = adminId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("ADMIN_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil

}

func (h *Handler) refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ADMIN_REFRESH_SECRET")), nil
	})
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Refresh token expired")
		return
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		userId, err := strconv.Atoi(fmt.Sprintf("%.f", claims["user_id"]))
		if err != nil {
			newErrorResponse(c, http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := h.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := h.CreateToken(userId)
		if createErr != nil {
			newErrorResponse(c, http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := h.CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
