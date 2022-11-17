package spotmanager

import "context"

type ISpotManager interface {
	Run(ctx context.Context, startC chan<- error)

	CheckHealth() map[string]string
}
