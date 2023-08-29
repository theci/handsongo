package main

import (
	"bufio"
	"errors"
	"flag"   // flag 패키지 사용. flag 패키지는 커맨드 라인 애플리케이션의 인자를 처리하기 위한 표준 동작의 타입과 메서드를 구현
	// 애플리케이션 실행 시에 -h를 지정해주면 그 외의 모든 인수를 무시하고 애플리케이션 사용법을 출력한다
	"fmt"
	"io"
	"os"
)

type config struct {	// parseArgs() 함수에서 -h, -help 인수를 처리하기 때문에 config 구조체에 더 이상 printUsage 필드가 없다
	numTimes int
}

func getName(r io.Reader, w io.Writer) (string, error) {
	scanner := bufio.NewScanner(r)
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprintf(w, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You didn't enter your name")
	}
	return name, nil
}

func greetUser(c config, name string, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, msg)
	}
}

func runCmd(r io.Reader, w io.Writer, c config) error {
	name, err := getName(r, w)
	if err != nil {
		return err
	}
	greetUser(c, name, w)
	return nil
}

func validateArgs(c config) error {
	if !(c.numTimes > 0) {
		return errors.New("Must specify a number greater than 0")
	}
	return nil
}

// mamual-parse와 비교되는 변경점
func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)		// 인수를 파싱하기 위해 FlagSet 객체 생성
	// Flagset의 첫 번째 매개변수는 애플리케이션 자체 이름. 두 번째 매개변수는 에러 처리 방법
	fs.SetOutput(w)		// FlagSet 객체의 출력메시지 값으로 w로 지정하여 함수의 동작을 검증하기 위한 유닛 테스트를 작성할 수 있다
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")	// 첫 번째 플래그 옵션 정의.
	// 첫 번째 매개변수는 해당 값을 저장할 포인터 값. 두 번째는 옵션 이름. 마지막 매개변수는 사용자에게 이 옵션의 목적을 나타내는 문자열
	err := fs.Parse(args)	// args[] 슬라이스를 매개변수로 함수 호출. 슬라이스의 요소들을 읽고 flag 값을 설정된 변수아 할당
	if err != nil {
		return c, err
	}
	if fs.NArg() != 0 {		// 인사 프로그램은 별도의 인수를 필요로 하지 않기 때문에 하나 이상 값이 지정된 경우 에러 출력
		return c, errors.New("Positional arguments specified")
	}
	return c, nil
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	err = validateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
