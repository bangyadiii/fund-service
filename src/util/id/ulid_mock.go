package id

type mockUlid struct {
	GenerateFunc func() string
}

func (m *mockUlid) Generate() string {
	return m.GenerateFunc()
}

func NewMockUlid(generateFunc func() string) IDGenerator {
	return &mockUlid{GenerateFunc: generateFunc}
}
