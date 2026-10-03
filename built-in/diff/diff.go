// Based on Myers diff

package diff

type PartType int

const (
	Equal PartType = iota
	Insert
	Delete
)

type DiffPart[T any] struct {
	Type  PartType
	Value T
}

type coord struct {
	x      int
	y      int
}

func GetDiff[T any](a, b []T, isEqual func(T, T) bool) []DiffPart[T] {
	// This can be improved a lot
	start := coord{0, 0}
	end := coord{len(a), len(b)}

	forwardVisited := map[coord]coord{start: {}}
	backwardVisited := map[coord]coord{end: {}}

	front := []coord{start}
	nextFront := []coord{}

	inverseFront := []coord{end}
	nextInverseFront := []coord{}

	for range len(a) + len(b) {
		for _, value := range front {
			snaked := value

			for snaked.x < len(a) && snaked.y < len(b) && isEqual(a[snaked.x], b[snaked.y]) {
				snaked.x++
				snaked.y++
			}

			if snaked != value {
				if _, found := forwardVisited[snaked]; !found {
					nextFront = append(nextFront, snaked)
					forwardVisited[snaked] = value
				}

				if _, found := backwardVisited[snaked]; found {
					return getPath(a, b, snaked, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.x < len(a) {
				coord := coord{snaked.x + 1, snaked.y}

				if _, found := forwardVisited[coord]; !found {
					nextFront = append(nextFront, coord)
					forwardVisited[coord] = snaked
				}

				if _, found := backwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.y < len(b) {
				coord := coord{snaked.x, snaked.y + 1}

				if _, found := forwardVisited[coord]; !found {
					nextFront = append(nextFront, coord)
					forwardVisited[coord] = snaked
				}

				if _, found := backwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}
		}

		front = front[:0]
		front, nextFront = nextFront, front

		for _, value := range inverseFront {
			snaked := value

			for snaked.x > 0 && snaked.y > 0 && isEqual(a[snaked.x - 1], b[snaked.y - 1]) {
				snaked.x--
				snaked.y--
			}

			if snaked != value {
				if _, found := backwardVisited[snaked]; !found {
					nextInverseFront = append(nextInverseFront, snaked)
					backwardVisited[snaked] = value
				}

				if _, found := forwardVisited[snaked]; found {
					return getPath(a, b, snaked, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.y > 0 {
				coord := coord{snaked.x, snaked.y - 1}

				if _, found := backwardVisited[coord]; !found {
					nextInverseFront = append(nextInverseFront, coord)
					backwardVisited[coord] = snaked
				}

				if _, found := forwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.x > 0 {
				coord := coord{snaked.x - 1, snaked.y}

				if _, found := backwardVisited[coord]; !found {
					nextInverseFront = append(nextInverseFront, coord)
					backwardVisited[coord] = snaked
				}

				if _, found := forwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}
		}

		inverseFront = inverseFront[:0]
		inverseFront, nextInverseFront = nextInverseFront, inverseFront
	}

	return []DiffPart[T]{}
}

func getPath[T any](
	a, b []T,
	commonCoord coord,
	forwardVisited, backwardVisited map[coord]coord,
	start, end coord,
) []DiffPart[T] {
	forwardLength := 0

	for coord := commonCoord; coord != start; {
		parent := forwardVisited[coord]

		if coord.x != parent.x && coord.y != parent.y {
			forwardLength += coord.y - parent.y
		} else {
			forwardLength++
		}

		coord = parent
	}

	backwardLength := 0

	for coord := commonCoord; coord != end; {
		parent := backwardVisited[coord]

		if coord.x != parent.x && coord.y != parent.y {
			backwardLength += parent.y - coord.y
		} else {
			backwardLength++
		}

		coord = parent
	}

	result := make([]DiffPart[T], forwardLength + backwardLength)

	pos := forwardLength - 1

	for coord := commonCoord; coord != start; {
		parent := forwardVisited[coord]

		if coord.x != parent.x && coord.y != parent.y {
			for i := coord.y - parent.y - 1; i >= 0; i-- {
				result[pos] = DiffPart[T]{Equal, b[parent.y + i]}
				pos--
			}
		} else if coord.x != parent.x {
			result[pos] = DiffPart[T]{Delete, a[parent.x]}
			pos--
		} else {
			result[pos] = DiffPart[T]{Insert, b[parent.y]}
			pos--
		}

		coord = parent
	}

	pos = forwardLength

	for coord := commonCoord; coord != end; {
		parent := backwardVisited[coord]

		if coord.x != parent.x && coord.y != parent.y {
			for i := 0; i < parent.y - coord.y; i++ {
				result[pos] = DiffPart[T]{Equal, b[coord.y + i]}
				pos++
			}
		} else if coord.x != parent.x {
			result[pos] = DiffPart[T]{Delete, a[parent.x - 1]}
			pos++
		} else {
			result[pos] = DiffPart[T]{Insert, b[parent.y - 1]}
			pos++
		}

		coord = parent
	}

	return result
}
