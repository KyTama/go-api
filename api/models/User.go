package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"strings"
	"time"
)

type User struct {
	ID            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username      string    `gorm:"size:255;not null;unique" json:"username"`
	Fullname      string    `gorm:"size:255;not null;" json:"fullname"`
	Email         string    `gorm:"size:100;not null;unique" json:"email"`
	Password      string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	LastUpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_updated_at"`
	Active        bool      `gorm:"default:true" json"active"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Fullname = html.EscapeString(strings.TrimSpace(u.Fullname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.LastUpdatedAt = time.Now()
	u.Active = true
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Fullname == "" {
			return errors.New("Required Fullname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Username == "" {
			return errors.New("Required Username")
		}
		if u.Fullname == "" {
			return errors.New("Required Fullname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) GetUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Where("active = true").Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) GetUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ? and active = true", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateUserByID(db *gorm.DB, uid uint32) (*User, error) {

	//To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":        u.Password,
			"username":        u.Username,
			"fullname":        u.Fullname,
			"email":           u.Email,
			"last_updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// To show the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteUserByID(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ? and active = true", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"active":          false,
			"last_updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) SignIn(db *gorm.DB, email string) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("email = ? and active = true", email).Take(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, err
}
