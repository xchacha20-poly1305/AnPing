//go:build !linux || android

package anping

func getPermission() error {
	return nil
}
