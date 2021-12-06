package fnruntime

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"strings"
	"time"
)

type httpContainerFn struct {
	ctx context.Context
	url string
}

func (c *httpContainerFn) Run(reader io.Reader, writer io.Writer) error {
	ctx := c.ctx

	if !strings.HasPrefix(c.url, "http") {
		return errors.New("not an http image")
	}

	startCall := time.Now()

	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fmt.Printf("%v: GetConn: %+v\n", time.Since(startCall), hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("%v: GotConn: %+v\n", time.Since(startCall), connInfo)
		},
		PutIdleConn: func(err error) {
			fmt.Printf("%v: PutIdleConn: %+v\n", time.Since(startCall), err)
		},
		GotFirstResponseByte: func() {
			fmt.Printf("%v: GotFirstResponseByte\n", time.Since(startCall))
		},
		Got100Continue: func() {
			fmt.Printf("%v: Got100Continue\n", time.Since(startCall))
		},
		Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
			fmt.Printf("%v: Got1xxResponse: %+v %+v\n", time.Since(startCall), code, header)
			return nil
		},
		DNSStart: func(x httptrace.DNSStartInfo) {
			fmt.Printf("%v: DNSStart: %+v\n", time.Since(startCall), x)
		},
		DNSDone: func(x httptrace.DNSDoneInfo) {
			fmt.Printf("%v: DNSDone: %+v\n", time.Since(startCall), x)
		},
		ConnectStart: func(network string, addr string) {
			fmt.Printf("%v: ConnectStart: %+v %+v\n", time.Since(startCall), network, addr)
		},
		ConnectDone: func(network string, addr string, err error) {
			fmt.Printf("%v: ConnectDone: %+v %+v %+v\n", time.Since(startCall), network, addr, err)
		},
		TLSHandshakeStart: func() {
			fmt.Printf("%v: TLSHandshakeStart\n", time.Since(startCall))
		},
		TLSHandshakeDone: func(x tls.ConnectionState, err error) {
			fmt.Printf("%v: TLSHandshakeDone: %+v %+v\n", time.Since(startCall), x, err)
		},
		WroteHeaderField: func(key string, value []string) {
			fmt.Printf("%v: WroteHeaderField: %+v %+v\n", time.Since(startCall), key, value)
		},
		WroteHeaders: func() {
			fmt.Printf("%v: WroteHeaders\n", time.Since(startCall))
		},
		Wait100Continue: func() {
			fmt.Printf("%v: Wait100Continue", time.Since(startCall))
		},
		WroteRequest: func(x httptrace.WroteRequestInfo) {
			fmt.Printf("%v: WroteRequest: %+v\n", time.Since(startCall), x)
		},
	}

	req, err := http.NewRequest("POST", c.url, reader)
	if err != nil {
		return fmt.Errorf("sending request %s: %w", c.url, err)
	}
	resp, err := http.DefaultTransport.RoundTrip(req.WithContext(httptrace.WithClientTrace(ctx, trace)))
	if err != nil {
		return fmt.Errorf("sending request %s: %w", c.url, err)
	}

	// TODO(dejardin) error for non-200

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return fmt.Errorf("receiving response %s: %w", c.url, err)
	}

	return nil
}

type HttpFnRunner interface {
	Run(reader io.Reader, writer io.Writer) error
}

func HttpFnRunnerForContainer(cfn *ContainerFn) HttpFnRunner {
	crcfn := &httpContainerFn{
		ctx: cfn.Ctx,
		url: cfn.Image,
	}
	return crcfn
}

func HttpFnRunnerForExec(efn *ExecFn) HttpFnRunner {
	crcfn := &httpContainerFn{
		ctx: context.Background(),
		url: efn.Path,
	}
	return crcfn
}
