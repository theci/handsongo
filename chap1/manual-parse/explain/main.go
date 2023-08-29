package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type config struct {
	numTimes   int
	printUsage bool
}

var usageString = fmt.Sprintf(`Usage: %s <integer> [-h|-help]

A greeter application which prints the name you entered <integer> number of times.
`, os.Args[0])

func printUsage(w io.Writer) {
	fmt.Fprintf(w, usageString)
}


// 3. 입력 값을 검증하는 함수
func validateArgs(c config) error {
	if c.numTimes <= 0 && !c.printUsage {   // 0이나 음수를 입력하면 에러를 반환
		return errors.New("Must specify a number greater than 0")
	}
	return nil
}

// 2. 입력 매개변수로 문자열의 슬라이스로 받고 config타입과 error 타입 값을 반환
func parseArgs(args []string) (config, error) {
	var numTimes int		// 인사 횟수를 결정하는 정숫값
	var err error
	c := config{} 			// 객체를 생성하여 인수 데이터를 저장
	if len(args) != 1 { 	// 입력이 빈 경우 에러를 반환
		return c, errors.New("Invalid number of arguments")
	}

	if args[0] == "-h" || args[0] == "-help" { 		// 입력한 인수가 -h, -help인 경우 c와 nil 에러를 반환
		c.printUsage = true			// 사용법을 보여주는지의 여부를 정하는 불리언 값. true면 사용법 출력
		return c, nil
	}

	numTimes, err = strconv.Atoi(args[0])   // 입력한 숫자를 정수로 변환.
	if err != nil {
		return c, err
	}
	c.numTimes = numTimes  // 변환된 정숫값을 c 객체에 할당

	return c, nil
}


// 1. 사용자의 이름을 입력받는 함수
func getName(r io.Reader, w io.Writer) (string, error) {  				// 애플리케이션의 입력과 출력 인터페이스를 매개변수로 받음
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprintf(w, msg)				// w에 프롬프트를 작성
	scanner := bufio.NewScanner(r)	// 스캐너 생성
	scanner.Scan()					// 입력 데이터를 스캔
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()			// 읽은 데이터를 문자열로 반환
	if len(name) == 0 {				// 입력 문자가 공백이면 에러 발생
		return "", errors.New("You didn't enter your name")
	}
	return name, nil			// 사용자의 입력으로 이름을 성공적으로 읽어 들인 경우 이름과 nil 에러를 반환
}

// 5. 화면에 사용자에게 인사를 표출하는 함수
func greetUser(c config, name string, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, msg)
	}
}

// 4. config 타입의 객체 c에 포함된 값에 따라 해당하는 동작을 수행
func runCmd(r io.Reader, w io.Writer, c config) error {
	if c.printUsage {		/// 사용자가 매개변수로 -help 또는 -h를 지정한 경우 함수 호출하고 에러로 nil 반환
		printUsage(w)
		return nil
	}

	name, err := getName(r, w)  // 그 외에는 사용자의 이름을 입력받음
	if err != nil {
		return err
	}
	greetUser(c, name, w)
	return nil
}

func main() {
	c, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		printUsage(os.Stdout)
		os.Exit(1)
	}
	err = validateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		printUsage(os.Stdout)
		os.Exit(1)
	}

	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
