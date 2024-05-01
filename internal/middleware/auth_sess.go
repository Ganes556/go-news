package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthMidleware interface {
	Authorized(*fiber.Ctx) error
}

type authMidleware struct {
	session *session.Store
}

func NewAuthMiddleware(session *session.Store) AuthMidleware {
	return &authMidleware{session}
}

func (a *authMidleware) Authorized(c *fiber.Ctx) error {
	sess, _ := a.session.Get(c)

	if c.Path() != "/news" {
		if sess.Get("id") == nil && c.Path() != "/user/login" {
			return c.Status(fiber.StatusSeeOther).Redirect("/user/login")
		}
		if c.Path() == "/user/login" && sess.Get("id") != nil {
			return c.Status(fiber.StatusSeeOther).Redirect("/user")
		}
	}

	if sess.Get("id") != nil {
		c.Set("user_name", sess.Get("name").(string))
		c.Set("user_username", sess.Get("username").(string))
	}

	return c.Next()
}
