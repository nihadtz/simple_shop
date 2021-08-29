package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/context"
	"github.com/nihadtz/simple_shop/models"
	"github.com/nihadtz/simple_shop/services"
)

type Provider struct {
	rules *casbin.Enforcer
}

var (
	pubkey []byte

	ErrUserNotFound = errors.New("This User doesn't exists")
)

func (m Provider) JWT(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	uri := strings.Split(req.RequestURI, ".")
	routes := strings.Split(uri[0], "/")
	fullAction := uri[0]

	action := "/"

	if len(routes) > 1 {
		action += strings.Split(uri[0], "/")[1]
	}

	if action == "/register" || action == "/products" || fullAction == "/dashboard/public" {
		//Publicly open
	} else if action == "/logout" {
		user, err := checkToken(req)

		if err != nil {
			services.LogError("Error checking token", err)
			services.Renderer.Error(res, http.StatusUnauthorized, "Unauthorized")
			return
		}

		context.Set(req, "user", user)

	} else if action == "/login" {
		var auth map[string]interface{}

		err := parseRequestData(req, &auth)

		if err != nil {
			services.LogError("Error parsing request data", err)
			services.Renderer.Error(res, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := issueToken(auth, req)

		if err != nil {

			if err == ErrUserNotFound {
				services.Renderer.Error(res, http.StatusNotFound, err.Error())
				return
			}

			services.Renderer.Error(res, http.StatusUnauthorized, "Email and password do not match")
			return
		}

		context.Set(req, "user", user)
	} else {

		user, err := checkToken(req)

		if err != nil {
			services.LogError("Error checking token", err)
			services.Renderer.Error(res, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if m.rules != nil && m.rules.Enforce(user.Type, fullAction, requestMethod2Mode(req.Method)) {
			services.LogInfo("RBAC allowed by type " + user.Type + " and email " + user.Email + " for action " + fullAction)
		} else {
			services.LogInfo("RBAC not allowed for email " + user.Email + " and type " + user.Type + " for action " + fullAction)
			services.Renderer.Error(res, http.StatusForbidden, "Access forbidden")
			return
		}

		context.Set(req, "user", user)
	}

	next(res, req)
}

func requestMethod2Mode(reqMethod string) string {
	if reqMethod == "POST" || reqMethod == "PUT" || reqMethod == "DELETE" || reqMethod == "PATCH" {
		return "write"
	}

	return "read"
}

func issueToken(auth map[string]interface{}, req *http.Request) (models.User, error) {
	user := models.User{}

	if auth["email"] == nil || auth["password"] == nil {
		err := errors.New("Email or Passwrod not provided")

		return user, err
	}

	err := user.GetByEmailPassword(auth["email"].(string), auth["password"].(string))

	if err != nil {
		return user, err
	}

	//Token has not expired yet
	if time.Now().Before(time.Unix(user.Issued, 0)) {
		return user, nil
	}

	user.Issued = time.Now().Add(time.Hour * 60).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  user.Email,
		"expiry": user.Issued,
	})

	user.Token, err = token.SignedString(pubkey)

	if err != nil {
		return user, err
	}

	user.Algorithm = jwt.SigningMethodHS256.Name
	user.Password = ""

	err = user.UpdateTokenInfo()

	if err != nil {
		return user, err
	}

	return user, nil
}

func checkToken(req *http.Request) (models.User, error) {
	var user models.User

	token, err := jwtreq.ParseFromRequest(req, jwtreq.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Not using correct algorithm")
			}
			return pubkey, nil
		})

	if err != nil || !token.Valid {
		return user, errors.New("Invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	err = user.GetByEmail(claims["email"].(string))

	if err != nil {
		return user, errors.New("User for provided token not found")
	}

	if user.Token != token.Raw {
		return models.User{}, errors.New("Invalid token")
	}

	if user.Issued != 0 && time.Now().After(time.Unix(user.Issued, 0)) {
		return user, errors.New("Token expired")
	}

	if user.Algorithm != token.Header["alg"] {
		return user, errors.New("Wrong algorithm in token")
	}

	return user, nil
}

func parseRequestData(req *http.Request, v interface{}) error {
	data := make([]byte, req.ContentLength)

	_, err := req.Body.Read(data)

	defer req.Body.Close()

	if err != nil && err != io.EOF {
		return err
	}

	err = json.Unmarshal(data, &v)

	return err
}

func (m *Provider) SetRBAC(model, policy string) error {
	path, err := os.Getwd()

	if err != nil {
		return err
	}

	m.rules, err = casbin.NewEnforcerSafe(path+model, path+policy, false)

	return err
}
