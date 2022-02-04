package keygen

//go:generate  mockgen -source ../repository/keyRepository/KeyRepository.go -destination ../../mocks/keygen_mocks.go -package=mocks

import (
	"email-archiver-cli/mocks"
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

			keygen, _ := NewKeygen(mockRepository)

			_, err := keygen.CreateKey(mockKeyName)

			assert.Error(t, err)
		})

		t.Run("Should generate and persist new key when not present yet", func(t *testing.T) {
			mockRepository := mocks.NewMockKeyRepository(mockController)
			mockRepository.EXPECT().Contains(mockKeyName).Return(false)
			mockRepository.EXPECT().Persist(mockKeyName, gomock.Any())

			keygen, _ := NewKeygen(mockRepository)

			_, err := keygen.CreateKey(mockKeyName)

			assert.NoError(t, err)
		})

	})
}
