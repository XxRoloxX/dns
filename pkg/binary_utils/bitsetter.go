package binary

func IntPow(n, m int) int {
	if m == 0 {
		return 1
	}

	if m == 1 {
		return n
	}

	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

type BitSetter struct {
	b byte
}

func NewBitSetter(b byte) *BitSetter {
	return &BitSetter{
		b: b,
	}
}

func (s *BitSetter) Byte() byte {
	return s.b
}

func (s *BitSetter) Set(idx uint8, value bool) {
	position := byte(IntPow(2, int(7-idx)))
	if value {
		s.b = s.b | position
		return
	}

	s.b = s.b & ^position
}

func toBinaryArray(value uint8) []bool {

	bits := make([]bool, 0)
	for value > 0 {
		bit := value % 2
		bits = append([]bool{bit != 0}, bits...)
		value /= 2
	}

	return bits
}

func (s *BitSetter) SetRange(startIndex uint8, endIndex uint8, value uint8) {

	// To prevent underflow
	var index int = int(endIndex)

	valueBits := toBinaryArray(value)

	valueBitIndex := len(valueBits) - 1

	for index >= int(startIndex) {

		if valueBitIndex < 0 {

			// Pad with zeros
			s.Set(uint8(index), false)
			index--
			continue
		}

		s.Set(uint8(index), valueBits[valueBitIndex])
		index--
		valueBitIndex--
	}
}
