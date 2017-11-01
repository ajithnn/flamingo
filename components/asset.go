package components


type Asset interface {
	Process(string, func())
}

var TypeMap =  map[string]Asset{
  "Media": Media{},
  "Meta": Meta{},
  "Transcode": Transcode{},
  "Graphics": Graphic{},
  "Subtitles": Subtitle{},
  "Audio": Audio{},
  "Track": Track{},
}
