package util

import (
  "skfw/papaya/koala"
)

func Banner(version koala.KVersionImpl) string {

  lines := []string{
    "",
    "      ____                                _   __     __",
    "     / __ \\____ _____  ____ ___  ______ _/ | / /__  / /_",
    "    / /_/ / __ `/ __ \\/ __ `/ / / / __ `/  |/ / _ \\/ __/",
    "   / ____/ /_/ / /_/ / /_/ / /_/ / /_/ / /|  /  __/ /_",
    "  /_/    \\__,_/ .___/\\__,_/\\__, /\\__,_/_/ |_/\\___/\\__/",
    "             /_/          /____/ v" + version.String(),
  }

  var context string

  for _, line := range lines {

    context += koala.KStrPadEnd(line, 57) + "\n"
  }

  return context
}
