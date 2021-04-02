package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/user"
	"redditclone/internal/domain/vote"
	"redditclone/internal/pkg/apperror"
	"redditclone/internal/pkg/errorshandler"

	"github.com/minipkg/selection_condition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *ApiTestSuite) TestPost_Create() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	newPost := &post.Post{}
	*newPost = *s.entities.post
	newPost.ID = ""
	newPost.Comments = nil
	newPost.Votes = nil

	s.repositoryMocks.post.On("Create", mock.Anything, newPost).Return(error(nil))

	b, err := json.Marshal(newPost)
	require.NoErrorf(err, "can not json.Marshal() a value: %v, error", newPost, err)

	reqBody := bytes.NewReader(b)
	uri := "/api/posts"
	expectedData := newPost
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
	require.NoErrorf(err, "can not unpack json %q, error: %v", string(resBody), err)

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	assert.Equalf(expected, result, "results not match\nGot: %#v\nExpected: %#v", result, expectedData)
}

func (s *ApiTestSuite) TestPost_Delete() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	newPost := &post.Post{}
	*newPost = *s.entities.post
	newPost.ID = ""
	newPost.Comments = nil
	newPost.Votes = nil

	s.repositoryMocks.post.On("Delete", mock.Anything, s.entities.post.ID).Return(error(nil))

	uri := "/api/post/" + s.entities.post.ID
	expectedData := errorshandler.SuccessMessage()
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodDelete, s.server.URL+uri, nil)
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

func (s *ApiTestSuite) TestPost_Get() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	p := &post.Post{}
	*p = *s.entities.post
	p.Views++

	s.repositoryMocks.post.On("Get", mock.Anything, s.entities.post.ID).Return(s.entities.post, error(nil))
	s.repositoryMocks.post.On("Update", mock.Anything, p).Return(error(nil))

	uri := "/api/post/" + s.entities.post.ID
	expectedData := p
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodGet, s.server.URL+uri, nil)
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

func (s *ApiTestSuite) TestPost_List() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	list := []post.Post{*s.entities.post}
	query := selection_condition.SelectionCondition{
		Where: &post.Post{},
	}

	s.repositoryMocks.post.On("Query", mock.Anything, query).Return(list, error(nil))

	uri := "/api/posts"
	expectedData := list
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodGet, s.server.URL+uri, nil)
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

func (s *ApiTestSuite) TestPost_ListByCategory() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	list := []post.Post{*s.entities.post}
	query := selection_condition.SelectionCondition{
		Where: &post.Post{
			Category: "category",
		},
	}

	s.repositoryMocks.post.On("Query", mock.Anything, query).Return(list, error(nil))

	uri := "/api/posts/category"
	expectedData := list
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodGet, s.server.URL+uri, nil)
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

func (s *ApiTestSuite) TestPost_ListByUser() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	list := []post.Post{*s.entities.post}
	searchedUser := &user.User{
		Name: s.entities.user.Name,
	}
	query := selection_condition.SelectionCondition{
		Where: &post.Post{
			UserID: s.entities.user.ID,
		},
	}

	s.repositoryMocks.user.On("First", mock.Anything, searchedUser).Return(s.entities.user, error(nil))
	s.repositoryMocks.post.On("Query", mock.Anything, query).Return(list, error(nil))

	uri := "/api/user/" + s.entities.user.Name
	expectedData := list
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodGet, s.server.URL+uri, nil)
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

func (s *ApiTestSuite) TestPost_Upvote() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	newVote := s.api.Domain.Vote.Service.NewEntity(s.entities.vote.UserID, s.entities.vote.PostID, 1)
	newVote.User = *s.entities.user

	searchedVote := &vote.Vote{
		PostID: s.entities.vote.PostID,
		UserID: s.entities.vote.UserID,
	}
	p := &post.Post{}
	*p = *s.entities.post
	p.Score--

	s.repositoryMocks.vote.On("First", mock.Anything, searchedVote).Return(nil, apperror.ErrNotFound)
	s.repositoryMocks.vote.On("Create", mock.Anything, newVote).Return(error(nil))

	s.repositoryMocks.post.On("Get", mock.Anything, p.ID).Return(p, error(nil))
	s.repositoryMocks.post.On("Update", mock.Anything, s.entities.post).Return(error(nil))

	uri := "/api/post/" + s.entities.post.ID + "/upvote"
	expectedData := s.entities.post
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodGet, s.server.URL+uri, nil)
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

func (s *ApiTestSuite) TestPost_Downvote() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	newVote := s.api.Domain.Vote.Service.NewEntity(s.entities.vote.UserID, s.entities.vote.PostID, -1)
	newVote.User = *s.entities.user

	searchedVote := &vote.Vote{
		PostID: s.entities.vote.PostID,
		UserID: s.entities.vote.UserID,
	}
	p := &post.Post{}
	*p = *s.entities.post
	p.Score++

	s.repositoryMocks.vote.On("First", mock.Anything, searchedVote).Return(nil, apperror.ErrNotFound)
	s.repositoryMocks.vote.On("Create", mock.Anything, newVote).Return(error(nil))

	s.repositoryMocks.post.On("Get", mock.Anything, p.ID).Return(p, error(nil))
	s.repositoryMocks.post.On("Update", mock.Anything, s.entities.post).Return(error(nil))

	uri := "/api/post/" + s.entities.post.ID + "/downvote"
	expectedData := s.entities.post
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodGet, s.server.URL+uri, nil)
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

func (s *ApiTestSuite) TestPost_Unvote() {
	var result interface{}
	var expected interface{}
	require := require.New(s.T())
	assert := assert.New(s.T())
	s.setupSession()

	newVote := s.api.Domain.Vote.Service.NewEntity(s.entities.vote.UserID, s.entities.vote.PostID, 1)
	newVote.User = *s.entities.user

	searchedVote := &vote.Vote{
		PostID: s.entities.vote.PostID,
		UserID: s.entities.vote.UserID,
	}
	p := &post.Post{}
	*p = *s.entities.post
	p.Score -= s.entities.vote.Value

	s.repositoryMocks.vote.On("First", mock.Anything, searchedVote).Return(s.entities.vote, nil)
	s.repositoryMocks.vote.On("Delete", mock.Anything, s.entities.vote.ID).Return(error(nil))

	s.repositoryMocks.post.On("Get", mock.Anything, p.ID).Return(p, error(nil))
	s.repositoryMocks.post.On("Update", mock.Anything, p).Return(error(nil))

	uri := "/api/post/" + s.entities.post.ID + "/unvote"
	expectedData := p
	expectedStatus := http.StatusOK

	req, _ := http.NewRequest(http.MethodGet, s.server.URL+uri, nil)
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
