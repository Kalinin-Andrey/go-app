package test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/session"
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Kalinin-Andrey/redditclone/pkg/config"

	"github.com/Kalinin-Andrey/redditclone/internal/app/api"
	"github.com/Kalinin-Andrey/redditclone/internal/test/mock"
)

type Client struct {
	client		*http.Client
	token		string
}

type App struct {
	app			*api.App
	sessionRepo	*mock.SessionRepository
	userRepo	*mock.UserRepository
	postRepo	*mock.PostRepository
	commentRepo	*mock.CommentRepository
	voteRepo	*mock.VoteRepository
}

var (
	c = Client{
		client:	&http.Client{},
	}
	a		App
	passhash, _ = hex.DecodeString("3a73acfdb534ddded4c0109383ee3e5a66314113d1ff691aaf4b3ee073c8fc2edd06d48f0555ec3783f4c479994e3eee3433734c29b05f08be0e9739b956b88d8fe872bd0a0942214e94fd4001e757fa3b66a2b9925de2e800c55ef49baa4c03")
	u		= &user.User{
		ID:			1,
		Name:		"demo1",
		Passhash:	string(passhash),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}
	co = &comment.Comment{
		ID:			10,
		PostID:		1,
		UserID:		1,
		User:		*u,
		Body:		"Who care about comments?",
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		DeletedAt:	nil,
	}
	v = &vote.Vote{
		ID:			10,
		PostID:		1,
		UserID:		1,
		User:		*u,
		Value:		1,
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		DeletedAt:	nil,
	}
	p = &post.Post{
		ID:			1,
		Score:		2,
		Views:		3,
		Title:		"What does a good programmer mean?",
		Type:		post.TypeText,
		Category:	post.CategoryProgramming,
		Text:		"Who can consider himself a good programmer?",
		UserID:		1,
		User:		user.User{
			ID:			1,
			Name:		"demo1",
			CreatedAt:	time.Now(),
			UpdatedAt:	time.Now(),
			DeletedAt:	nil,
		},
		Votes:		[]vote.Vote{*v},
		Comments:	[]comment.Comment{*co},
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		DeletedAt:	nil,
	}

)

func start() *httptest.Server {
	if a.app == nil {
		os.Chdir("../..")	//
		cfg, err := config.Get()
		if err != nil {
			log.Fatalln("Can not load the config")
		}
		cfg.Log.OutputPaths = []string{"stdout", "log/test_app.log"}

		a.app = api.New(NewCommonApp(*cfg), *cfg)

		a.sessionRepo	= a.app.SessionRepository.(*mock.SessionRepository)
		a.userRepo		= a.app.Domain.User.Repository.(*mock.UserRepository)
		a.postRepo		= a.app.Domain.Post.Repository.(*mock.PostRepository)
		a.commentRepo	= a.app.Domain.Comment.Repository.(*mock.CommentRepository)
		a.voteRepo		= a.app.Domain.Vote.Repository.(*mock.VoteRepository)

		a.userRepo.Response.First.Entity	= u
		a.userRepo.Response.First.Err		= nil
		a.userRepo.Response.Get.Entity		= u
		a.userRepo.Response.Get.Err			= nil
		a.userRepo.Response.Query.List		= []user.User{*u}
		a.userRepo.Response.Query.Err		= nil

		a.sessionRepo.Response.GetByUserID.Entity	= &session.Session{
			ID:			1,
			UserID:		1,
			User:		*u,
		}
		a.sessionRepo.Response.GetByUserID.Err	= nil
	}
	return httptest.NewServer(a.app.Server.Handler)
}

func TestRegister(t *testing.T) {
	caseName := "Register"
	var result		interface{}
	expectedStatus := http.StatusCreated

	ts := start()

	a.userRepo.Response.Create.Err	= nil
	reqBody := strings.NewReader(`{
	"username": "demo",
	"password": "demo"
}`)
	uri := "/api/register"

	req, _ := http.NewRequest(http.MethodPost, ts.URL + uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Errorf("[%s] result %v cant assign to map[string]interface{}", caseName, result)
		return
	}

	if _, ok := m["token"]; !ok {
		t.Errorf("[%s] result %v do not contain a token", caseName, m)
		return
	}

	ts.Close()
}

