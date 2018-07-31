package converter

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type numConverter struct {
}

func NewNumberConverter() NumberConverter {
	return &numConverter{}
}

func (numConv *numConverter) ConvertToWords(number int64) (string, error) {
	var result []string
	sn := strconv.FormatInt(number, 10)
	leftOver := len(sn) % 3
	result = append(result, sn[:leftOver])
	for i := leftOver; i < len(sn); i += 3 {
		result = append(result, sn[i:i+3])
	}

	nc := getWordIDNNumber()
	ml := len(result) - 1
	var s string

	for i, val := range result {
		if len(val) > 0 && IsLetter(val) {
			lvl := ml - i
			if val != "000" {
				s += getCardinal(val, nc)
				s += " "
				s += getCardinalPronoun(lvl)
			}
			s += " "
		}

	}

	return standardizeSpaces(s), nil
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func getCardinalPronoun(lvl int) string {
	// 312 triliun 133 milyar 412 juta 353 ribu 123
	switch lvl {
	case 1:
		return "ribu"
	case 2:
		return "juta"
	case 3:
		return "milyar"
	case 4:
		return "trilyun"
	default:
		return ""
	}
}

func getCardinal(number string, nc map[byte]string) string {
	var result string
	if len(number) == 3 {
		if number[0] == '1' {
			result += "seratus "
		} else {
			result += nc[number[0]]
			result += " ratus "
		}
	}

	if len(number) >= 2 {
		if number == "11" {
			result += "sebelas "
			return result
		} else {
			lastNum := number[len(number)-1]
			if len(number) > 2 {
				middleNum := number[1]
				if middleNum == '1' {
					// belasan
					result += nc[lastNum]
					result += " belas "
					return result

				} else if middleNum == '0' {
					result += nc[lastNum]
					result += " "
					return result

				} else {
					result += nc[middleNum]
					result += " puluh "
					result += nc[lastNum]
					return result

				}
			} else {
				middleNum := number[0]
				if middleNum == '1' {
					// belasan
					result += nc[lastNum]
					result += " belas "
					return result

				} else if middleNum == '0' {
					result += nc[lastNum]
					result += " "
					return result

				} else {
					result += nc[middleNum]
					result += " puluh "
					result += nc[lastNum]
					return result

				}
			}
		}

	}

	if len(number) == 1 {
		if len(number) > 1 {
			result += nc[number[1]]
		} else {
			result += nc[number[0]]
		}
		return result
	}

	return ""

}

func (numConv *numConverter) ConvertToNumber(input string) (int64, error) {
	var tmp []string

	skip := 0
	pointer := -1
	numberID := getIDNNumber()
	numbers := strings.Fields(input)
	for cur, num := range numbers {
		if skip > 0 {
			skip--
			continue
		}
		if skip == 0 {
			isHighCardinalPronouns := IsContainWordHighCardinalPronouns(num)
			if !isHighCardinalPronouns {
				numS, sk, err := scanWords(numbers[cur:], numberID)
				if err != nil {
					return -1, errors.New("error scanning words to numerical character")
				}

				skip = sk
				tmp = append(tmp, numS)
				continue
			}

			for in := range tmp {
				if in > pointer {
					res := tmp[in]
					switch num {
					case "ribu":
						res += "000"
						break
					case "juta":
						res += "000000"
						break
					case "miliar":
						res += "000000000"
						break
					case "milyar":
						res += "000000000"
						break
					case "triliun":
						res += "000000000000"
						break
					case "trilyun":
						res += "000000000000"
						break
					}
					tmp[in] = res
					pointer = in
				}
			}

		}
	}
	var result int64
	result = 0
	for _, val := range tmp {
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return -1, errors.New("error converting result to number value")
		}
		result += num
	}
	return result, nil
}

func IsContainWordHighCardinalPronouns(input string) bool {
	existOperation := []string{"ribu", "juta", "miliar", "milyar", "triliun", "trilyun"}
	for _, val := range existOperation {
		if input == val {
			return true
		}
	}
	return false
}

func scanWords(sisa []string, num map[string]string) (string, int, error) {
	var result string
	var tmp string
	if len(sisa) > 0 {
		// check first word
		tmp = num[sisa[0]]
	}

	if len(sisa) > 1 {
		// check cardinal pronouns
		switch sisa[1] {
		case "belas":
			result = "1" + tmp
			return result, 1, nil
		case "puluh":
			result = tmp + "0"
			break
		case "ratus":
			result = tmp + "00"
			return result, 1, nil
		default:
			result = tmp + num[sisa[1]]
			return result, 1, nil
		}
	} else {
		return tmp, 0, nil
	}

	if len(sisa) > 2 {
		// only puluh and ratus cardinal pronouns
		cardinalPronouns := sisa[1]
		if cardinalPronouns == "puluh" {
			if IsContainWordHighCardinalPronouns(sisa[2]) {
				return result, 1, nil
			}
			tmpNum, err := strconv.Atoi(result)
			if err != nil {
				return "", 0, err
			}
			leftOver, err := strconv.Atoi(num[sisa[2]])
			if err != nil {
				return "", 0, err
			}
			tens := tmpNum + leftOver
			result = strconv.Itoa(tens)

			return result, 2, nil
		}
	} else {
		return result, 1, nil
	}

	return "", 0, errors.New("invalid convert words to numerical character")
}

func getIDNNumber() map[string]string {
	numbers := make(map[string]string)
	numbers["satu"] = "1"
	numbers["dua"] = "2"
	numbers["tiga"] = "3"
	numbers["empat"] = "4"
	numbers["lima"] = "5"
	numbers["enam"] = "6"
	numbers["tujuh"] = "7"
	numbers["delapan"] = "8"
	numbers["sembilan"] = "9"
	numbers["sepuluh"] = "10"
	numbers["sebelas"] = "11"

	numbers["seratus"] = "100"
	numbers["seribu"] = "1000"

	return numbers
}

func getWordIDNNumber() map[byte]string {
	numbers := make(map[byte]string)
	numbers['1'] = "satu"
	numbers['2'] = "dua"
	numbers['3'] = "tiga"
	numbers['4'] = "empat"
	numbers['5'] = "lima"
	numbers['6'] = "enam"
	numbers['7'] = "tujuh"
	numbers['8'] = "delapan"
	numbers['9'] = "sembilan"
	numbers['0'] = ""

	return numbers
}
