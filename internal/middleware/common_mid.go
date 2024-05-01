package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type CommonMid interface {
	IsAdmin(c *fiber.Ctx) error
}

type commonMid struct{
	session *session.Store
}

func NewCommonMid(s *session.Store) CommonMid {
	return &commonMid{s}
}

func (cm *commonMid) IsAdmin(c *fiber.Ctx) error {	
	sess, _ := cm.session.Get(c)
	
	if sess.Get("id") != nil {
		return c.Redirect("/user")
	}

	return c.Next()
}
