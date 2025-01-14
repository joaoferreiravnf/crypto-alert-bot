package services

import (
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	mock_prompt "uphold-alert-bot/internal/adapters/mocks/mock_api"
	"uphold-alert-bot/internal/mocks/mock_scheduler"

	"github.com/stretchr/testify/assert"

	"uphold-alert-bot/internal/models"
)

func TestTickerScheduler(t *testing.T) {
	t.Run("Below threshold", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := mock_prompt.NewMockApiDataValidator(ctrl)
		mockRepo := mock_services.NewMockRecorder(ctrl)
		mockPublisher := mock_services.NewMockPublisher(ctrl)

		testTicker := &models.Ticker{
			Config: models.TickerConfig{
				RefreshRate:     1,
				PercOscillation: 5.0,
			},
		}
		testTicker.PreviousAsk = 100.0
		testTicker.CurrentAsk = 102.0

		mockAPI.EXPECT().
			FetchPairData(testTicker).
			Return(nil).
			AnyTimes()

		sched := NewTickerScheduler(mockAPI, testTicker, mockRepo, mockPublisher)

		sched.SchedulerStart()

		time.Sleep(2 * time.Second)

		sched.SchedulerStop()

		assert.Equal(t, float64(102.0), testTicker.CurrentAsk.Float64())
	})

	t.Run("Above threshold", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := mock_api.NewMockApiResponseInterface(ctrl)
		mockRepo := mock_repository.NewMockRepository(ctrl)
		mockPublisher := mock_services.NewMockPublisherInterface(ctrl)

		testTicker := &models.Ticker{
			Config: models.TickerConfig{
				RefreshRate:     1,
				PercOscillation: 5.0,
			},
			PreviousAsk: 100.0,
			CurrentAsk:  110.0,
		}

		mockAPI.EXPECT().
			FetchPairData(testTicker).
			Return(nil).
			AnyTimes()

		mockPublisher.EXPECT().
			Publish(gomock.Any(), testTicker).
			Times(1)

		mockRepo.EXPECT().
			Save(gomock.Any(), testTicker).
			Return(nil).
			Times(1)

		sched := NewTickerScheduler(mockAPI, testTicker, mockRepo, mockPublisher)
		sched.SchedulerStart()

		time.Sleep(2 * time.Second)

		sched.SchedulerStop()
	})
}
