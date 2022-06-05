package notification

import "fmt"

type FmtNotifier struct {
}

func NewFmtNotifier() *FmtNotifier {
	return &FmtNotifier{}
}

func (s *FmtNotifier) Send(m string) error {
	fmt.Println(m)
	return nil
}
