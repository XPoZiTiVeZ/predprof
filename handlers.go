package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "golang.org/x/crypto/bcrypt"
)


type LoginFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterFormData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	RPassword string `json:"rpassword"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewUser(Email, Password string, IsAuthenticated, IsActive, IsAdmin, IsSuperuser bool) User {
	return User{
		Email:           Email,
		Password:		 Password,
		IsAuthenticated: IsAuthenticated,
		IsActive:        IsActive,
		IsAdmin:         IsAdmin,
		IsSuperuser:     IsSuperuser,
	}
}

func Auth(w http.ResponseWriter, r *http.Request) User {
	c, err := r.Cookie("token")
	if err != nil { return NewUser("Anonymous", "", false, false, false, false) }

	tokenStr := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })

	if err != nil || !token.Valid { return NewUser("Anonymous", "", false, false, false, false) }

	user, err := GetUserByEmail(claims.Email)
	if err != nil { return NewUser("Anonymous", "", false, false, false, false) }

	return user
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	user := Auth(w, r)
	if !user.IsAuthenticated { http.Redirect(w, r, "/login", http.StatusSeeOther) }

	switch r.Method {
	case "GET":
		tmp := template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/home.html",
		))

		ctx := struct {
			User  User
		}{
			User:  user,
		}

		tmp.Execute(w, ctx)
	default:
		http.Error(w, "Метод запрещён", http.StatusMethodNotAllowed)
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	user := Auth(w, r)
	if user.IsAuthenticated { http.Redirect(w, r, "/", http.StatusSeeOther) }

	switch r.Method {
	case "GET":
		tmp := template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/register.html",
		))

		ctx := struct {
			User  User
		}{
			User:  user,
		}

		tmp.Execute(w, ctx)
	case "POST":
		var formData RegisterFormData
		err := json.NewDecoder(r.Body).Decode(&formData)
		if err != nil {
			http.Error(w, "Неправильный данные запроса", http.StatusBadRequest)
			return
		}

		if formData.Email == "" {
			http.Error(w, "Нет почты", http.StatusBadRequest)
			return
		} else if formData.Password == "" {
			http.Error(w, "Нет пароля", http.StatusBadRequest)
			return
		} else if formData.RPassword == "" {
			http.Error(w, "Нет повтора пароля", http.StatusBadRequest)
			return
		}

		_, err = RegisterUser(formData)
		if (err == UserExists{}) {
			http.Error(w, "Почта уже зарегистрирована", http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, "Ошибка", http.StatusBadGateway)
			return
		}

		http.Redirect(w, r, "/", http.StatusOK)
	default:
		http.Error(w, "Метод запрещён", http.StatusMethodNotAllowed)
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	user := Auth(w, r)
	if user.IsAuthenticated { http.Redirect(w, r, "/", http.StatusSeeOther) }

	switch r.Method {
	case "GET":
		tmp := template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/login.html",
		))

		ctx := struct {
			User  User
		}{
			User:  user,
		}

		tmp.Execute(w, ctx)
	case "POST":
		var formData LoginFormData
		err := json.NewDecoder(r.Body).Decode(&formData)
		if err != nil {
			http.Error(w, "Неправильный данные запроса", http.StatusBadRequest)
			return
		}

		if formData.Email == "" {
			http.Error(w, "Нет почты", http.StatusBadRequest)
			return
		} else if formData.Password == "" {
			http.Error(w, "Нет пароля", http.StatusBadRequest)
			return
		}

		_, err = AuthenticateUser(formData)
		if (err == UserNotExists{}) {
			http.Error(w, "Пользователь не найден", http.StatusBadRequest)
			return
		} else if (err == IncorrectPassword{}) {
			http.Error(w, "Пароль неправильный", http.StatusBadRequest)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Email: formData.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Неуспешное создание токена", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		http.Redirect(w, r, "/", http.StatusOK)
	default:
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	user := Auth(w, r)
	if !user.IsAuthenticated { http.Redirect(w, r, "/login", http.StatusSeeOther) }

	switch r.Method {
	case "GET":
		tmp := template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/profile.html",
		))

		ctx := struct {
			User  User
		}{
			User:  user,
		}

		tmp.Execute(w, ctx)
	}
}