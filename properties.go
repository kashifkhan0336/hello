package main

var ObserveStateProperty = MPVRequest{
	Command: []interface{}{"observe_property", 1, "pause"},
}
var ObserveVolumeProperty = MPVRequest{
	Command: []interface{}{"observe_property", 1, "volume"},
}
var ObserveTimePositionProperty = MPVRequest{
	Command: []interface{}{"observe_property", 1, "time-pos"},
}
var ObserveSeekProperty = MPVRequest{
	Command: []interface{}{"observe_property", 1, "playback-restart"},
}
var GetPosition = MPVRequest{
	Command: []interface{}{"get_property", "time-pos"},
}
var SetStatePause = MPVRequest{
	Command: []interface{}{"set_property", "pause", true},
}
var ShowStatePause = MPVRequest{
	Command: []interface{}{"show_text", "Paused!"},
}
var SetStatePlay = MPVRequest{
	Command: []interface{}{"set_property", "pause", false},
}
