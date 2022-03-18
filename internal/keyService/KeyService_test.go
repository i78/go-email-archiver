package keyService

//go:generate  mockgen -source KeyRepository.go -destination ../../mocks/KeyRepository_mocks.go -package=mocks

import (
	"email-archiver-cli/mocks"
	"email-archiver-cli/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKeyGeneration(t *testing.T) {

	t.Run("When new key requested", func(t *testing.T) {
		const mockKeyName = "MOCK_KEY_NAME"

		mockController := gomock.NewController(t)

		t.Run("Should raise error when key already exists and no rotation requested", func(t *testing.T) {
			mockRepository := mocks.NewMockKeyRepository(mockController)
			mockRepository.EXPECT().Contains(mockKeyName).Return(true)

			keygen := KeyService{
				keyRepository: mockRepository,
				keyGenerator:  mockKeyGenerator,
			}

			_, err := keygen.CreateKey(mockKeyName, false)

			assert.Error(t, err)
		})

		t.Run("Should generate and persist new key when not present yet", func(t *testing.T) {
			mockRepository := mocks.NewMockKeyRepository(mockController)
			mockRepository.EXPECT().Contains(mockKeyName).Return(false)
			mockRepository.EXPECT().Persist(mockKeyName, gomock.Any())

			keygen := KeyService{
				keyRepository: mockRepository,
				keyGenerator:  mockKeyGenerator,
			}

			_, err := keygen.CreateKey(mockKeyName, false)

			assert.NoError(t, err)
		})

	})
}

func mockKeyGenerator(int) (key *models.Key, err error) {
	return
}
