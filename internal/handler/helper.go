package helper_handler

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/news/helper"
	dto_response "github.com/news/internal/dto/response"
	"github.com/sujit-baniya/flash"
)

func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}

func ReturnErrFlash(c *fiber.Ctx, redirectTo string, errs []dto_response.Response) error {
	header := c.GetReqHeaders()
	if errs == nil {
		errs = append(errs, dto_response.Response{
			Message: fiber.ErrInternalServerError.Message,
			Code:    fiber.ErrInternalServerError.Code,
		}) 
	}
	
	if redirectTo == "" {
		redirectTo = c.Path()
	}

	if header["Hx-Request"] != nil && header["Hx-Request"][0] == "true" {
		c.Set("HX-Refresh", "true")
		return flash.WithError(c, fiber.Map{
			"error":    true,
			"messages": helper.JSONStringify(errs),
		}).SendString("")
	}
	return flash.WithError(c, fiber.Map{
		"error":    true,
		"messages": helper.JSONStringify(errs),
	}).Redirect(redirectTo)
}

func ReturnOkFlash(c *fiber.Ctx, redirectTo string , msg dto_response.Response) error {
	header := c.GetReqHeaders()

	if redirectTo == "" {
		redirectTo = c.Path()
	}

	if msg.Message == "" {
		msg = dto_response.Response{
			Message: "success",
			Code: fiber.StatusOK,
		}
	}

	if header["Hx-Request"] != nil && header["Hx-Request"][0] == "true" {
		c.Set("HX-Refresh", "true")
		return flash.WithSuccess(c, fiber.Map{
			"success":    true,
			"messages": helper.JSONStringify(msg),
		}).SendString("")
	}
	return flash.WithSuccess(c, fiber.Map{
		"error":    true,
		"messages": helper.JSONStringify(msg),
	}).Redirect(redirectTo)
}
