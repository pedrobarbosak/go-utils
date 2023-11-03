package defaults

const tag = "default"
const ignoreTag = "-"

type Config struct {
	IgnoreOnMissing bool
	Tag             string
}

func DefaultConfig() Config {
	return Config{
		IgnoreOnMissing: true,
		Tag:             tag,
	}
}
