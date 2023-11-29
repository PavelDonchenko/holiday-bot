package model

type State map[int64]ActiveFlags

type ActiveFlags struct {
	subscribeActiveFlag   bool
	UnsubscribeActiveFlag bool
}
