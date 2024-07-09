package goenv

type Options func(e *EnvParser)

func WithPrefix(prefix string) Options {
	return func(e *EnvParser) {
		e.prefix = prefix
	}
}

func WithSplitChar(splitChar string) Options {
	return func(e *EnvParser) {
		e.splitChar = splitChar
	}
}

func WithTag(tag string) Options {
	return func(e *EnvParser) {
		e.tag = tag
	}
}
