package reuse

type reuseImpl struct {
	recv   int
	buffer int
}

func (r reuseImpl) Recv() int {
	return r.recv
}
func (r reuseImpl) Buffer() int {
	return r.buffer
}
