package user

import (
	"fmt"
	"pnas/category"
	"pnas/db"
	"pnas/ptype"
	"pnas/utils"
	"sync"

	"github.com/pkg/errors"
)

const (
	AdminId = 1
)

const (
	AuthAdmin = iota + 0
	AuthDownload
	AuthCreateVideoLib
	AuthWatchVideo

	AuthMax
)

type UserBaseInfo struct {
	Id              ptype.UserID
	Name            string
	Email           string
	HomeDirectoryId ptype.CategoryID
}

type User struct {
	mtx      sync.Mutex
	userInfo UserBaseInfo
	auth     utils.AuthBitSet
}

type NewUserParams struct {
	Email           string
	Name            string
	Passwd          string
	Auth            utils.AuthBitSet
	CategoryService category.IService
}

func NewUser(params *NewUserParams) error {
	sql := "call new_user(?, ?, ?, ?, ?, @new_user_id, @new_home_id)"
	byteAuth, err := params.Auth.MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}
	homeAuth := utils.NewBitSet(category.AuthMax)
	homeByteAuth, err := homeAuth.MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}
	var userId ptype.UserID
	var homeId ptype.CategoryID
	err = db.QueryRow(sql,
		params.Name,
		params.Email,
		params.Passwd,
		byteAuth,
		homeByteAuth,
	).Scan(&userId, &homeId)

	if err != nil {
		return errors.WithStack(err)
	}
	params.CategoryService.RefreshItem(category.UsersId)
	return nil
}

func UsedEmail(email string) bool {
	sql := "select count(*) from pnas.user where email=?"
	var c int
	err := db.QueryRow(sql, email).Scan(&c)
	if err != nil {
		return true
	}
	if c > 0 {
		return true
	}
	return false
}

func loadUser(userId ptype.UserID) (*User, error) {
	var user User
	user.userInfo.Id = userId
	sql := "select name, email, auth, directory_id from pnas.user where id=?"
	var byteAuth []byte
	err := db.QueryRow(sql, userId).Scan(
		&user.userInfo.Name,
		&user.userInfo.Email,
		&byteAuth,
		&user.userInfo.HomeDirectoryId,
	)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("userId: %d", userId))
	}
	user.auth = utils.NewBitSet(AuthMax)
	user.auth.UnmarshalBinary(byteAuth)
	return &user, nil
}

func (user *User) ChangeUserName(name string) error {
	sql := "update pnas.user set name=? where id=?"
	_, err := db.Exec(sql, name, user.userInfo.Id)
	if err != nil {
		return err
	}

	user.mtx.Lock()
	defer user.mtx.Unlock()
	user.userInfo.Name = name
	return nil
}

func (user *User) GetHomeDirectoryId() ptype.CategoryID {
	user.mtx.Lock()
	defer user.mtx.Unlock()
	return user.userInfo.HomeDirectoryId
}

func (user *User) GetUserInfo() UserBaseInfo {
	user.mtx.Lock()
	defer user.mtx.Unlock()
	return user.userInfo
}

func (user *User) GetUserName() string {
	user.mtx.Lock()
	defer user.mtx.Unlock()
	return user.userInfo.Name
}
