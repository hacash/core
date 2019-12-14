package interfaces

type HacashNode interface {
	Launch(inicnffilepath string) error

	StartMining()
	StopMining()
}
