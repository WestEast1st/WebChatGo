package main

import (
	"errors"
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
