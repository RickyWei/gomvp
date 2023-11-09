// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package po

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID        int64  `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true;comment:id" json:"id"`                                      // id
	UserName  string `gorm:"column:user_name;type:varchar(64);not null;uniqueIndex:uk_user_name,priority:1;comment:user_name" json:"user_name"` // user_name
	NickName  string `gorm:"column:nick_name;type:varchar(64);not null;comment:nick_name" json:"nick_name"`                                     // nick_name
	Email     string `gorm:"column:email;type:varchar(64);not null;uniqueIndex:uk_email,priority:1;comment:email" json:"email"`                 // email
	Password  string `gorm:"column:password;type:varchar(64);not null;comment:password" json:"password"`                                        // password
	Mobile    string `gorm:"column:mobile;type:varchar(16);not null;comment:mobile" json:"mobile"`                                              // mobile
	CreatedAt int64  `gorm:"column:created_at;type:bigint(20);not null;comment:create time" json:"created_at"`                                  // create time
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint(20);not null;comment:update time" json:"updated_at"`                                  // update time
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint(20);not null;comment:delete time" json:"deleted_at"`                                  // delete time
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
