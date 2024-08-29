// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package view_news

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/news/internal/entity"
import "fmt"

func tabsStore(activeCategory string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_tabsStore_c07f`,
		Function: `function __templ_tabsStore_c07f(activeCategory){document.addEventListener('alpine:init', () => {
		Alpine.store('tabs', {
			active: activeCategory,
			init() {
				Alpine.effect(() => {
					let url = new URL(window.location.href);
					let params = url.searchParams;
					if(params.get('category')) {	
						document.title = this.active
					}
				})
			}
		})
	})
}`,
		Call:       templ.SafeScript(`__templ_tabsStore_c07f`, activeCategory),
		CallInline: templ.SafeScriptInline(`__templ_tabsStore_c07f`, activeCategory),
	}
}

func NewsTemplate(categories []entity.Categories, activeCategory string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = tabsStore(activeCategory).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"navbar navbar-expand-md border-bottom\" style=\"background-color: #F7F7FC;\"><div class=\"container\"><div class=\"d-flex\" style=\"height:100px\"><a href=\"#\" class=\"navbar-brand my-auto text-dark\">Brand</a></div><ul class=\"navbar-nav flex-row\"><li class=\"nav-item me-md-2 me-3\"><a href=\"#\" class=\"nav-link\"><i class=\"bi bi-twitter-x\"></i></a></li><li class=\"nav-item me-md-2 me-3\"><a href=\"#\" class=\"nav-link\"><i class=\"bi bi-whatsapp\"></i></a></li></ul></div></nav><nav id=\"navbar\" class=\"navbar navbar-expand-md justify-content-start sticky-top shadow\" style=\"background-color: #F7F7FC\" role=\"navigation\"><div class=\"container\" x-data=\"{btnNavActive: false}\"><div class=\"ms-auto navbar-toggler hamburger border-0\" :class=\"btnNavActive ? &#39;is-active&#39; : &#39;&#39;\" @click=\"btnNavActive = !btnNavActive\" id=\"hamburger-1\" data-bs-toggle=\"collapse\" data-bs-target=\"#navbar-collapse-x\"><span class=\"line\"></span> <span class=\"line\"></span> <span class=\"line\"></span></div><div class=\"collapse navbar-collapse\" id=\"navbar-collapse-x\"><ul class=\"nav w-100 nav-underline flex-column flex-md-row\" x-data>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, v := range categories {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"nav-item\"><a class=\"nav-link fs-6\" style=\"cursor: pointer;\" :class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("{ 'active fw-normal' : $store.tabs.active === '%s' }", v.Name))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/news/layout.templ`, Line: 64, Col: 92}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-on:click=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("$store.tabs.active = '%s' ", v.Name))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/news/layout.templ`, Line: 65, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs("/news?category=" + v.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/news/layout.templ`, Line: 66, Col: 43}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-replace-url=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/news?category=%s", v.Name))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/news/layout.templ`, Line: 67, Col: 65}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#content-news\" hx-trigger=\"click\" hx-swap=\"innerHTML\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(v.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/news/layout.templ`, Line: 71, Col: 16}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></li>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"nav-item ms-auto h-100 position-relative d-flex align-items-center flex-column flex-md-row\" :style=\"btnNavActive ? &#39;width: 100%&#39; : &#39;&#39;\" x-data=\"{open: false, notFound: false, inputVal: &#39;&#39;, mobile: false}\" x-effect=\"\r\n\t\t\t\t\t\tif (window.matchMedia(&#39;(max-width:768px)&#39;).matches) {\r\n\t\t\t\t\t\t\tmobile = true\r\n\t\t\t\t\t\t}\" @mousedown.outside=\"open = false\"><div class=\"input-group\"><input hx-get=\"/news\" hx-trigger=\"keyup delay:500ms changed\" :hx-target=\"mobile ? &#39;#res-search-mobile&#39; : &#39;#res-search&#39;\" x-init=\"$nextTick(() =&gt; htmx.process($el))\" hx-swap=\"innerHTML\" x-model=\"inputVal\" x-on:htmx:after-request=\"\r\n\t\t\t\t\t\t\t\t\tif(inputVal.trim() === &#39;&#39;) {\r\n\t\t\t\t\t\t\t\t\t\topen = false\r\n\t\t\t\t\t\t\t\t\t}else {\r\n\t\t\t\t\t\t\t\t\t\topen = true\r\n\t\t\t\t\t\t\t\t\t}\r\n\t\t\t\t\t\t\t\t\tif($event.detail.xhr.responseText.trim() === &#39;&#39;) {\r\n\t\t\t\t\t\t\t\t\t\tnotFound = true\r\n\t\t\t\t\t\t\t\t\t}else {\r\n\t\t\t\t\t\t\t\t\t\tnotFound = false\r\n\t\t\t\t\t\t\t\t\t}\r\n\t\t\t\t\t\t\t\t\" name=\"search\" class=\"form-control border-end-0 border rounded\" type=\"search\" placeholder=\"search...\" id=\"example-search-input\"> <span x-show=\"!mobile\" class=\"input-group-append\"><button class=\"btn btn-outline-secondary bg-white border-bottom-0 border rounded\" type=\"button\"><i class=\"bi bi-search\"></i></button></span></div><div x-show=\"open &amp;&amp; !mobile\" class=\"list-group position-absolute bg-body-tertiary p-3 shadow-sm\" :style=\"open &amp;&amp; !mobile \r\n                                ? (notFound \r\n                                    ? &#39;width: 180%; overflow-y: auto; max-height: 20rem; transform: translate(-50%, 105%);&#39; \r\n                                    : &#39;width: 180%; overflow-y: auto; max-height: 20rem; transform: translate(-50%, 59%);&#39;\r\n                                ) \r\n                                : &#39;&#39;\"><div id=\"res-search\"></div><div x-show=\"notFound &amp;&amp; !mobile\" class=\"m-auto h6\">Not Found</div><div id=\"error-search-news\" class=\"text-center text-danger\" hx-swap-oob=\"true\"></div></div><div x-show=\"open &amp;&amp; mobile\" x-transition class=\"mt-3 d-flex flex-column\" :style=\"!notFound ? &#39;overflow-y: auto;max-height: 50vh;&#39; : &#39;max-height: 50vh;&#39;\"><div id=\"res-search-mobile\"></div><div x-show=\"open &amp;&amp; notFound &amp;&amp; mobile\" class=\"m-auto h6\">Not Found</div></div></li></ul></div></div></nav><section class=\"container mt-5\"><div class=\"row\"><div x-data id=\"content-news\" class=\"col-12 col-md-7\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"error-get-content-news\" class=\"text-center text-danger\" hx-swap-oob=\"true\"></div></div><div class=\"col-1\"></div><div class=\"col-12 col-md-4 pt-3\" id=\"content-most-viewed\"><h5 class=\"mb-3 text-center\">Most Viewed</h5><div hx-get=\"/news?most_viewed=1\" hx-trigger=\"load once\" hx-target=\"this\" hx-swap=\"innerHTML ignoreTitle:true\"></div><div id=\"error-most-viewed-news\" class=\"text-center text-danger\" hx-swap-oob=\"true\"></div></div></div></section>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
