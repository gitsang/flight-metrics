package syntax

func If[T any](condition bool, trueOut, falseOut T) T {
	if condition {
		return trueOut
	}
	return falseOut
}
