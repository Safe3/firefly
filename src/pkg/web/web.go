package web

import (
	"embed"
	"fahi/pkg/config"
	"fahi/pkg/util"
	"fahi/pkg/wg"
	"fmt"
	"image/png"
	"io/fs"
	"net/http"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

//go:embed www
var embededFiles embed.FS

func assetHandler() http.Handler {
	fsys, err := fs.Sub(embededFiles, "www")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}

func Serve(cfg *config.Config, wgIface *wg.WgIface) error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(util.GenerateKey()))))

	e.GET("/*", echo.WrapHandler(assetHandler()))

	e.GET("/api/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &cfg.Version)
	})

	e.GET("/api/lang", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &cfg.Lang)
	})

	e.GET("/api/session", func(c echo.Context) error {
		var result struct {
			RequiresPassword bool `json:"requiresPassword"`
			Authenticated    bool `json:"authenticated"`
		}

		if cfg.Password != "" {
			result.RequiresPassword = true
			if verify(c) {
				result.Authenticated = true
			}
		}

		return c.JSON(http.StatusOK, &result)
	})

	e.POST("/api/session", func(c echo.Context) error {
		var req struct {
			Password string `json:"password"`
		}

		var result struct {
			Error string `json:"error"`
		}

		if err := c.Bind(&req); err != nil {
			result.Error = "Missing password"
			return c.JSON(http.StatusUnauthorized, &result)
		} else if req.Password != cfg.Password {
			result.Error = "Incorrect password"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		sess, _ := session.Get("connect.sid", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
		}
		sess.Values["authenticated"] = true
		sess.Save(c.Request(), c.Response())
		return c.NoContent(http.StatusNoContent)
	})

	e.DELETE("/api/session", func(c echo.Context) error {
		sess, err := session.Get("connect.sid", c)
		if err == nil {
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
		}
		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/api/wireguard/client", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		peers, err := wgIface.GetPeers()
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		return c.JSON(http.StatusOK, peers)
	})

	e.GET("/api/wireguard/client/:clientId/qrcode.svg", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		peerConfig, err := wgIface.GetPeerConfig(c.Param("clientId"))
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		c.Response().Header().Set("Content-Type", "image/png")
		qrCode, err := qr.Encode(peerConfig, qr.M, qr.Auto)
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}
		qrCode, err = barcode.Scale(qrCode, 200, 200)
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}
		png.Encode(c.Response().Writer, qrCode)

		return nil
	})

	e.GET("/api/wireguard/client/:clientId/configuration", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		clientId := c.Param("clientId")
		peerConfig, err := wgIface.GetPeerConfig(clientId)
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.conf"`, strings.ReplaceAll(clientId, "-", "")))

		return c.Blob(http.StatusOK, "text/plain", []byte(peerConfig))
	})

	e.POST("/api/wireguard/client", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		var req struct {
			Name string `json:"name"`
		}

		if err := c.Bind(&req); err != nil {
			result.Error = "Missing name"
			return c.JSON(http.StatusForbidden, &result)
		}

		peer, err := wgIface.AddPeer(req.Name)
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		return c.JSON(http.StatusOK, peer)
	})

	e.DELETE("/api/wireguard/client/:clientId", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		clientId := c.Param("clientId")

		err := wgIface.DelPeer(clientId)
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.POST("/api/wireguard/client/:clientId/:enabled", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		err := wgIface.SetPeer(c.Param("clientId"), c.Param("enabled"), "", "")
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.PUT("/api/wireguard/client/:clientId/name", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		var req struct {
			ClientId string `param:"clientId"`
			Name     string `json:"name"`
		}

		if err := c.Bind(&req); err != nil {
			result.Error = "Missing name"
			return c.JSON(http.StatusForbidden, &result)
		}

		err := wgIface.SetPeer(req.ClientId, "", req.Name, "")
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.PUT("/api/wireguard/client/:clientId/address", func(c echo.Context) error {
		var result struct {
			Error string `json:"error"`
		}

		if !verify(c) {
			result.Error = "Not logged in"
			return c.JSON(http.StatusUnauthorized, &result)
		}

		var req struct {
			ClientId string `param:"clientId"`
			Address  string `json:"address"`
		}

		if err := c.Bind(&req); err != nil {
			result.Error = "Missing address"
			return c.JSON(http.StatusForbidden, &result)
		}

		err := wgIface.SetPeer(req.ClientId, "", "", req.Address)
		if err != nil {
			result.Error = err.Error()
			return c.JSON(http.StatusInternalServerError, &result)
		}

		return c.NoContent(http.StatusNoContent)
	})

	address := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	if cfg.AutoSSL {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(cfg.Host)
		e.AutoTLSManager.Cache = autocert.DirCache(util.RootDir + "cert")
		return e.StartAutoTLS(address)
	} else {
		return e.Start(address)
	}
}

func verify(c echo.Context) bool {
	sess, err := session.Get("connect.sid", c)
	if err != nil {
		return false
	}

	authenticated := sess.Values["authenticated"]
	if authenticated == nil {
		return false
	} else if value, ok := authenticated.(bool); ok {
		if value {
			return true
		}
	}
	return false
}
