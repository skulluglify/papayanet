package kornet

import m "PapayaNet/papaya/koala/mapping"

// ref:https://www.iana.org/assignments/character-sets/character-sets.xhtml

var AvailableCharsets = m.Keys{

  "ASCII",
  "US-ASCII",
  //"LATIN", -> ISO-8859-1
  //"LATIN-1", -> ISO-8859-1
  "ISO-8859-1",
  "ISO-8859-2",
  "ISO-8859-3",
  "ISO-8859-4",
  "ISO-8859-5",
  "ISO-8859-6",
  "ISO-8859-6-E",
  "ISO-8859-6-I",
  "ISO-8859-7",
  "ISO-8859-8",
  "ISO-8859-8-E",
  "ISO-8859-8-I",
  "ISO-8859-9",
  "ISO-8859-10",
  "ISO-8859-11",
  //"ISO-8859-12", just kidding
  "ISO-8859-13",
  "ISO-8859-14",
  "ISO-8859-15",
  "ISO-8859-16",
  "ISO-10646-UCS-2",
  "ISO-10646-UCS-4",
  "ISO-10646-UCS-Basic",
  "ISO-10646-Unicode-Latin1",
  "UTF-8",
  "UTF-16",
  "UTF-16BE",
  "UTF-16LE",
  "UTF-32",
  "UTF-32BE",
  "UTF-32LE",
  // next type, find your self
}
