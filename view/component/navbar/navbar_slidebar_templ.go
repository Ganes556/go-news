// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package view_navbar

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func checkUrl(id string, page string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_checkUrl_9777`,
		Function: `function __templ_checkUrl_9777(id, page){const currentUrl = window.location.href;
	// Create a URL object (this automatically parses the URL)
	const url = new URL(currentUrl);

	// Get access to the search parameters
	const params = new URLSearchParams(url.search);
	const pageParam = params.get('page');
	if(pageParam === page) {
		htmx.trigger(id, 'htmx:abort')
	}else {
		document.title = page[0].toUpperCase() + page.slice(1);
	}
}`,
		Call:       templ.SafeScript(`__templ_checkUrl_9777`, id, page),
		CallInline: templ.SafeScriptInline(`__templ_checkUrl_9777`, id, page),
	}
}

func Slidebar(param ParamNavbar) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"d-flex flex-column align-items-center align-items-sm-start px-3 pt-2 text-white h-100\"><a href=\"/\" class=\"d-flex align-items-center pb-3 mb-md-0 me-md-auto text-white text-decoration-none\"><span class=\"fs-5 d-none d-sm-inline\">Menu</span></a><ul class=\"nav nav-pills flex-column mb-sm-auto mb-0 align-items-center align-items-sm-start\" id=\"menu\"><li>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, checkUrl("#page-dashboard", "dashboard"))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"page-dashboard\" hx-trigger=\"click\" hx-on:click=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 templ.ComponentScript = checkUrl("#page-dashboard", "dashboard")
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" style=\"cursor: pointer;\" hx-get=\"?page=dashboard&amp;partial=1\" hx-replace-url=\"?page=dashboard\" hx-swap=\"outerHTML\" hx-target=\"#content\" data-bs-toggle=\"collapse\" class=\"nav-link px-0 align-middle text-white\"><i class=\"fs-4 bi bi-speedometer2\"></i> <span class=\"ms-1 d-none d-sm-inline\">Dashboard</span></div></li><li>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, checkUrl("#page-news", "news"))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"page-news\" style=\"cursor: pointer;\" hx-trigger=\"click\" hx-on:click=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 templ.ComponentScript = checkUrl("#page-news", "news")
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-get=\"?page=news&amp;partial=1\" hx-replace-url=\"?page=news\" hx-swap=\"outerHTML\" hx-target=\"#content\" data-bs-toggle=\"collapse\" class=\"nav-link px-0 align-middle text-white\"><i class=\"fs-sm-5 fs-4 bi bi-newspaper\"></i> <span class=\"ms-1 d-none d-sm-inline\">News</span></div></li></ul><div class=\"dropdown pb-4\" style=\"margin-top: auto;\"><a href=\"#\" class=\"d-flex align-items-center text-white text-decoration-none dropdown-toggle\" data-bs-toggle=\"dropdown\" aria-expanded=\"false\"><i class=\"fs-4 bi bi-person-circle\"></i> <span class=\"d-none d-sm-inline mx-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(param.Name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `view\component\navbar\navbar_slidebar.templ`, Line: 37, Col: 54}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></a><ul class=\"dropdown-menu dropdown-menu-primary text-small shadow\"><li><a class=\"dropdown-item\" href=\"#\">Profile</a></li><li><hr class=\"dropdown-divider\"></li><li><a class=\"dropdown-item\" href=\"#\">Sign out</a></li></ul></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
