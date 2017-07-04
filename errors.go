package moldova

// UnsupportedTokenError is returned from the parser when it encounters an unknown token
type UnsupportedTokenError string

// Error implmenets the error interface
func (e UnsupportedTokenError) Error() string {
	return string(e)
}

// InvalidArgumentError is returned from the parser when it encounters an invalid argument
// to a known token
type InvalidArgumentError string

// Error implmenets the error interface
func (e InvalidArgumentError) Error() string {
	return string(e)
}
