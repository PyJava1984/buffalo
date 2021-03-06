package buffalo

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler is the basis for all of Buffalo. A Handler
// will be given a Context interface that represents the
// give request/response. It is the responsibility of the
// Handler to handle the request/response correctly. This
// could mean rendering a template, JSON, etc... or it could
// mean returning an error.
/*
	func (c Context) error {
		return c.Render(200, render.String("Hello World!"))
	}

	func (c Context) error {
		return c.Redirect(301, "http://github.com/gobuffalo/buffalo")
	}

	func (c Context) error {
		return c.Error(422, errors.New("oops!!"))
	}
*/
type Handler func(Context) error

func (a *App) handlerToHandler(info RouteInfo, h Handler) http.Handler {
	hf := func(res http.ResponseWriter, req *http.Request) {
		ws := res.(*buffaloResponse)
		params := req.URL.Query()
		vars := mux.Vars(req)
		for k, v := range vars {
			params.Set(k, v)
		}

		c := &DefaultContext{
			response: ws,
			request:  req,
			params:   params,
			logger:   a.Logger,
			session:  a.getSession(req, ws),
			notFound: a.notFound(),
			data: map[string]interface{}{
				"routes":        a.Routes(),
				"current_route": info,
			},
		}

		err := a.Middleware.handler(h)(c)

		if err != nil {
			err := c.Error(500, err)
			a.Logger.Error(err)
		}
	}
	return http.HandlerFunc(hf)
}
