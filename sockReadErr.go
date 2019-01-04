package fasthttp

import (
	"fmt"

	"github.com/fasthttp/websocket"
)

// ErrSockRead implements the webwire.ErrSockRead interface
type ErrSockRead struct {
	cause error
}

// Error implements the Go error interface
func (err ErrSockRead) Error() string {
	return fmt.Sprintf("reading socket failed: %s", err.cause)
}

// IsCloseErr implements the ErrSockRead interface
func (err ErrSockRead) IsCloseErr() bool {
	return websocket.IsCloseError(
		err.cause,
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}

// SockReadWrongMsgTypeErr implements the ErrSockRead interface
type SockReadWrongMsgTypeErr struct {
	messageType int
}

// Error implements the Go error interface
func (err SockReadWrongMsgTypeErr) Error() string {
	return fmt.Sprintf("invalid websocket message type: %d", err.messageType)
}

// IsCloseErr implements the ErrSockRead interface
func (err SockReadWrongMsgTypeErr) IsCloseErr() bool {
	return false
}
