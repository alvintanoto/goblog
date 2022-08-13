package route

import (
	"net/http"
	"time"

	"alvintanoto.id/blog/internal/database"
	"alvintanoto.id/blog/internal/database/connection"
	model "alvintanoto.id/blog/internal/model/response"
	t "alvintanoto.id/blog/internal/template"
	"alvintanoto.id/blog/pkg/forms"
	"alvintanoto.id/blog/pkg/helper"
	"alvintanoto.id/blog/pkg/log"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Logger *log.Log
}

var flash string = "Server error, please try again later"

func (h *Handler) Healthz(c echo.Context) error {
	resp := &model.HealthzResponse{
		BaseResponse: model.BaseResponse{
			Code:    http.StatusOK,
			Message: "Success",
		},
		Time: time.Now(),
	}

	b, err := helper.PrettyJson(resp)
	if err != nil {
		h.Logger.ErrorLog.Println(err)
	}

	h.Logger.InfoLog.Printf("Response JSON: \n%s", string(b))
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Home(c echo.Context) error {
	var userID *int = nil
	sess, _ := session.Get("session", c)
	if sess.Values["userID"] != nil {
		id := sess.Values["userID"].(int)
		userID = &id
	}

	posts, err := new(database.PostDB).GetHomePosts(userID)
	if err != nil {
		return c.Render(http.StatusOK, "home.page.html", &t.TemplateData{
			FlashError: flash,
		})
	}

	return c.Render(http.StatusOK, "home.page.html", &t.TemplateData{
		Posts: posts,
	})
}

func (h *Handler) CreatePostForm(c echo.Context) error {
	return c.Render(http.StatusOK, "create_post.page.html", &t.TemplateData{
		Form: forms.New(nil),
	})
}

func (h *Handler) CreatePost(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		return c.Render(http.StatusOK, "create_post.page.html", &t.TemplateData{
			Form:       forms.New(nil),
			FlashError: flash,
		})
	}

	title := c.Request().PostForm.Get("title")
	content := c.Request().PostForm.Get("content")
	isPublicValue := c.Request().PostForm.Get("is_public")
	isPublic := false
	if len(isPublicValue) > 0 {
		isPublic = true
	}

	form := forms.New(c.Request().PostForm)
	form.Required("title", "content")
	form.MaxLength("title", 30)
	if !form.Valid() {
		return c.Render(http.StatusOK, "create_post.page.html", &t.TemplateData{
			Form: form,
		})
	}

	sess, _ := session.Get("session", c)
	userID := sess.Values["userID"].(int)

	_, err = new(database.PostDB).Insert(title, content, isPublic, userID)
	if err != nil {
		return c.Render(http.StatusOK, "signup.page.html", &t.TemplateData{
			Form:       form,
			FlashError: flash,
		})
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) SignupForm(c echo.Context) error {
	return c.Render(http.StatusOK, "signup.page.html", &t.TemplateData{
		Form: forms.New(nil),
	})
}

func (h *Handler) Signup(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		return c.Render(http.StatusOK, "signup.page.html", &t.TemplateData{
			Form:       forms.New(nil),
			FlashError: flash,
		})
	}

	username := c.Request().PostForm.Get("username")
	password := c.Request().PostForm.Get("password")

	form := forms.New(c.Request().PostForm)
	form.Required("username", "password", "confirm_password")
	form.MinLength("password", 8)
	form.MinLength("confirm_password", 8)
	form.MaxLength("username", 20)
	form.Match("password", "confirm_password")
	if !form.Valid() {
		return c.Render(http.StatusOK, "signup.page.html", &t.TemplateData{
			Form: form,
		})
	}

	_, err = new(database.UserDB).Insert(username, password)
	if err != nil {
		if err == connection.ErrConflictData {
			flash = "User already exist"
			return c.Render(http.StatusOK, "signup.page.html", &t.TemplateData{
				Form:       form,
				FlashError: flash,
			})
		}

		return c.Render(http.StatusOK, "signup.page.html", &t.TemplateData{
			Form:       form,
			FlashError: flash,
		})
	}

	sess, _ := session.Get("session", c)
	sess.Values["flash"] = "Your signup was successful!"
	sess.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusSeeOther, "/user/login")
}

func (h *Handler) LoginForm(c echo.Context) error {
	sess, _ := session.Get("session", c)
	var flash string

	if sess.Values["flash"] != nil {
		flash = sess.Values["flash"].(string)
		delete(sess.Values, "flash")
	}

	sess.Save(c.Request(), c.Response().Writer)

	return c.Render(http.StatusOK, "login.page.html", &t.TemplateData{
		Form:  forms.New(nil),
		Flash: flash,
	})
}

func (h *Handler) Login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		c.JSON(echo.ErrBadRequest.Code, map[string]interface{}{"error": "bad request"})
	}

	form := forms.New(c.Request().PostForm)
	form.Required("username", "password")
	if !form.Valid() {
		if err == connection.ErrConflictData {
			flash = "User already exist"
			return c.Render(http.StatusOK, "login.page.html", &t.TemplateData{
				Form:       form,
				FlashError: flash,
			})
		}

		return c.Render(http.StatusOK, "login.page.html", &t.TemplateData{
			Form:       form,
			FlashError: flash,
		})
	}

	id, err := new(database.UserDB).Authenticate(form.Get("username"), form.Get("password"))
	if err != nil {
		if err == connection.ErrInvalidCredential {
			flash = "Username or password is incorrect"
			return c.Render(http.StatusOK, "login.page.html", &t.TemplateData{
				Form:       form,
				FlashError: flash,
			})
		}
	}

	sess, _ := session.Get("session", c)
	sess.Values["userID"] = id
	sess.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	delete(sess.Values, "userID")

	sess.Save(c.Request(), c.Response().Writer)

	return c.Redirect(http.StatusSeeOther, "/")
}
