package session

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"go.khulnasoft.com/velocity"
)

func Test_Session_Middleware(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New())

	app.Get("/get", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess == nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		value, ok := sess.Get("key").(string)
		if !ok {
			return c.Status(velocity.StatusNotFound).SendString("key not found")
		}
		return c.SendString("value=" + value)
	})

	app.Post("/set", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess == nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		// get a value from the body
		value := c.FormValue("value")
		sess.Set("key", value)
		return c.SendStatus(velocity.StatusOK)
	})

	app.Post("/delete", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess == nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		sess.Delete("key")
		return c.SendStatus(velocity.StatusOK)
	})

	app.Post("/reset", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess == nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		if err := sess.Reset(); err != nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		return c.SendStatus(velocity.StatusOK)
	})

	app.Post("/destroy", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess == nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		if err := sess.Destroy(); err != nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		return c.SendStatus(velocity.StatusOK)
	})

	app.Post("/fresh", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess == nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		// Reset the session to make it fresh
		if err := sess.Reset(); err != nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		if sess.Fresh() {
			return c.SendStatus(velocity.StatusOK)
		}
		return c.SendStatus(velocity.StatusInternalServerError)
	})

	// Test GET, SET, DELETE, RESET, DESTROY by sending requests to the respective routes
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.SetRequestURI("/get")
	h := app.Handler()
	h(ctx)
	require.Equal(t, velocity.StatusNotFound, ctx.Response.StatusCode())
	token := string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.NotEmpty(t, token, "Expected Set-Cookie header to be present")
	tokenParts := strings.SplitN(strings.SplitN(token, ";", 2)[0], "=", 2)
	require.Len(t, tokenParts, 2, "Expected Set-Cookie header to contain a token")
	token = tokenParts[1]
	require.Equal(t, "key not found", string(ctx.Response.Body()))

	// Test POST /set
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.SetRequestURI("/set")
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Set the Content-Type
	ctx.Request.SetBodyString("value=hello")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())

	// Test GET /get to check if the value was set
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.SetRequestURI("/get")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	require.Equal(t, "value=hello", string(ctx.Response.Body()))

	// Test POST /delete to delete the value
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.SetRequestURI("/delete")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())

	// Test GET /get to check if the value was deleted
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.SetRequestURI("/get")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusNotFound, ctx.Response.StatusCode())
	require.Equal(t, "key not found", string(ctx.Response.Body()))

	// Test POST /reset to reset the session
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.SetRequestURI("/reset")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	// verify we have a new session token
	newToken := string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.NotEmpty(t, newToken, "Expected Set-Cookie header to be present")
	newTokenParts := strings.SplitN(strings.SplitN(newToken, ";", 2)[0], "=", 2)
	require.Len(t, newTokenParts, 2, "Expected Set-Cookie header to contain a token")
	newToken = newTokenParts[1]
	require.NotEqual(t, token, newToken)
	token = newToken

	// Test POST /destroy to destroy the session
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.SetRequestURI("/destroy")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())

	// Verify the session cookie has expired
	setCookieHeader := string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.Contains(t, setCookieHeader, "max-age=0")

	// Sleep so that the session expires
	time.Sleep(1 * time.Second)

	// Test GET /get to check if the session was destroyed
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.SetRequestURI("/get")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusNotFound, ctx.Response.StatusCode())
	// check that we have a new session token
	newToken = string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.NotEmpty(t, newToken, "Expected Set-Cookie header to be present")
	parts := strings.Split(newToken, ";")
	require.Greater(t, len(parts), 1)
	valueParts := strings.Split(parts[0], "=")
	require.Greater(t, len(valueParts), 1)
	newToken = valueParts[1]
	require.NotEqual(t, token, newToken)
	token = newToken

	// Test POST /fresh to check if the session is fresh
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.SetRequestURI("/fresh")
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	// check that we have a new session token
	newToken = string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.NotEmpty(t, newToken, "Expected Set-Cookie header to be present")
	newTokenParts = strings.SplitN(strings.SplitN(newToken, ";", 2)[0], "=", 2)
	require.Len(t, newTokenParts, 2, "Expected Set-Cookie header to contain a token")
	newToken = newTokenParts[1]
	require.NotEqual(t, token, newToken)
}

func Test_Session_NewWithStore(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New())

	app.Get("/", func(c velocity.Ctx) error {
		sess := FromContext(c)
		id := sess.ID()
		return c.SendString("value=" + id)
	})
	app.Post("/", func(c velocity.Ctx) error {
		sess := FromContext(c)
		id := sess.ID()
		c.Cookie(&velocity.Cookie{
			Name:  "session_id",
			Value: id,
		})
		return nil
	})

	h := app.Handler()

	// Test GET request without cookie
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	// Get session cookie
	token := string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.NotEmpty(t, token, "Expected Set-Cookie header to be present")
	tokenParts := strings.SplitN(strings.SplitN(token, ";", 2)[0], "=", 2)
	require.Len(t, tokenParts, 2, "Expected Set-Cookie header to contain a token")
	token = tokenParts[1]
	require.Equal(t, "value="+token, string(ctx.Response.Body()))

	// Test GET request with cookie
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	require.Equal(t, "value="+token, string(ctx.Response.Body()))
}

