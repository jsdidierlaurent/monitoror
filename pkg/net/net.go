//go:generate mockery --name Conn

package net

import (
	"net"
	"time"
)

// Conn is a copy og net/net.go Conn interface (only for generate mock faster)
type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}
