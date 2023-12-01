package model

type State map[int64]ActiveFlags

type ActiveFlags struct {
	UnsubscribeActiveFlag bool
	UpdateTimeActiveFlag  bool
	HolidayActiveFlag     bool
	WeatherActiveFlag     bool
	RegularActiveFlag     bool
	SubscribeActiveFlag   bool
}
