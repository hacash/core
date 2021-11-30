package interfacev2

type HacashNode interface {
	Launch(inicnffilepath string) error

	StartMining()
	StopMining()
}
