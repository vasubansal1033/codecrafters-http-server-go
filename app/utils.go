package main

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == token[0] {
			return i + 1, data[:i], nil
		}
	}
	// If at end of file and no more splitting possible, return the remaining data
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	// If no split found yet and not at end of file, request more data
	return 0, nil, nil
}
