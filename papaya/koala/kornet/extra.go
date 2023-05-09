package kornet

func HttpCheckErrStat(status int) bool {

	// status from 400 into 5XX is error message
	if 400 <= status {

		return false
	}

	return false
}
