package service

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
)

func InitUI(s *Server) {
	if s.store.CoreConfig.DisableUI {
		return
	}

	// Remove LocateAppend and LocateFS because we don't use it and it cause tests issue
	riceConfig := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded},
	}
	// Never use constant or variable according to docs : https://github.com/GeertJohan/go.rice#calling-findbox-and-mustfindbox
	uiAssets, err := riceConfig.FindBox("../ui/dist")
	if err != nil {
		panic("static ui/dist not found. Build them with `cd ui && npm build` first.")
	}

	assetHandler := http.FileServer(uiAssets.HTTPBox())
	s.Echo.GET("/", echo.WrapHandler(assetHandler))
	s.Echo.GET("/index*.js", echo.WrapHandler(assetHandler))
	s.Echo.GET("/style*", echo.WrapHandler(assetHandler))
	s.Echo.GET("/favicon*", echo.WrapHandler(assetHandler))
	s.Echo.GET("/icons.svg", echo.WrapHandler(assetHandler))
	s.Echo.GET("/fonts/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
}
