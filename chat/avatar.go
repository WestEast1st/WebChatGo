package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("char: アバターのURLを取得できません。")

// Avatarはユーザーのプロフィール画像を表す型です
type Avatar interface {
	// GetAvatarURLは使用されたクライアントのアバターのURLを返します
	// 問題が発生した場合にはエラーを返します。特にURLを取得できなかった場合にはErrNoAvatarURLを返す。
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

// AvatarURLを取得
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlstr, ok := url.(string); ok {
			return urlstr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	/*
		if email, ok := c.userData["email"]; ok {
			if emailstr, ok := email.(string); ok {
				m := md5.New()
				io.WriteString(m, strings.ToLower(emailstr))
				return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
			}
		}
	*/
	if userid, ok := c.userData["userid"]; ok {
		if useridstr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridstr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["userid"]; ok {
		if userIdStr, ok := userID.(string); ok {
			// return "/avatars/" + useridstr + ".jpg", nil
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				fmt.Println(files)
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := filepath.Match(userIdStr+"*", file.Name()); match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
