package processor

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCountStatsLines(t *testing.T) {
	fileJob := FileJob{
		Content: []byte(""),
		Lines:   0,
	}

	// Both tokei and sloccount count this as 0 so lets follow suit
	// cloc ignores the file itself because it is empty
	countStats(&fileJob)
	if fileJob.Lines != 0 {
		t.Errorf("Zero lines expected got %d", fileJob.Lines)
	}

	// Interestingly this file would be 0 lines in "wc -l" because it only counts newlines
	// all others count this as 1
	fileJob.Lines = 0
	fileJob.Content = []byte("a")
	countStats(&fileJob)
	if fileJob.Lines != 1 {
		t.Errorf("One line expected got %d", fileJob.Lines)
	}

	fileJob.Lines = 0
	fileJob.Content = []byte("a\n")
	countStats(&fileJob)
	if fileJob.Lines != 1 {
		t.Errorf("One line expected got %d", fileJob.Lines)
	}

	// tokei counts this as 1 because its still on a single line unless something follows
	// the newline its still 1 line
	fileJob.Lines = 0
	fileJob.Content = []byte("1\n")
	countStats(&fileJob)
	if fileJob.Lines != 1 {
		t.Errorf("One line expected got %d", fileJob.Lines)
	}

	fileJob.Lines = 0
	fileJob.Content = []byte("1\n2\n")
	countStats(&fileJob)
	if fileJob.Lines != 2 {
		t.Errorf("Two lines expected got %d", fileJob.Lines)
	}

	fileJob.Lines = 0
	fileJob.Content = []byte("1\n2\n3")
	countStats(&fileJob)
	if fileJob.Lines != 3 {
		t.Errorf("Three lines expected got %d", fileJob.Lines)
	}

	content := ""
	for i := 0; i < 5000; i++ {
		content += "a\n"
		fileJob.Lines = 0
		fileJob.Content = []byte(content)
		countStats(&fileJob)
		if fileJob.Lines != int64(i+1) {
			t.Errorf("Expected %d got %d", i+1, fileJob.Lines)
		}
	}
}

func TestCountStatsCode(t *testing.T) {
	fileJob := FileJob{
		Content: []byte(""),
		Code:    0,
	}

	// Both tokei and sloccount count this as 0 so lets follow suit
	// cloc ignores the file itself because it is empty
	countStats(&fileJob)
	if fileJob.Code != 0 {
		t.Errorf("Zero lines expected got %d", fileJob.Code)
	}

	// Interestingly this file would be 0 lines in "wc -l" because it only counts newlines
	// all others count this as 1
	fileJob.Code = 0
	fileJob.Content = []byte("a")
	countStats(&fileJob)
	if fileJob.Code != 1 {
		t.Errorf("One line expected got %d", fileJob.Code)
	}

	fileJob.Code = 0
	fileJob.Content = []byte("a\n")
	countStats(&fileJob)
	if fileJob.Code != 1 {
		t.Errorf("One line expected got %d", fileJob.Code)
	}

	// tokei counts this as 1 because its still on a single line unless something follows
	// the newline its still 1 line
	fileJob.Code = 0
	fileJob.Content = []byte("1\n")
	countStats(&fileJob)
	if fileJob.Code != 1 {
		t.Errorf("One line expected got %d", fileJob.Code)
	}

	fileJob.Code = 0
	fileJob.Content = []byte("1\n2\n")
	countStats(&fileJob)
	if fileJob.Code != 2 {
		t.Errorf("Two lines expected got %d", fileJob.Code)
	}

	fileJob.Code = 0
	fileJob.Content = []byte("1\n2\n3")
	countStats(&fileJob)
	if fileJob.Code != 3 {
		t.Errorf("Three lines expected got %d", fileJob.Code)
	}

	content := ""
	for i := 0; i < 5000; i++ {
		content += "a\n"
		fileJob.Code = 0
		fileJob.Content = []byte(content)
		countStats(&fileJob)
		if fileJob.Code != int64(i+1) {
			t.Errorf("Expected %d got %d", i+1, fileJob.Code)
		}
	}
}

