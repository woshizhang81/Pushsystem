package utils

import (
	"bytes"
	"fmt"
)

// 用于存储数字的集和
type IntSet struct {
	words  []uint64
	length int
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && (s.words[word]&(1<<bit)) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	// 判断 x 是否已经存在 s 中
	if s.words[word]&(1<<bit) == 0 {
		s.words[word] |= 1 << bit
		s.length++
	}
}

func (s *IntSet) UnionWith(t *IntSet) {
	originLenOfSet := len(s.words)
	for i, v := range t.words {
		if i < originLenOfSet {
			if v != 0 {
				for j := uint(0); j < 64; j++ {
					if s.words[i]&(1<<j) == 0 && (v&(1<<j) != 0) {
						s.length++
					}
				}
				s.words[i] |= v
			}
		} else {
			s.words = append(s.words, v)
			for j := uint(0); j < 64; j++ {
				if v&(1<<j) != 0 {
					s.length++
				}
			}
		}
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, v := range s.words {
		if v == 0 {
			continue
		}
		for j := uint(0); j < 64; j++ {
			if v&(1<<j) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*uint(i)+j)
			}
		}
	}
	buf.WriteByte('}')
	fmt.Fprintf(&buf, "\nLength: %d", s.length)
	return buf.String()
}

func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	if word < len(s.words) {
		if s.words[word] & (1 << bit) != 0 {
			s.words[word] &= ^(1 << bit)
			s.length--
		}
	}
}

func (s *IntSet) Len() int {
	return s.length
}

func (s *IntSet) Clear() {
	s.words = nil
	s.length = 0
}

func (s *IntSet) Copy() (cp *IntSet) {
	cp = new(IntSet)
	cp.words = make([]uint64, len(s.words))
	cp.length = s.length
	copy(cp.words, s.words)
	return cp
}