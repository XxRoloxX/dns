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

	index := startIndex

	valueBits := toBinaryArray(value)

	valueBitIndex := 0

	for index <= endIndex {

		if valueBitIndex >= len(valueBits) {

			// Pad with zeros
			s.Set(index, false)
			index++
			continue
		}

		s.Set(index, valueBits[valueBitIndex])
		index++
		valueBitIndex++
	}
}