func TestLogin(t *testing.T) {
	caseName := "Login"
	var result		interface{}
	expectedStatus := http.StatusOK
	ts := start()

	reqBody := strings.NewReader(`{
	"username": "demo1",
	"password": "demo1"
}`)
	uri := "/api/login"

	req, _ := http.NewRequest(http.MethodPost, ts.URL + uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Errorf("[%s] result %v cant assign to map[string]interface{}", caseName, result)
		return
	}

	token, ok := m["token"]
	if !ok {
		t.Errorf("[%s] result %v do not contain a token", caseName, m)
		return
	}

	if c.token, ok = token.(string); !ok {
		t.Errorf("[%s] can not assign to string token %v", caseName, token)
		return
	}

	ts.Close()
}

func TestCreatePost(t *testing.T) {
	caseName := "CreatePost"
	var result		interface{}
	var expected	interface{}
	ts := start()

	reqBody := strings.NewReader(`{
	"category": "programming",
	"type": "text",
	"title": "What does a good programmer mean?",
	"text": "Who can consider himself a good programmer?"
}`)
	uri := "/api/posts"
	expectedData := p
	expectedStatus := http.StatusCreated

	a.postRepo.Response.Create.Entity	= p
	a.postRepo.Response.Create.Err		= nil

	req, _ := http.NewRequest(http.MethodPost, ts.URL + uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestDeletePost(t *testing.T) {
	caseName :=		"DeletePost"
	var result		interface{}
	var expected	interface{}
	ts := start()

	uri := "/api/post/1"
	expectedData	:= errorshandler.SuccessMessage()
	expectedStatus := http.StatusOK

	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil
	a.voteRepo.Response.Query.List		= []vote.Vote{*v}
	a.voteRepo.Response.Query.Err		= nil
	a.commentRepo.Response.Query.List	= []comment.Comment{*co}
	a.commentRepo.Response.Query.Err	= nil
	a.postRepo.Response.Delete.Err		= nil

	req, _ := http.NewRequest(http.MethodDelete, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestGetPost(t *testing.T) {
	caseName :=		"GetPost"
	var result		interface{}
	var expected	interface{}
	ts := start()

	uri := "/api/post/1"
	expectedStatus := http.StatusOK
	post := *p
	a.postRepo.Response.Get.Entity		= &post
	a.postRepo.Response.Get.Err			= nil
	a.voteRepo.Response.Query.List		= []vote.Vote{*v}
	a.voteRepo.Response.Query.Err		= nil
	a.commentRepo.Response.Query.List	= []comment.Comment{*co}
	a.commentRepo.Response.Query.Err	= nil
	a.postRepo.Response.Update.Entity	= &post
	a.postRepo.Response.Update.Err		= nil
	p.Views++
	expectedData	:= p

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestGetPosts(t *testing.T) {
	caseName :=		"GetPosts"
	var result		interface{}
	var expected	interface{}
	ts := start()

	uri := "/api/posts"
	expectedData	:= []post.Post{*p}
	expectedStatus := http.StatusOK
	a.postRepo.Response.Query.List		= []post.Post{*p}
	a.postRepo.Response.Query.Err		= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestGetPostsByCategory(t *testing.T) {
	caseName :=		"GetPostsByCategory"
	var result		interface{}
	var expected	interface{}
	ts := start()

	uri := "/api/posts"
	expectedData	:= []post.Post{*p}
	expectedStatus := http.StatusOK
	a.postRepo.Response.Query.List		= []post.Post{*p}
	a.postRepo.Response.Query.Err		= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestGetPostsByUser(t *testing.T) {
	caseName :=		"GetPostsByUser"
	var result		interface{}
	var expected	interface{}
	ts := start()

	uri := "/api/posts"
	expectedData	:= []post.Post{*p}
	expectedStatus := http.StatusOK
	a.postRepo.Response.Query.List		= []post.Post{*p}
	a.postRepo.Response.Query.Err		= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestCreateComment(t *testing.T) {
	caseName := "CreateComment"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusCreated
	ts := start()

	reqBody := strings.NewReader(`{
	"body": "Who care about comments?"
}`)
	uri := "/api/post/1"
	expectedData := p

	a.commentRepo.Response.Create.Entity	= co
	a.commentRepo.Response.Create.Err		= nil
	a.postRepo.Response.Get.Entity			= p
	a.postRepo.Response.Get.Err				= nil

	req, _ := http.NewRequest(http.MethodPost, ts.URL + uri, reqBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestDeleteComment(t *testing.T) {
	caseName := "DeleteComment"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/10"
	expectedData := p

	a.commentRepo.Response.Delete.Err		= nil
	a.postRepo.Response.Get.Entity			= p
	a.postRepo.Response.Get.Err				= nil

	req, _ := http.NewRequest(http.MethodDelete, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestUpvoteCreate(t *testing.T) {
	caseName := "UpvoteCreate"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/upvote"
	expectedData := p

	a.voteRepo.Response.First.Entity	= nil
	a.voteRepo.Response.First.Err		= apperror.ErrNotFound

	a.voteRepo.Response.Create.Err		= nil

	a.postRepo.Response.Update.Entity	= p
	a.postRepo.Response.Update.Err		= nil

	p.Score++
	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestDownvoteCreate(t *testing.T) {
	caseName := "DownvoteCreate"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/downvote"
	expectedData := p

	a.voteRepo.Response.First.Entity	= nil
	a.voteRepo.Response.First.Err		= apperror.ErrNotFound

	a.voteRepo.Response.Create.Err		= nil

	a.postRepo.Response.Update.Entity	= p
	a.postRepo.Response.Update.Err		= nil

	p.Score--
	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestUpvoteSecond(t *testing.T) {
	caseName := "UpvoteSecond"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/upvote"
	expectedData := p

	a.voteRepo.Response.First.Entity	= v
	a.voteRepo.Response.First.Err		= nil

	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestDownvoteSecond(t *testing.T) {
	caseName := "DownvoteSecond"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/downvote"
	expectedData := p

	a.voteRepo.Response.First.Entity	= v
	a.voteRepo.Response.First.Err		= nil

	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestUpvoteAfterDownvote(t *testing.T) {
	caseName := "UpvoteAfterDownvote"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/upvote"
	expectedData := p
	vo := *v
	vo.Value = -1
	a.voteRepo.Response.First.Entity	= &vo
	a.voteRepo.Response.First.Err		= nil

	a.voteRepo.Response.Update.Err		= nil

	a.postRepo.Response.Update.Entity	= p
	a.postRepo.Response.Update.Err		= nil

	p.Score += 2
	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestDownvoteAfterUpvote(t *testing.T) {
	caseName := "DownvoteAfterUpvote"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/downvote"
	expectedData := p
	vo := *v
	vo.Value = -1
	a.voteRepo.Response.First.Entity	= &vo
	a.voteRepo.Response.First.Err		= nil

	a.voteRepo.Response.Update.Err		= nil

	a.postRepo.Response.Update.Entity	= p
	a.postRepo.Response.Update.Err		= nil

	p.Score -= 2
	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestUnvoteAfterUpvote(t *testing.T) {
	caseName := "UnvoteAfterUpvote"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/unvote"
	expectedData := p

	a.voteRepo.Response.First.Entity	= v
	a.voteRepo.Response.First.Err		= nil

	a.voteRepo.Response.Delete.Err		= nil

	a.postRepo.Response.Update.Entity	= p
	a.postRepo.Response.Update.Err		= nil

	p.Score--
	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

func TestUnvoteAfterDownvote(t *testing.T) {
	caseName := "UnvoteAfterDownvote"
	var result		interface{}
	var expected	interface{}
	expectedStatus := http.StatusOK
	ts := start()

	uri := "/api/post/1/unvote"
	expectedData := p
	v.Value = -1
	a.voteRepo.Response.First.Entity	= v
	a.voteRepo.Response.First.Err		= nil

	a.voteRepo.Response.Delete.Err		= nil

	a.postRepo.Response.Update.Entity	= p
	a.postRepo.Response.Update.Err		= nil

	p.Score++
	a.postRepo.Response.Get.Entity		= p
	a.postRepo.Response.Get.Err			= nil

	req, _ := http.NewRequest(http.MethodGet, ts.URL + uri, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + c.token)
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("[%s] request error: %v", caseName, err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Errorf("[%s] expected http status %v, got %v", caseName, expectedStatus, resp.StatusCode)
	}

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		t.Errorf("[%s] cant unpack json: %v", caseName, err)
		return
	}

	jsonData, err := json.Marshal(expectedData)
	json.Unmarshal(jsonData, &expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("[%s] results not match\nGot: %#v\nExpected: %#v", caseName, result, expectedData)
	}

	ts.Close()
}

