package repository

import (
  "errors"
  "github.com/google/uuid"
  "gorm.io/gorm"
  "skfw/papaya/pigeon/templates/basic/models"
  "time"
)

type SessionRepository struct {
  DB *gorm.DB
}

type SessionRepositoryImpl interface {
  Init(db *gorm.DB) error
  Create(session *models.SessionModel) error
  Search(session *models.SessionModel) bool
  Delete(session *models.SessionModel) error
  SearchFast(userid uuid.UUID, token string) (*models.SessionModel, bool)
  CreateFast(userid uuid.UUID, token string, secret string, expired time.Time) (*models.SessionModel, error)
  DeleteFast(userid uuid.UUID, token string) error

  CreateFastAutoToken(user *models.UserModel, expired time.Time) (*models.SessionModel, error)
}

func SessionRepositoryNew(db *gorm.DB) SessionRepositoryImpl {

  sessionRepo := &SessionRepository{}
  err := sessionRepo.Init(db.Model(&models.SessionModel{}))
  if err != nil {
    return nil
  }
  return sessionRepo
}

func (s *SessionRepository) Init(db *gorm.DB) error {

  if db == nil {

    return errors.New("gorm.DB is NULL")
  }

  s.DB = db
  return nil
}

func (s *SessionRepository) Search(session *models.SessionModel) bool {

  if session == nil {

    return false
  }

  if session.UserID.String() == "" && session.Token == "" {

    return false
  }

  var sessions []models.SessionModel
  if s.DB.Where(session).Limit(1).Find(sessions).Error != nil {

    return false
  }

  if len(sessions) > 0 {

    session = &sessions[0]
    return true
  }

  return false
}

func (s *SessionRepository) Create(session *models.SessionModel) error {

  if session == nil {

    return errors.New("session is NULL")
  }

  if session.UserID.String() == "" || session.Token == "" || session.SecretKey == "" {

    return errors.New("userid, token, and secret key can't be empty")
  }

  if s.Search(session) {

    return errors.New("session has been added")
  }

  if s.DB.Create(session).Error != nil {

    return errors.New("unable to add session")
  }

  return nil
}

func (s *SessionRepository) Delete(session *models.SessionModel) error {

  if session == nil {

    return errors.New("session is NULL")
  }

  if session.UserID.String() == "" && session.Token == "" {

    return errors.New("userid, or token can't be empty")
  }

  if !s.Search(session) {

    return errors.New("session has been deleted")
  }

  // force delete with unscoped
  if s.DB.Unscoped().Delete(session).Error != nil {

    return errors.New("unable to delete session")
  }

  return nil
}

func (s *SessionRepository) SearchFast(userid uuid.UUID, token string) (*models.SessionModel, bool) {

  if userid.String() == "" && token == "" {

    return nil, false
  }

  var sessions []models.SessionModel

  if s.DB.Where("user_id = ? OR token = ?", userid, token).Limit(1).Find(&sessions).Error != nil {

    return nil, false
  }

  if len(sessions) > 0 {

    return &sessions[0], true
  }

  return nil, false
}

func (s *SessionRepository) CreateFast(userid uuid.UUID, token string, secret string, expired time.Time) (*models.SessionModel, error) {

  if userid.String() == "" && token == "" {

    return nil, errors.New("userid, or token can't be empty")
  }

  var session models.SessionModel

  session = models.SessionModel{
    UserID:    userid,
    Token:     token,
    SecretKey: secret,
    Expired:   expired,
  }

  if _, found := s.SearchFast(userid, token); found {

    return nil, errors.New("session has been added")
  }

  if s.DB.Create(&session).Error != nil {

    return nil, errors.New("unable to add session")
  }

  return &session, nil
}

func (s *SessionRepository) DeleteFast(userid uuid.UUID, token string) error {

  if userid.String() == "" && token == "" {

    return errors.New("userid, or token can't be empty")
  }

  if _, found := s.SearchFast(userid, token); !found {

    return errors.New("session has been deleted")
  }

  var session models.SessionModel

  session = models.SessionModel{
    UserID: userid,
    Token:  token,
  }

  // force delete with unscoped
  if s.DB.Unscoped().Delete(session).Error != nil {

    return errors.New("unable to delete session")
  }

  return nil
}

func (s *SessionRepository) CreateFastAutoToken(user *models.UserModel, expired time.Time) (*models.SessionModel, error) {

  if user == nil {

    return nil, errors.New("user is NULL")
  }

  if user.ID.String() == "" || user.Username == "" || user.Email == "" {

    return nil, errors.New("userid, username, and email can't be empty")
  }

  data := map[string]any{

    "username": user.Username,
    "email":    user.Email,
    "admin":    user.Admin,
    "iat":      time.Now().UTC().Unix(),
    "exp":      expired.Unix(),
  }

  secret, err := CreateSecretKey()
  if err != nil {

    return nil, errors.New("unable to generate secret key")
  }

  token, err := EncodeJWT(data, secret)
  if err != nil {

    return nil, errors.New("unable to create token")
  }

  return s.CreateFast(user.ID, token, secret, expired)
}
