package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)


/*
개선방법
1. 중복된 오류 메시지를 제거
2. 도움말 사용법 메시지를 사용자 정의할 수 있도록 한다
3. 위치 인수를 통해 사용자의 이름을 입력 받을 수 있도록 한다
*/

type config struct {
	numTimes int
	name     string		// 3. 위치 인수를 통해 이름 받기
}


// 중복 오류 메시지가 나온 이유 : fs.SetOutput() 함수에 설정된 writer에 해당 에러를 출력하지만 main() 함수에서도 또 출력된다
var errInvalidPosArgSpecified = errors.New("More than one positional argument specified") 	// 1-1. 사용자 정의 에러값

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

func greetUser(c config, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", c.name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, msg)
	}
}

func runCmd(r io.Reader, w io.Writer, c config) error {
	var err error
	if len(c.name) == 0 {		// 이름이 지정되지 않았거나 이름이 공백인 경우에만 사용자에게 이름을 받도록 함
		c.name, err = getName(r, w)
		if err != nil {
			return err
		}
	}
	greetUser(c, w)
	return nil
}

func validateArgs(c config) error {
	if !(c.numTimes > 0) {
		return errors.New("Must specify a number greater than 0")
	}
	return nil
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Usage = func() {	// 2. FlagSet 객체의 Usage 속성값에 함수를 정의하여 사용자가 메시지를 직접 정의할 수 있다
		var usageString = `
A greeter application which prints the name you entered a specified number of times.

Usage of %s: <options> [name]`
		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if fs.NArg() > 1 {		// parseArg() 함수가 위치 인수를 찾도록 업데이트하고, 찾은 경우에는 config 객체의 name 속성값을 설정
		return c, errInvalidPosArgSpecified		// 1-2. 사용자 정의 에러를 반환
	}
	if fs.NArg() == 1 {
		c.name = fs.Arg(0)
	}
	return c, nil
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		if errors.Is(err, errInvalidPosArgSpecified) {	
			// 1-3. rrors.Is() 함수를 사용하여 반환된 에러값 err가 errInvalidPosArgSpecified값과 일치하는지 확인하고 일치하는 경우에만 오류 메시지를 화면에 출력
			fmt.Fprintln(os.Stdout, err)
		}
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
