package key

// todo
type Key string

func (k Key) String() string {
	return string(k)
}

func Parse(s string) Key {
	return Key(s)
}
