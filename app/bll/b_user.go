package bll

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/casbin/casbin"
	"github.com/pkg/errors"
	"moddns/app/models"
	"moddns/app/schema"
	"moddns/app/util"
)

// User 用户管理
type User struct {
	UserModel models.IUser     `inject:"IUser"`
	RoleModel models.IRole     `inject:"IRole"`
	Enforcer  *casbin.Enforcer `inject:""`
}

// QueryPage 查询分页数据
func (a *User) QueryPage(ctx context.Context, params schema.UserQueryParam, pageIndex, pageSize uint) (int64, []*schema.UserQueryResult, error) {
	total, items, err := a.UserModel.QueryPage(ctx, params, pageIndex, pageSize)
	if err != nil {
		return 0, nil, err
	}

	for i, item := range items {
		user, err := a.UserModel.Get(ctx, item.RecordID, true)
		if err == nil && user != nil && len(user.RoleIDs) > 0 {
			roleItems, err := a.RoleModel.QuerySelect(ctx, schema.RoleSelectQueryParam{RecordIDs: user.RoleIDs})
			if err == nil {
				roleNames := make([]string, len(roleItems))
				for i, item := range roleItems {
					roleNames[i] = item.Name
				}
				items[i].RoleNames = roleNames
			}
		}
	}

	return total, items, nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, recordID string) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, recordID, true)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, util.ErrNotFound
	}

	// 查询不返回密码
	item.Password = ""
	return item, nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item *schema.User) error {
	exists, err := a.UserModel.CheckUserName(ctx, item.UserName)
	if err != nil {
		return err
	} else if exists {
		return errors.New("用户名已经存在")
	}

	item.Password = util.SHA1HashString(item.Password)
	item.ID = 0
	item.RecordID = uuid.New().String()
	item.Created = time.Now().Unix()
	item.Deleted = 0
	err = a.UserModel.Create(ctx, item)
	if err != nil {
		return err
	}

	return a.LoadPolicy(ctx, item.RecordID)
}

// Update 更新数据
func (a *User) Update(ctx context.Context, recordID string, item *schema.User) error {
	oldItem, err := a.UserModel.Get(ctx, recordID, false)
	if err != nil {
		return err
	} else if oldItem == nil {
		return util.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		exists, err := a.UserModel.CheckUserName(ctx, item.UserName)
		if err != nil {
			return err
		} else if exists {
			return errors.New("用户名已经存在")
		}
	}

	info := util.StructToMap(item)
	delete(info, "id")
	delete(info, "record_id")
	delete(info, "creator")
	delete(info, "created")
	delete(info, "updated")
	delete(info, "deleted")
	delete(info, "password")

	if item.Password != "" {
		info["password"] = util.SHA1HashString(item.Password)
	}

	err = a.UserModel.UpdateWithRoleIDs(ctx, recordID, info, item.RoleIDs)
	if err != nil {
		return err
	}

	return a.LoadPolicy(ctx, recordID)
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, recordID string) error {
	exists, err := a.UserModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return util.ErrNotFound
	}

	err = a.UserModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}

	a.Enforcer.DeleteRolesForUser(recordID)
	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, recordID string, status int) error {
	exists, err := a.UserModel.Check(ctx, recordID)
	if err != nil {
		return err
	} else if !exists {
		return util.ErrNotFound
	}

	info := map[string]interface{}{
		"status": status,
	}
	err = a.UserModel.Update(ctx, recordID, info)
	if err != nil {
		return err
	}

	if status == 2 {
		a.Enforcer.DeleteRolesForUser(recordID)
	} else {
		err = a.LoadPolicy(ctx, recordID)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadAllPolicy 加载所有的用户策略
func (a *User) LoadAllPolicy() error {
	ctx := context.Background()

	userRoles, err := a.UserModel.QueryUserRoles(ctx, schema.UserRoleQueryParam{})
	if err != nil {
		return err
	}

	for _, ur := range userRoles {
		a.Enforcer.AddRoleForUser(ur.UserID, ur.RoleID)
	}

	return nil
}

// LoadPolicy 加载用户权限策略
func (a *User) LoadPolicy(ctx context.Context, userID string) error {
	userRoles, err := a.UserModel.QueryUserRoles(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return err
	}

	a.Enforcer.DeleteRolesForUser(userID)
	for _, ur := range userRoles {
		a.Enforcer.AddRoleForUser(ur.UserID, ur.RoleID)
	}
	return nil
}