func TestCountStatsBlankLines(t *testing.T) {
	fileJob := FileJob{
		Content: []byte(""),
		Blank:   0,
	}

	countStats(&fileJob)
	if fileJob.Blank != 0 {
		t.Errorf("Zero lines expected got %d", fileJob.Blank)
	}

	fileJob.Blank = 0
	fileJob.Content = []byte(" ")
	countStats(&fileJob)
	if fileJob.Blank != 1 {
		t.Errorf("One line expected got %d", fileJob.Blank)
	}

	fileJob.Blank = 0
	fileJob.Content = []byte("\n")
	countStats(&fileJob)
	if fileJob.Blank != 1 {
		t.Errorf("One line expected got %d", fileJob.Blank)
	}

	fileJob.Blank = 0
	fileJob.Content = []byte("\n ")
	countStats(&fileJob)
	if fileJob.Blank != 2 {
		t.Errorf("Two line expected got %d", fileJob.Blank)
	}

	fileJob.Blank = 0
	fileJob.Content = []byte("            ")
	countStats(&fileJob)
	if fileJob.Blank != 1 {
		t.Errorf("One line expected got %d", fileJob.Blank)
	}

	fileJob.Blank = 0
	fileJob.Content = []byte("            \n             ")
	countStats(&fileJob)
	if fileJob.Blank != 2 {
		t.Errorf("Two lines expected got %d", fileJob.Blank)
	}

	fileJob.Blank = 0
	fileJob.Content = []byte("\r\n\r\n")
	countStats(&fileJob)
	if fileJob.Blank != 2 {
		t.Errorf("Two lines expected got %d", fileJob.Blank)
	}

	fileJob.Blank = 0
	fileJob.Content = []byte("\r\n")
	countStats(&fileJob)
	if fileJob.Blank != 1 {
		t.Errorf("One line expected got %d", fileJob.Blank)
	}
}

func TestCountStatsComplexityCount(t *testing.T) {
	fileJob := FileJob{
		Content: []byte("f"),
		Blank:   0,
	}

	countStats(&fileJob)
}

//////////////////////////////////////////////////
// Benchmarks Below
//////////////////////////////////////////////////

func BenchmarkCountStatsLinesEmpty(b *testing.B) {
	fileJob := FileJob{
		Content: []byte(""),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesSingleChar(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("a"),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesTwoLines(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("a\na"),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesThreeLines(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("a\na\na"),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesShortLine(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("1234567890"),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesShortEmptyLine(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("          "),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesThreeShortLines(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("1234567890\n1234567890\n1234567890"),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesThreeShortEmptyLines(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("          \n          \n          "),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesLongLine(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesLongMixedLine(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("1234567890          1234567890          1234567890          1234567890          1234567890          "),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesLongAlternateLine(b *testing.B) {
	fileJob := FileJob{
		Content: []byte("a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a a "),
	}

	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesFiveHundredLongLines(b *testing.B) {
	b.StopTimer()
	content := ""
	for i := 0; i < 500; i++ {
		content += "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890\n"
	}

	fileJob := FileJob{
		Content: []byte(content),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesFiveHundredLongLinesTriggerComplexityIf(b *testing.B) {
	b.StopTimer()
	content := ""
	for i := 0; i < 500; i++ {
		content += "iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii\n"
	}

	fileJob := FileJob{
		Content: []byte(content),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesFiveHundredLongLinesTriggerComplexityFor(b *testing.B) {
	b.StopTimer()
	content := ""
	for i := 0; i < 500; i++ {
		content += "fofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofofo\n"
	}

	fileJob := FileJob{
		Content: []byte(content),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCountStatsLinesSixHundredLongLinesMixed(b *testing.B) {
	b.StopTimer()
	content := ""
	for i := 0; i < 200; i++ {
		content += "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890\n"
		content += "1234567890          1234567890          1234567890          1234567890          1234567890          \n"
		content += "                                                                                                    \n"
	}

	fileJob := FileJob{
		Content: []byte(content),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		countStats(&fileJob)
	}
}

func BenchmarkCheckByteEqualityReflect(b *testing.B) {
	b.StopTimer()
	one := []byte("for")
	two := []byte("for")

	count := 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		equal := reflect.DeepEqual(one, two)

		if equal {
			count++
		}
	}

	b.Log(count)
}

func BenchmarkCheckByteEqualityBytes(b *testing.B) {
	b.StopTimer()
	one := []byte("for")
	two := []byte("for")

	count := 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		equal := bytes.Equal(one, two)

		if equal {
			count++
		}
	}

	b.Log(count)
}

func BenchmarkCheckByteEqualityLoop(b *testing.B) {
	b.StopTimer()
	one := []byte("for")
	two := []byte("for")

	count := 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		equal := true
		for in, by := range one {
			if by != two[in] {
				equal = false
				break
			}
		}

		if equal {
			count++
		}
	}

	b.Log(count)
}
