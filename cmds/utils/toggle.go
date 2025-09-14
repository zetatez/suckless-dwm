package utils

func Toggle(proc string) {
	if IsRunning(proc) {
		Kill(proc)
	} else {
		RunScript("bash", proc)
	}
}
