package fasthttp

import (
	"bytes"

	"github.com/fasthttp/websocket"
	"github.com/qbeon/webwire-go"
	"github.com/valyala/fasthttp"
)

var methodNameOptions = []byte("OPTIONS")

func (srv *Transport) handleAccept(ctx *fasthttp.RequestCtx) {
	// Reject incoming connections during shutdown, pretend the server is
	// temporarily unavailable
	if srv.isShuttingdown() {
		ctx.Response.Header.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return
	}

	// Handle OPTION requests
	method := ctx.Method()
	if bytes.Equal(method, methodNameOptions) {
		if srv.OnOptions != nil {
			srv.OnOptions(ctx)
		}
		return
	}

	connectionOptions := webwire.ConnectionOptions{
		Connection:       webwire.Accept,
		ConcurrencyLimit: 0,
	}
	if srv.BeforeUpgrade != nil {
		connectionOptions = srv.BeforeUpgrade(ctx)
	}

	// Abort connection establishment if the connection was refused
	if connectionOptions.Connection != webwire.Accept {
		return
	}

	// Copy the user agent string
	ua := ctx.UserAgent()
	userAgent := make([]byte, len(ua))
	copy(userAgent, ua)
	connectionOptions.Info[0] = userAgent

	if err := srv.Upgrader.Upgrade(
		ctx,
		func(conn *websocket.Conn) {
			srv.handleConnection(connectionOptions, conn)
		},
	); err != nil {
		// Establish connection
		srv.ErrorLog.Print("upgrade failed:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}
