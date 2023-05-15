package tui

func (*stateError) Intermediate() bool {
	return true
}

func (*stateLoading) Intermediate() bool {
	return true
}

func (*stateLogin) Intermediate() bool {
	return false
}

func (*stateCourseSelection) Intermediate() bool {
	return false
}
