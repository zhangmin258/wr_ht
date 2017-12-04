package models

import (
	"github.com/astaxie/beego/orm"
)

type SysRole struct {
	Id          int
	Displayname string
	Remark      string
	Org_id      int
}

func SysRoleList() (list []SysRole, err error) {
	sql := `SELECT * FROM sys_role`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}

func SysRoleByRid(rid int) (role *SysRole, err error) {
	sql := `SELECT * FROM sys_role WHERE id=?`
	err = orm.NewOrm().Raw(sql, rid).QueryRow(&role)
	return
}

func DelRole(rid int) error {
	sql := `DELETE FROM sys_role WHERE id=?`
	o := orm.NewOrm()
	o.Begin()
	_, err := o.Raw(sql, rid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	sql = `DELETE FROM sys_role_menu WHERE role_id=?`
	_, err = o.Raw(sql, rid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func (sr *SysRole) Insert(menu_ids []string) error {
	sql := `INSERT INTO sys_role (displayname, remark, org_id)
			values(?, ?, ?)`
	o := orm.NewOrm()
	o.Begin()
	res, err := o.Raw(sql, sr.Displayname, sr.Remark, sr.Org_id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	if len(menu_ids) > 0 {
		rid, err := res.LastInsertId()
		if err != nil {
			o.Rollback()
			return err
		}

		sql = ` INSERT INTO sys_role_menu (role_id, menu_id) values(?, ?)`
		for i := 0; i < len(menu_ids); i++ {
			_, err = o.Raw(sql, rid, menu_ids[i]).Exec()
			if err != nil {
				break
			}
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func (sr *SysRole) Update(menu_ids []string) error {
	sql := `UPDATE sys_role SET displayname=?, remark=?, org_id=?
			WHERE id=?`
	o := orm.NewOrm()
	o.Begin()
	_, err := o.Raw(sql, sr.Displayname, sr.Remark, sr.Org_id, sr.Id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	sql = `DELETE FROM sys_role_menu WHERE role_id=?`
	_, err = o.Raw(sql, sr.Id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	if len(menu_ids) > 0 {
		sql = ` INSERT INTO sys_role_menu (role_id, menu_id) values(?, ?)`
		for i := 0; i < len(menu_ids); i++ {
			_, err = o.Raw(sql, sr.Id, menu_ids[i]).Exec()
			if err != nil {
				break
			}
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}
