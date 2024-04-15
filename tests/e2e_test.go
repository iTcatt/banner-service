package main

import (
	"banner-service/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
	"testing"
)

var params = model.BannerParams{
	TagIDs:    []int{rand.Int()},
	FeatureID: rand.Int(),
	Content:   `{"title":"test"}`,
	IsActive:  true,
}

func getAdminToken() (string, error) {
	resp, err := http.Get("http://localhost:8888/admin")
	if err != nil {
		return "", err
	}
	if http.StatusOK != resp.StatusCode {
		return "", errors.New("wrong status code")
	}

	var token struct {
		Token string `json:"token"`
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(body, &token); err != nil {
		return "", err
	}
	return token.Token, nil
}

func createBanner(token string) (int, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8888/banner", bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusCreated {
		return 0, errors.New("wrong status code")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	res := struct {
		BannerID int `json:"banner_id"`
	}{}
	if err = json.Unmarshal(body, &res); err != nil {
		return 0, err
	}
	return res.BannerID, nil
}

func TestCreateBannerWithoutToken(t *testing.T) {
	data, err := json.Marshal(params)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8888/banner", bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetBannersWithoutToken(t *testing.T) {
	resp, err := http.Get("http://localhost:8888/user_banner/100")
	if err != nil {
		t.Errorf("Get request without header failed: %v", err)
	}
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetBanners(t *testing.T) {
	token, err := getAdminToken()
	if err != nil {
		t.Error(err)
	}
	_, err = createBanner(token)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", "http://localhost:8888/user_banner", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	q := req.URL.Query()
	q.Add("tag_id", fmt.Sprint(params.TagIDs[0]))
	q.Add("feature_id", fmt.Sprint(params.FeatureID))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("token", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Get request failed: %v", err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll error: %v", err)
	}
	result := strings.TrimRight(string(body), "\n")
	assert.Equal(t, params.Content, result)
}