func Test_Session_FromSession(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	sess := FromContext(app.AcquireCtx(&fasthttp.RequestCtx{}))
	require.Nil(t, sess)

	app.Use(New())
}

func Test_Session_WithConfig(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	app.Use(New(Config{
		Next: func(c velocity.Ctx) bool {
			return c.Get("key") == "value"
		},
		IdleTimeout: 1 * time.Second,
		KeyLookup:   "cookie:session_id_test",
		KeyGenerator: func() string {
			return "test"
		},
		source:      "cookie_test",
		sessionName: "session_id_test",
	}))

	app.Get("/", func(c velocity.Ctx) error {
		sess := FromContext(c)
		id := sess.ID()
		return c.SendString("value=" + id)
	})

	app.Get("/isFresh", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess.Fresh() {
			return c.SendStatus(velocity.StatusOK)
		}
		return c.SendStatus(velocity.StatusInternalServerError)
	})

	app.Post("/", func(c velocity.Ctx) error {
		sess := FromContext(c)
		id := sess.ID()
		c.Cookie(&velocity.Cookie{
			Name:  "session_id_test",
			Value: id,
		})
		return nil
	})

	h := app.Handler()

	// Test GET request without cookie
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	// Get session cookie
	token := string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.NotEmpty(t, token, "Expected Set-Cookie header to be present")
	tokenParts := strings.SplitN(strings.SplitN(token, ";", 2)[0], "=", 2)
	require.Len(t, tokenParts, 2, "Expected Set-Cookie header to contain a token")
	token = tokenParts[1]
	require.Equal(t, "value="+token, string(ctx.Response.Body()))

	// Test GET request with cookie
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.Header.SetCookie("session_id_test", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	require.Equal(t, "value="+token, string(ctx.Response.Body()))

	// Test POST request with cookie
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.Header.SetCookie("session_id_test", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())

	// Test POST request without cookie
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())

	// Test POST request with wrong key
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.Header.SetCookie("session_id", token)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())

	// Test POST request with wrong value
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodPost)
	ctx.Request.Header.SetCookie("session_id_test", "wrong")
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())

	// Check idle timeout not expired
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.Header.SetCookie("session_id_test", token)
	ctx.Request.SetRequestURI("/isFresh")
	h(ctx)
	require.Equal(t, velocity.StatusInternalServerError, ctx.Response.StatusCode())

	// Test idle timeout
	time.Sleep(1200 * time.Millisecond)
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	ctx.Request.Header.SetCookie("session_id_test", token)
	ctx.Request.SetRequestURI("/isFresh")
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
}

func Test_Session_Next(t *testing.T) {
	t.Parallel()

	var (
		doNext bool
		muNext sync.RWMutex
	)

	app := velocity.New()

	app.Use(New(Config{
		Next: func(_ velocity.Ctx) bool {
			muNext.RLock()
			defer muNext.RUnlock()
			return doNext
		},
	}))

	app.Get("/", func(c velocity.Ctx) error {
		sess := FromContext(c)
		if sess == nil {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		id := sess.ID()
		return c.SendString("value=" + id)
	})

	h := app.Handler()

	// Test with Next returning false
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
	// Get session cookie
	token := string(ctx.Response.Header.Peek(velocity.HeaderSetCookie))
	require.NotEmpty(t, token, "Expected Set-Cookie header to be present")
	tokenParts := strings.SplitN(strings.SplitN(token, ";", 2)[0], "=", 2)
	require.Len(t, tokenParts, 2, "Expected Set-Cookie header to contain a token")
	token = tokenParts[1]
	require.Equal(t, "value="+token, string(ctx.Response.Body()))

	// Test with Next returning true
	muNext.Lock()
	doNext = true
	muNext.Unlock()

	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	h(ctx)
	require.Equal(t, velocity.StatusInternalServerError, ctx.Response.StatusCode())
}

func Test_Session_Middleware_Store(t *testing.T) {
	t.Parallel()
	app := velocity.New()

	handler, sessionStore := NewWithStore()

	app.Use(handler)

	app.Get("/", func(c velocity.Ctx) error {
		sess := FromContext(c)
		st := sess.Store()
		if st != sessionStore {
			return c.SendStatus(velocity.StatusInternalServerError)
		}
		return c.SendStatus(velocity.StatusOK)
	})

	h := app.Handler()

	// Test GET request
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(velocity.MethodGet)
	h(ctx)
	require.Equal(t, velocity.StatusOK, ctx.Response.StatusCode())
}
