package usecase

import (
	"cakestore/internal/entity"
	"cakestore/internal/repository"

	"github.com/sirupsen/logrus"
)

type CakeUseCase interface {
	GetAllCakes() ([]entity.Cake, error)
	GetCakeByID(id int) (*entity.Cake, error)
	CreateCake(cake *entity.Cake) error
	UpdateCake(cake *entity.Cake) error
	DeleteCake(id int) error
}

type cakeUseCase struct {
	repo   repository.CakeRepository
	logger *logrus.Logger
}

func NewCakeUseCase(repo repository.CakeRepository, logger *logrus.Logger) CakeUseCase {
	return &cakeUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *cakeUseCase) GetAllCakes() ([]entity.Cake, error) {
	cakes, err := uc.repo.GetAll()
	if err != nil {
		uc.logger.Errorf("Error fetching all cakes: %v", err)
		return nil, err
	}
	uc.logger.Info("Successfully fetched all cakes")
	return cakes, nil
}

func (uc *cakeUseCase) GetCakeByID(id int) (*entity.Cake, error) {
	cake, err := uc.repo.GetByID(id)
	if err != nil {
		uc.logger.Errorf("Error fetching cake with ID %d: %v", id, err)
		return nil, err
	}
	uc.logger.Infof("Successfully fetched cake with ID %d", id)
	return cake, nil
}

func (uc *cakeUseCase) CreateCake(cake *entity.Cake) error {
	err := uc.repo.Create(cake)
	if err != nil {
		uc.logger.Errorf("Error creating cake: %v", err)
		return err
	}
	uc.logger.Infof("Successfully created a new cake: %s", cake.Title)
	return nil
}

func (uc *cakeUseCase) UpdateCake(cake *entity.Cake) error {
	err := uc.repo.UpdateCake(cake)
	if err != nil {
		uc.logger.Errorf("Error updating cake: %v", err)
		return err
	}
	uc.logger.Infof("Successfully updated cake with ID %d", cake.ID)
	return nil
}

func (uc *cakeUseCase) DeleteCake(id int) error {
	err := uc.repo.Delete(id)
	if err != nil {
		uc.logger.Errorf("Error deleting cake with ID %d: %v", id, err)
		return err
	}
	uc.logger.Infof("Successfully deleted cake with ID %d", id)
	return nil
}
