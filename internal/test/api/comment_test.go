package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"redditclone/internal/domain/comment"
)

func (s *ApiTestSuite) TestComment_Create() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	newComment := &comment.Comment{}
	*newComment = *s.entities.comment
	newComment.ID = ""

	s.repositoryMocks.comment.On("Create", mock.Anything, newComment).Return(error(nil))
	s.repositoryMocks.post.On("Get", mock.Anything, newComment.PostID).Return(s.entities.post, error(nil))

	b, err := json.Marshal(newComment)
	require.NoErrorf(err, "can not json.Marshal() a value: %v, error", newComment, err)

	reqBody := bytes.NewReader(b)
	uri := "/api/post/" + s.entities.comment.PostID
	expectedData := s.entities.post
	expectedStatus := http.StatusCreated

	req, _ := http.NewRequest(http.MethodPost, s.server.URL+uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token)
	resp, err := s.client.Do(req)
	require.NoErrorf(err, "request error: %v", err)
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	require.NoErrorf(err, "read body error: %v", err)

	assert.Equalf(expectedStatus, resp.StatusCode, "expected http status %v, got %v", expectedStatus, resp.StatusCode)

	err = json.Unmarshal(resBody, &result)
	require.NoErrorf(err, "can not unpack json, error: %v", err)

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	assert.Equalf(expected, result, "results not match\nGot: %#v\nExpected: %#v", result, expectedData)
}

func (s *ApiTestSuite) TestComment_Delete() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	newComment := &comment.Comment{}
	*newComment = *s.entities.comment
	newComment.ID = ""

	s.repositoryMocks.comment.On("Delete", mock.Anything, s.entities.comment.ID).Return(error(nil))
	s.repositoryMocks.post.On("Get", mock.Anything, newComment.PostID).Return(s.entities.post, error(nil))

	b, err := json.Marshal(newComment)
	require.NoErrorf(err, "can not json.Marshal() a value: %v, error", newComment, err)

	reqBody := bytes.NewReader(b)
	uri := "/api/post/" + s.entities.comment.PostID + "/" + s.entities.comment.ID
	expectedData := s.entities.post
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodDelete, s.server.URL+uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token)
	resp, err := s.client.Do(req)
	require.NoErrorf(err, "request error: %v", err)
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	require.NoErrorf(err, "read body error: %v", err)

	assert.Equalf(expectedStatus, resp.StatusCode, "expected http status %v, got %v", expectedStatus, resp.StatusCode)

	err = json.Unmarshal(resBody, &result)
	require.NoErrorf(err, "can not unpack json, error: %v", err)

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	assert.Equalf(expected, result, "results not match\nGot: %#v\nExpected: %#v", result, expectedData)
}
