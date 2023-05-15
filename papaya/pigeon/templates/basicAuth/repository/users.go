package repository

import (
  "encoding/hex"
  "errors"
  "skfw/papaya/pigeon/templates/basicAuth/models"
  "skfw/papaya/pigeon/templates/basicAuth/util"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type UserRepository struct {
  DB *gorm.DB
}

type UserRepositoryImpl interface {
  Init(db *gorm.DB) error
  Create(user *models.UserModel) error
  Search(user *models.UserModel) bool
  Delete(user *models.UserModel) error
  SearchFast(username string, email string) (*models.UserModel, bool)
  CreateFast(name string, username string, email string, password string) (*models.UserModel, error)
  DeleteFast(username string, email string) error

  SearchFastById(userId uuid.UUID) (*models.UserModel, bool)

  SessionNew()
}

func UserRepositoryNew(db *gorm.DB) (UserRepositoryImpl, error) {

  userRepo := &UserRepository{}
  err := userRepo.Init(db.Model(&models.UserModel{}))
  if err != nil {
    return nil, err
  }
  return userRepo, nil
}

func (u *UserRepository) Init(db *gorm.DB) error {

  if db == nil {

    return errors.New("gorm.DB is NULL")
  }

  u.DB = db
  return nil
}

func (u *UserRepository) Search(user *models.UserModel) bool {

  u.SessionNew()

  if user == nil {

    return false
  }

  if user.Username == "" && user.Email == "" {

    return false
  }

  var users []models.UserModel
  if u.DB.Where(user).Limit(1).Find(&users).Error != nil {

    return false
  }

  if len(users) > 0 {

    user = &users[0]
    return true
  }

  return false
}

func (u *UserRepository) Create(user *models.UserModel) error {

  var err error

  if user == nil {

    return errors.New("user is NULL")
  }

  if user.Username == "" || user.Email == "" || user.Password == "" {

    return errors.New("username, email, and password can't be empty")
  }

  if u.Search(user) {

    return errors.New("user has been added")
  }

  var id []byte

  id, err = uuid.New().MarshalBinary()

  if err != nil {

    return err
  }

  user.ID = hex.EncodeToString(id)

  if u.DB.Create(user).Error != nil {

    return errors.New("unable to create user")
  }

  return nil
}

func (u *UserRepository) Delete(user *models.UserModel) error {

  if user == nil {

    return errors.New("user is NULL")
  }

  if user.Username == "" && user.Email == "" {

    return errors.New("username, or email can't be empty")
  }

  if !u.Search(user) {

    return errors.New("user has been deleted")
  }

  if u.DB.Delete(user).Error != nil {

    return errors.New("unable to delete user")
  }

  return nil
}

func (u *UserRepository) SearchFast(username string, email string) (*models.UserModel, bool) {

  u.SessionNew()

  // solve about ignoring attempts to get all data
  username = util.EmptyAsterisk(username)
  email = util.EmptyAsterisk(email)

  var users []models.UserModel
  users = make([]models.UserModel, 0)

  if username != "" {

    if email != "" {

      // have both, search by username and email, vulnerability
      if u.DB.Where("username = ? AND email = ?", username, email).Limit(1).Find(&users).Error != nil {

        return nil, false
      }

    } else {

      if u.DB.Where("username = ?", username).Limit(1).Find(&users).Error != nil {

        return nil, false
      }
    }

  } else {

    if email != "" {

      if u.DB.Where("email = ?", email).Limit(1).Find(&users).Error != nil {

        return nil, false
      }

    } else {

      // username, or email is empty
      return nil, false
    }
  }

  if len(users) > 0 {

    return &users[0], true
  }

  return nil, false
}

func (u *UserRepository) CreateFast(name string, username string, email string, password string) (*models.UserModel, error) {

  var err error

  if username == "" || email == "" || password == "" {

    return nil, errors.New("username, email, and password can't be empty")
  }

  password, err = util.HashPassword(password)
  if err != nil {

    return nil, errors.New("password can't be hashed")
  }

  var user models.UserModel

  sID := util.Idx(uuid.New())

  if err != nil {

    return nil, err
  }

  user = models.UserModel{
    Name:     name,
    Username: username,
    Email:    email,
    Password: password,
  }

  user.Model.ID = sID

  if _, found := u.SearchFast(username, email); found {

    return nil, errors.New("user has been added")
  }

  if u.DB.Create(&user).Error != nil {

    return nil, errors.New("unable to create user")
  }

  return &user, nil
}

func (u *UserRepository) DeleteFast(username string, email string) error {

  if username == "" && email == "" {

    return errors.New("username, or email can't be empty")
  }

  if _, found := u.SearchFast(username, email); !found {

    return errors.New("user has been deleted")
  }

  var user models.UserModel

  user = models.UserModel{
    Username: username,
    Email:    email,
  }

  if u.DB.Delete(user).Error != nil {

    return errors.New("unable to delete user")
  }

  return nil
}

func (u *UserRepository) SearchFastById(userId uuid.UUID) (*models.UserModel, bool) {

  u.SessionNew()

  if util.EmptyIdx(userId) {

    return nil, false
  }

  var users []models.UserModel

  if u.DB.Where("id = ?", util.Idx(userId)).Limit(1).Find(&users).Error != nil {

    return nil, false
  }

  if len(users) > 0 {

    return &users[0], true
  }

  return nil, false
}

func (u *UserRepository) SessionNew() {

  u.DB = u.DB.Session(&gorm.Session{})
}
