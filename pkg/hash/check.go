package hash

import (
	"fmt"
)

func Check(should string) error {
	is, err := Calc()
	if err != nil {
		return fmt.Errorf("cannot calculate own hash: %w", err)
	}
	if is != should {
		fmt.Printf("%s (is)\n%s (should)\n", is, should)
		return fmt.Errorf("hashes to not match")
	}
	return nil
}
