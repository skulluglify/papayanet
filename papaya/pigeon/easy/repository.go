package easy

import (
  "errors"
  "fmt"
  "github.com/google/uuid"
  "gorm.io/gorm"
  "skfw/papaya/koala/mapping"
)

type Repository[T any] struct {
  DB   *gorm.DB
  Name string
}

type RepositoryImpl[T any] interface {
  Init(DB *gorm.DB, model *T) error
  SessionNew()
  Find(query any, args ...any) (*T, error)
  FindAll(size int, page int, query any, args ...any) ([]T, error)
  CatchAll(size int, page int) ([]T, error)
  Create(model *T) (*T, error)
  Update(model *T, query any, args ...any) error
  Remove(query any, args ...any) error
  Delete(query any, args ...any) error
  Unscoped() RepositoryImpl[T]
  GORM() *gorm.DB
}

func RepositoryNew[T any](DB *gorm.DB, model *T) (RepositoryImpl[T], error) {

  modelRepo := &Repository[T]{}
  return modelRepo, modelRepo.Init(DB, model)
}

func (u *Repository[T]) Init(DB *gorm.DB, model *T) error {

  if DB == nil {

    return errors.New("DB is NULL")
  }
  u.DB = DB.Model(model)
  u.Name = TableName(model)
  return nil
}

func (u *Repository[T]) SessionNew() {

  // reset session with new session
  u.DB = u.DB.Session(&gorm.Session{})
}

func (u *Repository[T]) Find(query any, args ...any) (*T, error) {

  var err error
  var models []T

  if models, err = u.FindAll(1, 1, query, args); err != nil {

    return nil, err
  }

  if len(models) > 0 {

    return &models[0], nil
  }

  return nil, errors.New(fmt.Sprintf("%s not found", u.Name))
}

func (u *Repository[T]) FindAll(size int, page int, query any, args ...any) ([]T, error) {

  u.SessionNew()

  var err error

  models := make([]T, 0)

  if page > 0 {

    offset := size * (page - 1)
    limit := size

    if err = u.DB.Where(query, args).Offset(offset).Limit(limit).Find(&models).Error; err != nil {

      return models, errors.New(fmt.Sprintf("unable to catch %ss", u.Name))
    }
  }

  return models, nil
}

func (u *Repository[T]) CatchAll(size int, page int) ([]T, error) {

  u.SessionNew()

  var err error

  models := make([]T, 0)

  if page > 0 {

    offset := size * (page - 1)
    limit := size

    if err = u.DB.Offset(offset).Limit(limit).Find(&models).Error; err != nil {

      return models, errors.New(fmt.Sprintf("unable to catch %ss", u.Name))
    }
  }

  return models, nil
}

func (u *Repository[T]) Create(model *T) (*T, error) {

  var err error

  if model != nil {

    if _, err = u.Find(model); err != nil {

      // bind random ID
      if err = StructSet(model, "ID", Idx(uuid.New())); err != nil {

        return nil, errors.New("invalid model")
      }

      if err = u.DB.Create(model).Error; err != nil {

        // fallback data model
        return model, errors.New(fmt.Sprintf("unable to create new %s", u.Name))
      }

      return model, nil
    }

    return model, errors.New("user has been added")
  }

  return model, errors.New(fmt.Sprintf("%s is NULL", u.Name))
}

func (u *Repository[T]) Update(model *T, query any, args ...any) error {

  var err error
  var info *T

  if model != nil {

    if info, err = u.Find(query, args...); err != nil {

      return err
    }

    // bind ID
    var xID any
    var infoId string

    if xID, err = StructGet(info, "ID"); err != nil {

      return errors.New(fmt.Sprintf("unable to get ID from %s information", u.Name))
    }

    infoId = mapping.KValueToString(xID)

    if infoId != "" {

      // bind random ID
      if err = StructSet(model, "ID", Idx(uuid.New())); err != nil {

        return errors.New(fmt.Sprintf("invalid %s", u.Name))
      }

    } else {

      return errors.New(fmt.Sprintf("invalid ID from %s information", u.Name))
    }

    prepared := u.DB.Where(query, args...).Updates(&model)

    if err = prepared.Error; err != nil {

      return errors.New(fmt.Sprintf("unable to update %s information", u.Name))
    }

    if prepared.RowsAffected > 0 {

      return nil
    }

    return errors.New(fmt.Sprintf("not match any data from %s ", u.Name))
  }

  return errors.New(fmt.Sprintf("%s is NULL", u.Name))
}

func (u *Repository[T]) Remove(query any, args ...any) error {

  var err error
  var info *T

  if info, err = u.Find(query, args...); err != nil {

    return err
  }

  prepared := u.DB.Where(query, args...).Delete(&info)

  if err = prepared.Error; err != nil {

    return errors.New(fmt.Sprintf("unable to remove %s", u.Name))
  }

  if prepared.RowsAffected > 0 {

    return nil
  }

  return errors.New(fmt.Sprintf("not match any data from %s ", u.Name))
}

func (u *Repository[T]) Delete(query any, args ...any) error {

  var err error
  var info *T

  if info, err = u.Find(query, args...); err != nil {

    return err
  }

  prepared := u.DB.Unscoped().Where(query, args...).Delete(&info)

  if err = prepared.Error; err != nil {

    return errors.New(fmt.Sprintf("unable to delete %s", u.Name))
  }

  if prepared.RowsAffected > 0 {

    return nil
  }

  return errors.New(fmt.Sprintf("not match any data from %s ", u.Name))
}

func (u *Repository[T]) Unscoped() RepositoryImpl[T] {

  return &Repository[T]{
    DB:   u.DB.Unscoped(),
    Name: u.Name,
  }
}

func (u *Repository[T]) GORM() *gorm.DB {

  return u.DB
}
