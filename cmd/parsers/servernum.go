package parsers

import (
	"errors"
	"fmt"
	"github.com/keyCat/srcds-manager/config"
	"regexp"
	"strconv"
	"strings"
)

// ParseAsServerNum validate and parse string to an array of server numbers. Returns an array of all server numbers from config, if empty string was passed
// example values: 1	1,2,3 	 1,1-2		2-
func ParseAsServerNum(args []string) ([]int, error) {
	var numbers []int
	if len(args) > 0 {
		// parse first argument
		arg := args[0]
		var split = strings.Split(arg, ",")
		for _, elem := range split {
			elem = strings.TrimSpace(elem)
			if elem == "" {
				continue
			}

			matched, _ := regexp.MatchString("^([0-9]{1,}-[0-9]{0,})|([0-9]{1,})$", elem)
			if !matched {
				return numbers, errors.New(fmt.Sprintf("invalid server number argument \"%v\" (must contain only numeric characters, commas or ranges)", elem))
			}
			if strings.Contains(elem, "-") {
				// unwind range
				nums := strings.Split(elem, "-")
				num1, err := strconv.Atoi(nums[0])
				num2 := len(config.Value.Servers)
				if len(nums) > 1 && nums[1] != "" {
					num2, err = strconv.Atoi(nums[1])
				}
				if err != nil {
					return numbers, err
				}
				if num1 >= num2 {
					return numbers, errors.New(fmt.Sprintf("invalid server number argument \"%v\" (incorrect range)", elem))
				}
				for i := num1; i <= num2; i++ {
					numbers = append(numbers, i)
				}
			} else {
				// push number
				num, err := strconv.Atoi(elem)
				if err != nil {
					return numbers, err
				}
				numbers = append(numbers, num)
			}
		}
	} else {
		// received an empty arg, fill in with all possible server numbers from the config
		for _, server := range config.Value.Servers {
			numbers = append(numbers, server.Number)
		}
		return numbers, nil
	}

	return numbers, nil
}
