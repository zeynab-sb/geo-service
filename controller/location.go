package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/zeynab-sb/geoolocation"
	"net"
	"net/http"
)

type Location struct {
	Geo *geoolocation.Geo
}

type GetByIPRes struct {
	IpAddress    string  `json:"ip_address"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	MysteryValue int     `json:"mystery_value"`
}

func (l *Location) GetByIP(ctx echo.Context) error {
	ip := ctx.Param("ip")
	if ip == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "IP is required")
	}

	if net.ParseIP(ip) == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ip")
	}

	loc, err := l.Geo.Repository.GetLocationByIP(ip)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	if loc.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Location with specified ip not found")
	}

	return ctx.JSON(http.StatusOK, &GetByIPRes{
		IpAddress:    loc.IPAddress,
		Country:      loc.Country,
		City:         loc.City,
		Lat:          loc.Lat,
		Lng:          loc.Lng,
		MysteryValue: loc.MysteryValue,
	})
}
