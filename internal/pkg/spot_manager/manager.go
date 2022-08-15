package spotmanager

type ISpotManager interface {
	Run(startC chan<- error)
}
