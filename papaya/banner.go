package papaya

import "PapayaNet/papaya/utils"

func Banner(version *utils.PnVersion) string {

	lines := []string{
		"",
		"      ____                                _   __     __",
		"     / __ \\____ _____  ____ ___  ______ _/ | / /__  / /_",
		"    / /_/ / __ `/ __ \\/ __ `/ / / / __ `/  |/ / _ \\/ __/",
		"   / ____/ /_/ / /_/ / /_/ / /_/ / /_/ / /|  /  __/ /_",
		"  /_/    \\__,_/ .___/\\__,_/\\__, /\\__,_/_/ |_/\\___/\\__/",
		"             /_/          /____/ " + version.Stringify(),
	}

	var context string

	for _, line := range lines {

		context += utils.PnStrPadEnd(line, 57) + "\n"
	}

	return context
}
