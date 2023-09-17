package controller

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/zeynab-sb/geoolocation"
	"github.com/zeynab-sb/geoolocation/repository"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type LocationTestSuite struct {
	suite.Suite
	e        *echo.Echo
	endpoint string
	ctx      context.Context
	patch    *gomonkey.Patches
	location Location
}

func (suite *LocationTestSuite) SetupSuite() {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		log.Fatal("error in new connection", zap.Error(err))
	}

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	suite.e = echo.New()
	suite.endpoint = "/locations/"
	suite.ctx = context.Background()
	repo := repository.NewLocationRepository(mockDB)
	suite.location = Location{Geo: &geoolocation.Geo{Repository: repo}}

	suite.patch = gomonkey.NewPatches()
}

func (suite *LocationTestSuite) TearDownSuit() {
	suite.patch.Reset()
}

func (suite *LocationTestSuite) CallHandler(ip string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, suite.endpoint, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, rec)
	ctx.SetParamNames("ip")
	ctx.SetParamValues(ip)
	err := suite.location.GetByIP(ctx)

	return rec, err
}

func (suite *LocationTestSuite) TestGetByIP_GetByIP_EmptyIP_Failure() {
	require := suite.Require()
	expectedError := "code=400, message=IP is required"

	_, err := suite.CallHandler("")

	require.EqualError(err, expectedError)
}

func (suite *LocationTestSuite) TestGetByIP_GetByIP_InvalidIP_Failure() {
	require := suite.Require()
	expectedError := "code=400, message=Invalid ip"

	_, err := suite.CallHandler("127.0")

	require.EqualError(err, expectedError)
}

func (suite *LocationTestSuite) TestGetByIP_GetByIP_RepositoryError_Failure() {
	require := suite.Require()
	expectedError := "code=500, message=Internal Server Error"

	suite.patch.ApplyMethodReturn(suite.location.Geo.Repository, "GetLocationByIP", nil, errors.New("error"))

	_, err := suite.CallHandler("163.123.45.12")

	require.EqualError(err, expectedError)
}

func (suite *LocationTestSuite) TestGetByIP_GetByIP_LocationNotFound_Failure() {
	require := suite.Require()
	expectedError := "code=404, message=Location with specified ip not found"

	suite.patch.ApplyMethodReturn(suite.location.Geo.Repository, "GetLocationByIP", &repository.Location{}, nil)

	_, err := suite.CallHandler("163.123.45.12")

	require.EqualError(err, expectedError)
}

func (suite *LocationTestSuite) TestGetByIP_GetByIP_Success() {
	require := suite.Require()
	expectedResponse := `{"ip_address":"163.123.45.12","country":"test","city":"test","lat":1.867547585,"lng":1.867547585,"mystery_value":1356535367}
`

	suite.patch.ApplyMethodReturn(suite.location.Geo.Repository, "GetLocationByIP", &repository.Location{
		ID:           1,
		IPAddress:    "163.123.45.12",
		CountryCode:  "AB",
		Country:      "test",
		City:         "test",
		Lat:          1.867547585,
		Lng:          1.867547585,
		MysteryValue: 1356535367,
		UpdatedAt:    time.Time{},
		CreatedAt:    time.Time{},
	}, nil)

	res, err := suite.CallHandler("163.123.45.12")

	require.NoError(err)
	require.Equal(http.StatusOK, res.Code)
	require.Equal(expectedResponse, res.Body.String())
}

func TestGetByIP(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(LocationTestSuite))
}
