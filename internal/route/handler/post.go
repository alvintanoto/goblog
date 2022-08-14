package route

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"alvintanoto.id/blog/internal/database"
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

func (h *Handler) ReadPost(c echo.Context) error {
	id := c.Param("id")

	if len(id) == 0 {
		return c.Render(http.StatusNotFound, "not_found.page.html", &t.TemplateData{})
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		return c.Render(http.StatusNotFound, "not_found.page.html", &t.TemplateData{})
	}

	post, err := new(database.PostDB).Get(i)
	if err != nil {
		return c.Render(http.StatusNotFound, "not_found.page.html", &t.TemplateData{})
	}

	if post == nil {
		return c.Render(http.StatusNotFound, "not_found.page.html", &t.TemplateData{})
	}

	sess, _ := session.Get("session", c)
	userID := sess.Values["userID"].(int)

	if !post.IsPublic {
		if post.CreatedBy != userID {
			return c.Render(http.StatusNotFound, "not_found.page.html", &t.TemplateData{})
		}
	}

	fmt.Println(userID == post.CreatedBy)
	return c.Render(http.StatusOK, "read_post.page.html", &t.TemplateData{
		Post:        post,
		IsPostOwner: userID == post.CreatedBy,
	})
}
