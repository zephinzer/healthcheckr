package scheduler

import (
	"fmt"
	"healthcheckr/internal/alert"
	"healthcheckr/internal/check"
	"healthcheckr/internal/utils"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type ScheduleOpts struct {
	Interval             time.Duration
	FailureThreshold     int
	AlertMinimumInterval time.Duration
}

type ChannelOpts struct {
	Channels []alert.Channel
}

func ScheduleHttp(
	httpCheckOpts check.HttpBasedUrlOpts,
	scheduleOpts ScheduleOpts,
	channelOpts ChannelOpts,
) error {
	urlInstance, err := utils.MakeHttpUrl(httpCheckOpts.Scheme, httpCheckOpts.Hostname, httpCheckOpts.Path, httpCheckOpts.Queries)
	if err != nil {
		return fmt.Errorf("failed to make url: %s", err)
	}
	logrus.Infof("scheduling url check for url[%s]...", urlInstance.String())
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	go func(
		channels []alert.Channel,
	) {
		isInAlertingMode := false
		lastSentNotificationAt := time.Time{}
		checkFailures := 0
		for {
			if isTerminated {
				logrus.Debugf("isTerminated flag was set to truthy, terminating...")
				waiter.Done()
				return
			}
			if err := check.HttpBasedUrl(httpCheckOpts); err != nil {
				checkFailures++
				logrus.Warnf("failed http-based url check: %s", err)
				if checkFailures > scheduleOpts.FailureThreshold {
					if time.Since(lastSentNotificationAt) > scheduleOpts.AlertMinimumInterval {
						if len(channels) > 0 {
							logrus.Infof("triggering failure notification to %v external channels...", len(channels))
							for _, channel := range channels {
								logrus.Debugf("triggering failure notification to channel of type[%s]...", channel.GetType())
								if err := channel.Send(
									alert.Message{
										Title:   fmt.Sprintf("`%s` HTTP check failed", urlInstance.Hostname()),
										Body:    fmt.Sprintf("HTTP check failed for full URL `%s`", urlInstance.String()),
										Details: err.Error(),
										Type:    alert.TypeError,
									},
								); err != nil {
									logrus.Warnf("failed to send notification: %s", err)
								}
							}
						}
						isInAlertingMode = true
						lastSentNotificationAt = time.Now()
					} else {
						logrus.Debugf("skipped triggering notification, minimum alert interval %v has not elapsed", scheduleOpts.AlertMinimumInterval)
					}
				}
			} else {
				checkFailures = 0
				if isInAlertingMode {
					if len(channels) > 0 {
						logrus.Infof("triggering success notification to %v external channels...", len(channels))
						for _, channel := range channels {
							logrus.Debugf("triggering success notification to channel of type[%s]...", channel.GetType())
							if err := channel.Send(
								alert.Message{
									Title: fmt.Sprintf("`%s` HTTP check succeeded", urlInstance.Hostname()),
									Body:  fmt.Sprintf("HTTP check succeeded for full URL `%s`", urlInstance.String()),
									Type:  alert.TypeSuccess,
								},
							); err != nil {
								logrus.Warnf("failed to send notification: %s", err)
							}
						}
					}
					lastSentNotificationAt = time.Now()
					isInAlertingMode = false
				} else {
					logrus.Debugf("skipped triggering notification, success condition reached within failure threshold")
				}
			}
			logrus.Debugf("next run in duration[%v]...", scheduleOpts.Interval)
			<-time.After(scheduleOpts.Interval)
		}
	}(
		channelOpts.Channels,
	)
	waiter.Wait()
	logrus.Infof("check for url[%s] terminated gracefully", urlInstance.String())
	return nil
}
