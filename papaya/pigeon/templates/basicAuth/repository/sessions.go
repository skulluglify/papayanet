package repository

import (
  "encoding/hex"
  "errors"
  "skfw/papaya/pigeon/templates/basicAuth/models"
  "time"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

var TokenExpiredOrUserNoLongerActive = errors.New("token has expired or user is no longer active")
var SessionReachedLimit = errors.New("create session has reached a limit")

type SessionRepository struct {
  DB *gorm.DB
}

type SessionRepositoryImpl interface {
  Init(db *gorm.DB) error
  Create(session *models.SessionModel, activeDuration time.Duration, maxSessions int) error
  Search(session *models.SessionModel) bool
  Delete(session *models.SessionModel) error

  SearchFast(userId uuid.UUID, token string) (*models.SessionModel, bool)
  CreateFast(userId uuid.UUID, clientIP string, userAgent string, token string, secret string, expired time.Time, activeDuration time.Duration, maxSessions int) (*models.SessionModel, error)
  DeleteFast(userId uuid.UUID, token string) error

  SearchFastById(sessionId uuid.UUID) (*models.SessionModel, bool)

  RecoveryFast(userId uuid.UUID, token string, activeDuration time.Duration, maxSessions int) error
  CreateFastAutoToken(user *models.UserModel, clientIP string, userAgent string, expired time.Time, activeDuration time.Duration, maxSessions int) (*models.SessionModel, error)

  CheckIn(session *models.SessionModel) error

  GetAll(userId uuid.UUID, maxSessions int) ([]models.SessionModel, error)

  SessionNew()
}

func SessionRepositoryNew(db *gorm.DB) (SessionRepositoryImpl, error) {

  sessionRepo := &SessionRepository{}
  err := sessionRepo.Init(db.Model(&models.SessionModel{}))
  if err != nil {
    return nil, err
  }
  return sessionRepo, nil
}

func (s *SessionRepository) Init(db *gorm.DB) error {

  if db == nil {

    return errors.New("gorm.DB is NULL")
  }

  s.DB = db
  return nil
}

func (s *SessionRepository) Search(session *models.SessionModel) bool {

  s.SessionNew()

  if session == nil {

    return false
  }

  if EmptyIds(session.UserID) && session.Token == "" {

    return false
  }

  var sessions []models.SessionModel
  sessions = make([]models.SessionModel, 0)

  if s.DB.Where(session).Limit(1).Find(&sessions).Error != nil {

    return false
  }

  if len(sessions) > 0 {

    session = &sessions[0]
    return true
  }

  return false
}

func (s *SessionRepository) Create(session *models.SessionModel, activeDuration time.Duration, maxSessions int) error {

  var err error

  if session == nil {

    return errors.New("session is NULL")
  }

  if EmptyIds(session.UserID) || session.Token == "" || session.SecretKey == "" {

    return errors.New("userId, token, and secret key can't be empty")
  }

  var idx uuid.UUID

  idx = Ids(session.UserID)

  if EmptyIdx(idx) {

    return errors.New("ids doesn't match id format")
  }

  if err := s.RecoveryFast(idx, session.Token, activeDuration, maxSessions); err != nil {

    return err
  }

  var id []byte

  id, err = uuid.New().MarshalBinary()

  if err != nil {

    return err
  }

  session.ID = hex.EncodeToString(id)

  if err != nil {

    return err
  }

  if s.DB.Create(session).Error != nil {

    return errors.New("unable to create session")
  }

  return nil
}

func (s *SessionRepository) Delete(session *models.SessionModel) error {

  if session == nil {

    return errors.New("session is NULL")
  }

  if EmptyIds(session.UserID) && session.Token == "" {

    return errors.New("userId, or token can't be empty")
  }

  if !s.Search(session) {

    return errors.New("session has been deleted")
  }

  if s.DB.Delete(session).Error != nil {

    return errors.New("unable to delete session")
  }

  return nil
}

func (s *SessionRepository) SearchFast(userId uuid.UUID, token string) (*models.SessionModel, bool) {

  s.SessionNew()

  token = EmptyAsterisk(token)

  var sessions []models.SessionModel
  sessions = make([]models.SessionModel, 0)

  if !EmptyIdx(userId) {

    if token != "" {

      if s.DB.Where("user_id = ? OR token = ?", Idx(userId), token).Limit(1).Find(&sessions).Error != nil {

        return nil, false
      }

    } else {

      if s.DB.Where("user_id = ?", Idx(userId)).Limit(1).Find(&sessions).Error != nil {

        return nil, false
      }
    }

  } else {

    if token != "" {

      if s.DB.Where("token = ?", token).Limit(1).Find(&sessions).Error != nil {

        return nil, false
      }
    } else {

      // userId, or token is empty
      return nil, false
    }
  }

  if len(sessions) > 0 {

    return &sessions[0], true
  }

  return nil, false
}

func (s *SessionRepository) CreateFast(userId uuid.UUID, clientIP string, userAgent string, token string, secret string, expired time.Time, activeDuration time.Duration, maxSessions int) (*models.SessionModel, error) {

  if EmptyIdx(userId) && token == "" {

    return nil, errors.New("userId, or token can't be empty")
  }

  var sID, uID string

  sID = Idx(uuid.New())

  if EmptyIds(sID) {

    return nil, errors.New("unable convert id to string")
  }

  uID = Idx(userId)

  if EmptyIds(uID) {

    return nil, errors.New("unable convert id to string")
  }

  var session models.SessionModel

  session = models.SessionModel{
    ID:        sID,
    UserID:    uID,
    ClientIP:  clientIP,
    UserAgent: userAgent,
    Token:     token,
    SecretKey: secret,
    Expired:   expired,
  }

  if err := s.RecoveryFast(userId, token, activeDuration, maxSessions); err != nil {

    return nil, err
  }

  if s.DB.Create(&session).Error != nil {

    return nil, errors.New("unable to create session")
  }

  return &session, nil
}

func (s *SessionRepository) DeleteFast(userId uuid.UUID, token string) error {

  if EmptyIdx(userId) && token == "" {

    return errors.New("userId, or token can't be empty")
  }

  if _, found := s.SearchFast(userId, token); !found {

    return errors.New("session has been deleted")
  }

  if s.DB.Where("user_id = ? OR token = ?", Idx(userId), token).Delete(&models.SessionModel{}).Error != nil {

    return errors.New("unable to delete session")
  }

  return nil
}

func (s *SessionRepository) GetAll(userId uuid.UUID, maxSessions int) ([]models.SessionModel, error) {

  s.SessionNew()

  if EmptyIdx(userId) {

    return nil, errors.New("userId is empty")
  }

  var sessions []models.SessionModel
  sessions = make([]models.SessionModel, 0)

  if s.DB.Where("user_id = ?", Idx(userId)).Limit(maxSessions).Find(&sessions).Error != nil {

    return nil, errors.New("unable to search all sessions")
  }

  return sessions, nil
}

func (s *SessionRepository) SearchFastById(sessionId uuid.UUID) (*models.SessionModel, bool) {

  s.SessionNew()

  if EmptyIdx(sessionId) {

    return nil, false
  }

  var sessions []models.SessionModel

  if s.DB.Where("id = ?", Idx(sessionId)).Limit(1).Find(&sessions).Error != nil {

    return nil, false
  }

  if len(sessions) > 0 {

    return &sessions[0], true
  }

  return nil, false
}

func (s *SessionRepository) RecoveryFast(userId uuid.UUID, token string, activeDuration time.Duration, maxSessions int) error {

  s.SessionNew()

  if EmptyIdx(userId) {

    return errors.New("userId is empty")
  }

  var k int
  var ok bool
  var sessions []models.SessionModel

  k = 0
  sessions = make([]models.SessionModel, 0)
  currentTime := time.Now().UTC()

  // limit process to get max sessions requirement
  if s.DB.Where("user_id = ?", Idx(userId)).Limit(maxSessions).Find(&sessions).Error != nil {

    return errors.New("unable to find session")
  }

  for _, session := range sessions {

    ok = currentTime.Before(session.Expired) &&
      currentTime.Before(session.LastActivated.Add(activeDuration))

    if !ok {

      if err := s.Delete(&session); err != nil {

        return err
      }

      if token == session.Token {

        return TokenExpiredOrUserNoLongerActive
      }

      continue
    }

    k++
  }

  if k == maxSessions {

    return SessionReachedLimit
  }

  return nil
}

func (s *SessionRepository) CreateFastAutoToken(user *models.UserModel, clientIP string, userAgent string, expired time.Time, activeDuration time.Duration, maxSessions int) (*models.SessionModel, error) {

  if user == nil {

    return nil, errors.New("user is NULL")
  }

  if EmptyIds(user.ID) || user.Username == "" || user.Email == "" {

    return nil, errors.New("userId, username, and email can't be empty")
  }

  currentTime := time.Now().UTC()

  data := map[string]any{

    "username": user.Username,
    "email":    user.Email,
    "admin":    user.Admin,
    "iat":      currentTime.Unix(),
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

  var idx uuid.UUID

  idx = Ids(user.ID)

  if EmptyIdx(idx) {

    return nil, errors.New("ids doesn't match id format")
  }

  return s.CreateFast(idx, clientIP, userAgent, token, secret, expired, activeDuration, maxSessions)
}

func (s *SessionRepository) CheckIn(session *models.SessionModel) error {

  currentTime := time.Now().UTC()

  if s.DB.Where(session).Update("last_activated", currentTime).Error != nil {

    return errors.New("unable to perform check-in session")
  }

  return nil
}

func (s *SessionRepository) SessionNew() {

  s.DB = s.DB.Session(&gorm.Session{})
}
