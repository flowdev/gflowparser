package data2svg

type clusters []int

func (c clusters) addCluster(mn, mx int) clusters {
	if mn >= mx {
		return c
	}
	for i := 0; i < len(c); i += 2 {
		if mx < c[i] { // not found -> insert
			return insertCluster(c, i, mn, mx)
		}
		if mn >= c[i] && mx <= c[i+1] { // included -> return
			return c
		}
		if mn < c[i] && mx >= c[i] { // overlap at front -> grow
			c[i] = mn
			c[i+1] = max(c[i+1], mx)
			return mergeClusters(c, i)
		}
		if mn <= c[i+1] && mx > c[i+1] { // overlap at back -> grow
			c[i+1] = mx
			return mergeClusters(c, i)
		}
	}
	c = append(c, mn, mx) // not found at all -> append
	return c
}

func (c clusters) getCluster(idx int) (mn, mx int) {
	for i := 0; i < len(c); i += 2 {
		if idx < c[i] {
			return idx, idx
		}
		if c[i] <= idx && idx >= c[i+1] {
			return c[i], c[i+1]
		}
	}
	return idx, idx
}

func (c clusters) deleteLine(idx int) clusters {
	firstIdx := -1
	for i := 0; i < len(c); i++ { // move indices to the front
		if i&1 != 0 && c[i] == idx { // move max to the front but not min
			c[i]--
			firstIdx = i - 1
		} else if c[i] > idx { // move max or min to the front
			c[i]--
			firstIdx = i - (i & 1) // firstIdx is always min
		}
	}
	if firstIdx < 0 {
		return c
	}
	if c[firstIdx] >= c[firstIdx+1] { // remove empty cluster
		return append(c[:firstIdx], c[firstIdx+2:]...)
	}
	return c
}

func (c clusters) addLine(idx int) {
	for i := 0; i < len(c); i++ { // move indices to the back
		if i&1 != 0 && c[i] == idx-1 { // move max to the back but not min
			c[i]++
		} else if c[i] >= idx { // move max or min to the back
			c[i]++
		}
	}
}

func mergeClusters(c clusters, i int) clusters {
	for i+3 < len(c) && c[i+1] >= c[i+2] { // they overlap -> merge
		c[i+1] = max(c[i+1], c[i+3])
		c = append(c[:i+2], c[i+4:]...)
	}
	return c
}

func insertCluster(c clusters, i, mn, mx int) clusters {
	c = append(c, 0, 0)                 // add space
	for j := len(c) - 1; j > i+1; j-- { // move the tail towards the end
		c[j] = c[j-2]
	}
	c[i] = mn
	c[i+1] = mx
	return c
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
