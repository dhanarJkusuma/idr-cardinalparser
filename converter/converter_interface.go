package converter

type NumberConverter interface {
	ConvertToNumber(number string) (int64, error)
	ConvertToWords(number int64) (string, error)
}
