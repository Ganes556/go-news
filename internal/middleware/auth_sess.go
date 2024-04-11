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
	sess, err := a.session.Get(c)
	
	if err != nil {
		return c.Status(fiber.StatusSeeOther).Redirect("/user/login")
	}
	if sess.Get("username") == nil {
		return c.Status(fiber.StatusSeeOther).Redirect("/user/login")
	}
	if c.Path() == "/user/login" {
		return c.Status(fiber.StatusSeeOther).Redirect("/user/dashboard")
	}
	return c.Next()
}
