package utility

const (
	defaultLength   = 0
	defaultCapacity = 64
)

// Configs ... elements of a slice
type configs struct {
	length   int
	capacity int
}

// Option ... option of a slice
type Option interface {
	Apply(*configs)
}

// Length ... length of a slice
type Length int

// Apply ... apply length of a slice
func (o Length) Apply(c *configs) {
	c.length = int(o)
}

// WithLength ... optional value - length of a slice
func WithLength(v int) Length {
	return Length(v)
}

// Capacity ... capacity of a slice
type Capacity int

// Apply ... apply capacity of a slice
func (o Capacity) Apply(c *configs) {
	c.capacity = int(o)
}

// WithCapacity ... optional value - capacity of a slice
func WithCapacity(v int) Capacity {
	return Capacity(v)
}

var args *configs

func init() {
	args = &configs{
		length:   defaultLength,
		capacity: defaultCapacity,
	}
}

// CombineStrings ... combine strings
func CombineStrings(s []string, options ...Option) string {
	for _, option := range options {
		option.Apply(args)
	}

	c := make([]byte, args.length, args.capacity)

	for _, v := range s {
		c = append(c, []byte(v)...)
	}
	return string(c)
}
