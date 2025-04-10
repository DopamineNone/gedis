package utils

func ToCommandLine(cmdName string, args ...[]byte) [][]byte {
	result := make([][]byte, len(args)+1)
	result[0] = []byte(cmdName)
	for i := range args {
		result[i+1] = args[i]
	}
	return result
}
