package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/internal/pkg/apperror"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"redditclone/internal/domain/user"
	"redditclone/internal/pkg/session"
)

func (s *ApiTestSuite) Test00Identity_Register() {
	var result interface{}
	expectedStatus := http.StatusCreated
	require := require.New(s.T())
	assert := assert.New(s.T())

	newUser := user.New()
	newUser.Name = s.entities.user.Name
	newUser.Passhash = s.entities.user.Passhash

	s.repositoryMocks.user.On("Create", mock.Anything, mock.Anything).Return(error(nil))
	s.repositoryMocks.session.On("NewEntity", mock.Anything, uint(0)).Return(session.New(), error(nil))
	s.repositoryMocks.session.On("Create", mock.Anything, mock.Anything).Return(error(nil))

	reqBody := strings.NewReader(`{
	"username": "demo1",
	"password": "demo1"
}`)
	uri := "/api/register"

	req, _ := http.NewRequest(http.MethodPost, s.server.URL+uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	require.NoErrorf(err, "request error: %v", err)
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	require.NoErrorf(err, "read body error: %v", err)

	assert.Equalf(expectedStatus, resp.StatusCode, "expected http status %v, got %v", expectedStatus, resp.StatusCode)

	err = json.Unmarshal(resBody, &result)
	require.NoErrorf(err, "can not unpack json, error: %v", err)

	m, ok := result.(map[string]interface{})
	require.Truef(ok, "result %v cant assign to map[string]interface{}", result)

	_, ok = m["token"]
	require.Truef(ok, "result %v do not contain a token", m)
}

func (s *ApiTestSuite) Test01Identity_Login() {
	var result interface{}
	expectedStatus := http.StatusOK
	require := require.New(s.T())
	assert := assert.New(s.T())

	searchedUser := user.New()
	searchedUser.Name = s.entities.user.Name
	newSession := &session.Session{
		UserID: s.entities.user.ID,
		User:   *s.entities.user,
		Data: session.Data{
			UserID:              s.entities.user.ID,
			UserName:            s.entities.user.Name,
			ExpirationTokenTime: time.Now().Local().Add(time.Hour),
		},
	}

	s.repositoryMocks.user.On("First", mock.Anything, searchedUser).Return(s.entities.user, error(nil))

	s.repositoryMocks.session.On("Get", mock.Anything, s.entities.user.ID).Return(nil, apperror.ErrNotFound)
	s.repositoryMocks.session.On("NewEntity", mock.Anything, s.entities.user.ID).Return(newSession, error(nil))
	s.repositoryMocks.session.On("Create", mock.Anything, mock.Anything).Return(error(nil))

	reqBody := strings.NewReader(`{
	"username": "demo1",
	"password": "demo1"
}`)
	uri := "/api/login"

	req, _ := http.NewRequest(http.MethodPost, s.server.URL+uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	require.NoErrorf(err, "request error: %v", err)
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	require.NoErrorf(err, "read body error: %v", err)

	assert.Equalf(expectedStatus, resp.StatusCode, "expected http status %v, got %v", expectedStatus, resp.StatusCode)

	err = json.Unmarshal(resBody, &result)
	require.NoErrorf(err, "can not unpack json, error: %v", err)

	m, ok := result.(map[string]interface{})
	require.Truef(ok, "result %v cant assign to map[string]interface{}", result)

	token, ok := m["token"]
	require.Truef(ok, "result %v do not contain a token", m)

	s.token, ok = token.(string)
	require.Truef(ok, "can not assign to string token %v", token)
}
