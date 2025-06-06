// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Home(name string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"flex flex-col items-center justify-center h-screen\" x-data=\"{ lang: &#39;python3&#39;, theme: &#39;vs-dark&#39; }\"><h1 class=\"text-4xl font-bold\">CodeEngine</h1><p class=\"text-gray-500\">CodeEngine is a platform for executing code in a sandboxed environment.</p><!-- Editor Controls --><div class=\"flex gap-4 mb-4\"><!-- Language Selector --><div class=\"flex flex-col\"><label class=\"text-sm font-medium text-gray-700 mb-1\">Language:</label> <select x-model=\"lang\" @change=\"$store.editorState.setLanguage(lang)\" class=\"p-2 border border-gray-300 rounded-md\"><option value=\"python3\">Python</option> <option value=\"nodejs\">JavaScript</option></select></div><!-- Theme Selector --><div class=\"flex flex-col\"><label class=\"text-sm font-medium text-gray-700 mb-1\">Theme:</label> <select x-model=\"theme\" @change=\"$store.editorState.setTheme(theme)\" class=\"p-2 border border-gray-300 rounded-md\"><option value=\"vs\">Light</option> <option value=\"vs-dark\">Dark</option> <option value=\"hc-black\">High Contrast</option></select></div></div><form action=\"/execute\" method=\"post\" class=\"flex flex-col w-full max-w-4xl items-center justify-center\"><div id=\"container\" style=\"min-height: 400px; width: 100%;\" class=\"tailwind-ignore border border-gray-300 rounded-md mb-4\"></div><input type=\"hidden\" name=\"code\" id=\"code\"> <input type=\"hidden\" name=\"language\" x-bind:value=\"lang\"> <button type=\"submit\" class=\"bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded-md transition-colors\">Execute Code</button></form><!-- Status Display --><div class=\"mt-4 text-sm text-gray-600\"><span>Language: <span x-text=\"lang\" class=\"font-medium\"></span></span> <span class=\"mx-2\">|</span> <span>Theme: <span x-text=\"theme\" class=\"font-medium\"></span></span></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
