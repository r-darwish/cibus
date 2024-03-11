package internal

import (
	"fmt"
	"os"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

func login(username, password string) (*req.Client, error) {
	var auth struct {
		Company  string `json:"company"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	auth.Username = username
	auth.Password = password

	client := req.C()
	if _, ok := os.LookupEnv("CIBUS_DEV"); ok {
		client.DevMode()
	}
	client.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0")
	client.SetCommonHeaders(map[string]string{
		"Application-Id": "E5D5FEF5-A05E-4C64-AEBA-BA0CECA0E402",
	})

	response, err := client.R().SetBody(&auth).Post("https://api.capir.pluxee.co.il/auth/authToken")
	if err != nil {
		return nil, fmt.Errorf("Failed logging in: %w", err)
	}

	if response.IsErrorState() {
		return nil, fmt.Errorf("HTTP Error: %d", response.StatusCode)
	}

	return client, nil
}

func getFriends(client *req.Client, query string) error {
	var request struct {
		Query string `json:"query"`
		Type  string `json:"type"`
	}

	type User struct {
		UserId int `json:"user_id"`
	}

	var data struct {
		List []User `json:"list"`
	}

	request.Query = query
	request.Type = "autocomp_friend"
	response, err := client.R().SetBody(&request).SetSuccessResult(&data).Post("https://api.consumers.pluxee.co.il/api/main.py")
	if err != nil {
		return nil
	}

	if response.IsErrorState() {
		return nil
	}

	userIds := lo.Map(data.List, func(item User, _ int) int {
		return item.UserId
	})
	if len(userIds) == 0 {
		return nil
	}

	var addFriendRequest struct {
		Type  string `json:"type"`
		Users []int  `json:"user_id_list"`
	}
	addFriendRequest.Type = "prx_share_user"
	addFriendRequest.Users = userIds

	request.Query = query
	request.Type = "autocomp_friend"
	response, err = client.R().SetBody(&addFriendRequest).Post("https://api.consumers.pluxee.co.il/api/main.py")
	if err != nil {
		return nil
	}

	if response.IsErrorState() {
		return nil
	}

	return nil
}

func AddAllFriends(username, password string) error {
	client, err := login(username, password)
	if err != nil {
		return err
	}

	err = getFriends(client, ".io")
	if err != nil {
		return err
	}

	err = getFriends(client, "com")
	if err != nil {
		return err
	}

	err = getFriends(client, "security")
	if err != nil {
		return err
	}

	return nil
}
