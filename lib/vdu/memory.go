package vdu

type Memory struct {
	data []byte
}

func (r *Memory) LoadProgram(program []byte) {
	copy(r.data[:], program)
}

func NewMemory(size int) *Memory {
	return &Memory{data: make([]byte, size)}
}
