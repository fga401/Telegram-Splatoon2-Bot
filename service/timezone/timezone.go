package timezone

// Minute returns the offset of timezone against UTC.
func (t Timezone) Minute() int {
	return int(t)
}
