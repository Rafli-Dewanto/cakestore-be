package repository

import (
	"cakestore/internal/entity"
	"cakestore/internal/repository"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CakeRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository repository.CakeRepository
}

func (suite *CakeRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&entity.Cake{})
	assert.NoError(suite.T(), err)

	logger := logrus.New()
	suite.db = db
	suite.repository = repository.NewCakeRepository(db, logger)
}

func (suite *CakeRepositoryTestSuite) TearDownTest() {
	// Clean up after each test
	sqlDB, err := suite.db.DB()
	assert.NoError(suite.T(), err)
	err = sqlDB.Close()
	assert.NoError(suite.T(), err)
}

func TestCakeRepositorySuite(t *testing.T) {
	suite.Run(t, new(CakeRepositoryTestSuite))
}

func (suite *CakeRepositoryTestSuite) TestCreate() {
	cake := &entity.Cake{
		Title:       "Test Cake",
		Description: "Test Description",
		Rating:      4.5,
		Image:       "test.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := suite.repository.Create(cake)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), cake.ID)
}

func (suite *CakeRepositoryTestSuite) TestGetByID() {
	// Create a test cake
	cake := &entity.Cake{
		Title:       "Test Cake",
		Description: "Test Description",
		Rating:      4.5,
		Image:       "test.jpg",
	}
	err := suite.repository.Create(cake)
	assert.NoError(suite.T(), err)

	// Test getting the cake
	foundCake, err := suite.repository.GetByID(int(cake.ID))
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundCake)
	assert.Equal(suite.T(), cake.Title, foundCake.Title)

	// Test getting non-existent cake
	_, err = suite.repository.GetByID(9999)
	assert.Error(suite.T(), err)
}

func (suite *CakeRepositoryTestSuite) TestGetAll() {
	cakes := []entity.Cake{
		{Title: "Cake A", Rating: 4.5},
		{Title: "Cake B", Rating: 5.0},
		{Title: "Cake C", Rating: 4.0},
	}

	for _, cake := range cakes {
		err := suite.repository.Create(&cake)
		assert.NoError(suite.T(), err)
	}

	foundCakes, err := suite.repository.GetAll()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(cakes), len(foundCakes))


	assert.Greater(suite.T(), foundCakes[0].Rating, foundCakes[1].Rating)
	assert.Greater(suite.T(), foundCakes[1].Rating, foundCakes[2].Rating)
}

func (suite *CakeRepositoryTestSuite) TestUpdateCake() {
	cake := &entity.Cake{
		Title:       "Original Cake",
		Description: "Original Description",
		Rating:      4.0,
		Image:       "original.jpg",
	}
	err := suite.repository.Create(cake)
	assert.NoError(suite.T(), err)

	cake.Title = "Updated Cake"
	cake.Rating = 4.5
	err = suite.repository.UpdateCake(cake)
	assert.NoError(suite.T(), err)

	updatedCake, err := suite.repository.GetByID(int(cake.ID))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Cake", updatedCake.Title)
	assert.Equal(suite.T(), 4.5, updatedCake.Rating)

	nonExistentCake := &entity.Cake{ID: 9999}
	err = suite.repository.UpdateCake(nonExistentCake)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "no rows updated, cake not found", err.Error())
}

func (suite *CakeRepositoryTestSuite) TestDelete() {
	cake := &entity.Cake{
		Title:       "Test Cake",
		Description: "Test Description",
	}
	err := suite.repository.Create(cake)
	assert.NoError(suite.T(), err)

	err = suite.repository.Delete(int(cake.ID))
	assert.NoError(suite.T(), err)

	_, err = suite.repository.GetByID(int(cake.ID))
	assert.Error(suite.T(), err)
}
