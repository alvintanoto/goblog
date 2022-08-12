package route

import (
	"net/http"

	model "alvintanoto.id/blog/internal/model/response"
	h "alvintanoto.id/blog/internal/route/handler"
	m "alvintanoto.id/blog/internal/route/middleware"
	t "alvintanoto.id/blog/internal/template"
	"alvintanoto.id/blog/pkg/helper"
	"alvintanoto.id/blog/pkg/log"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Init(port string) {
	logger := log.Get()

	middleware := &m.Middleware{
		Logger: logger,
	}

	handler := &h.Handler{
		Logger: logger,
	}

	t := &t.Template{
		Templates: t.NewTemplateCache("./ui/html/"),
	}

	echo.NotFoundHandler = func(c echo.Context) error {
		resp := model.BaseResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		}

		b, err := helper.PrettyJson(resp)
		if err != nil {
			logger.ErrorLog.Println(err)
		}

		logger.InfoLog.Printf("Response JSON: \n%s", string(b))
		return c.Render(http.StatusOK, "not_found.page.html", nil)
	}

	e := echo.New()
	e.Renderer = t

	store := sessions.NewCookieStore([]byte("xT11TWwO60c*b3&*j42coY9eSPdzJ77W"))
	store.Options = &sessions.Options{
		MaxAge:   12 * 3600,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	e.Use(session.MiddlewareWithConfig(session.Config{
		Store: store,
	}))

	e.Use(middleware.LogRequest)

	e.GET("/", handler.Home)
	e.GET("/healthz", handler.Healthz)

	// posts
	e.GET("/create-post", handler.CreatePostForm)

	// user
	e.GET("/user/signup", handler.SignupForm)
	e.POST("/user/signup", handler.Signup)
	e.GET("/user/login", handler.LoginForm)
	// e.POST("/user/login", handler.Login)

	e.Static("/static", "./ui/static")

	logger.InfoLog.Fatal(e.Start(port))
}
