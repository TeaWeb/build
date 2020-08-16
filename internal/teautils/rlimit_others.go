// +build !linux,!darwin

package teautils

// set resource limit
func SetRLimit(limit uint64) error {
	return nil
}

// set best resource limit value
func SetSuitableRLimit() {

}
