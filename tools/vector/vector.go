package vector

type (
	Vector struct {
		ElementCount int
		ArraySize    int
		Array        []interface{}
	}
)

func (this *Vector) Len() int {
	return this.ElementCount
}
