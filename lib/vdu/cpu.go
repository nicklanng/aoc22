package vdu

const (
	InstNil  byte = iota
	InstNoop
	InstAddX
)

type CPU struct {
	ram    *Memory
	memPtr int

	cycleCount int
	cycleDebt uint8

	instRegister byte
	xRegister    int
}

func NewCPU(ram *Memory) *CPU {
	return &CPU{ram: ram, xRegister: 1}
}

func (cpu *CPU) ExecuteCycle() bool {
	if cpu.cycleDebt == 0 {
		// load next instruction
		cpu.instRegister = cpu.ram.data[cpu.memPtr]
		cpu.memPtr++
		switch cpu.instRegister {
		case InstNil:
			return false
		case InstNoop:
			cpu.cycleDebt = 1
		case InstAddX:
			cpu.cycleDebt = 2
		default:
			panic("unknown instruction")
		}
	}

	cpu.cycleCount++
	cpu.cycleDebt--
	if cpu.cycleDebt > 0 {
		return true
	}

	// perform instruction
	switch cpu.instRegister {
	case InstNil:
		return false
	case InstNoop:
		return true
	case InstAddX:
		x := int8(cpu.ram.data[cpu.memPtr])
		cpu.xRegister += int(x)
		cpu.memPtr++
	}

	return true
}

func (cpu *CPU) Cycle() int {
	return cpu.cycleCount
}

func (cpu *CPU) XRegister() int {
	return cpu.xRegister
}