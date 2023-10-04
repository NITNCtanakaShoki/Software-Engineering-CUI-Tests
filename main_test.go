package main

import (
	"os/exec"
	"testing"
)

func Test正常系(t *testing.T) {
	// Arrange
	testCases := []TestCase{
		// 足し算のテストケース
		{"1", "+", "2", "3\n"},                              // 自然数同士の足し算
		{"1", "+", "0.2", "1.2\n"},                          // 自然数と小数の足し算
		{"0.3", "+", "2", "2.3\n"},                          // 小数と自然数の足し算
		{"0.1", "+", "0.2", "0.3\n"},                        // Decimalでないと対応できない小数同士の足し算
		{"-1", "+", "0.2", "0.8\n"},                         // 負と正の足し算
		{"3", "+", "-4.5", "-1.5\n"},                        // 正と負の足し算
		{"-3", "+", "-4.5", "-7.5\n"},                       // 負と負の足し算
		{"9999.9999", "+", "9999.9999", "19,999.9998\n"},    // 最大値同士の足し算
		{"-9999.9999", "+", "-9999.9999", "-19,999.9998\n"}, // 最小値同士の足し算
		{"4.5", "+", "0", "4.5\n"},                          // 数値とゼロとの足し算
		{"0", "+", "5.123", "5.123"},                        // ゼロと数値との足し算

		// 引き算のテストケース
		{"5", "-", "3", "2\n"},                             // 自然数同士の引き算
		{"5.5", "-", "3.2", "2.3\n"},                       // 小数同士の引き算
		{"5", "-", "0.2", "4.8\n"},                         // 自然数と小数の引き算
		{"0.5", "-", "0.3", "0.2\n"},                       // 小数と自然数の引き算
		{"-5", "-", "3", "-8\n"},                           // 負と正の引き算
		{"5", "-", "-3", "8\n"},                            // 正と負の引き算
		{"-5", "-", "-3", "-2\n"},                          // 負と負の引き算
		{"9999.9999", "-", "-9999.9999", "19,999.9998\n"},  // 最大値-最小値
		{"-9999.9999", "-", "9999.9999", "-19,999.9998\n"}, // 最小値-最大値
		{"4.5", "-", "0", "4.5\n"},                         // 数値とゼロとの足し算
		{"0", "-", "5.123", "-5.123"},                      // ゼロと数値との足し算

		// 掛け算のテストケース
		{"5", "*", "3", "15\n"},                           // 自然数同士の掛け算
		{"5.5", "*", "3.2", "17.6\n"},                     // 小数同士の掛け算
		{"5", "*", "0.2", "1\n"},                          // 自然数と小数の掛け算
		{"0.5", "*", "3", "1.5\n"},                        // 小数と自然数の掛け算
		{"-5", "*", "3", "-15\n"},                         // 負と正の掛け算
		{"5", "*", "-3", "-15\n"},                         // 正と負の掛け算
		{"-5", "*", "-3", "15\n"},                         // 負と負の掛け算
		{"9999.9999", "*", "-9999.9999", "-99,999,998\n"}, // 最大値×最小値
		{"4.5", "*", "0", "0n"},                           // 数値とゼロとの掛け算
		{"0", "*", "5.123", "0"},                          // ゼロと数値との掛け算

		// 割り算のテストケース
		{"15", "/", "3", "5\n"},                  // 自然数同士の割り算
		{"5.5", "/", "1.1", "5\n"},               // 小数同士の割り算
		{"5", "/", "0.2", "25\n"},                // 自然数と小数の割り算
		{"0.5", "/", "0.1", "5\n"},               // 小数と自然数の割り算
		{"-15", "/", "3", "-5\n"},                // 負と正の割り算
		{"15", "/", "-3", "-5\n"},                // 正と負の割り算
		{"-15", "/", "-3", "5\n"},                // 負と負の割り算
		{"9999.9999", "/", "-9999.9999", "-1\n"}, // 最大値÷最小値
		{"0", "/", "5.123", "0"},                 // ゼロを割る
	}

	for _, testCase := range testCases {
		// Act
		got, status, err := execute(testCase.lhs, testCase.operator, testCase.rhs)
		if err != nil {
			t.Fatalf("command output err: %s", err.Error())
		}

		// Assert
		if status != 0 {
			t.Errorf("when %s %s %s, want exit status 0 but actual %d", testCase.lhs, testCase.operator, testCase.rhs, status)
		}

		if testCase.want != got {
			t.Errorf("want %s %s %s = %s but actual %s", testCase.lhs, testCase.operator, testCase.rhs, testCase.want, got)
		}
	}
}

func Test引数の数が4つでない場合(t *testing.T) {
	// Arrange
	argCases := [][]string{
		{},
		{"1"},
		{"1", "+"},
		{"1", "+", "2", "3"},
		{"1", "+", "2", "3", "4"},
	}
	for _, argCase := range argCases {
		// Act
		got, status, err := execute(argCase...)
		if err != nil {
			t.Fatalf("command output err: %s", err.Error())
		}

		// Assert
		if status != 1 {
			t.Errorf("when %v, want exit status 1 but actual %d", argCase, status)
		}

		want := "引数値の数は3つである必要があります\n"
		if want != got {
			t.Errorf("want: \"%s\" but got: \"%s\"", want, got)
		}
	}
}

func Test数値の制限(t *testing.T) {
	// Arrange
	testCases := []TestCase{
		{"10000", "+", "2", "左辺: 10000は最大値9999.9999を上回っています\n"},           // 左辺が最大値を超える
		{"10000", "-", "10000", "左辺: 10000は最大値9999.9999を上回っています\n"},       // 両辺が最大値を超える
		{"2", "*", "10000", "右辺: 10000は最大値9999.9999を上回っています\n"},           // 右辺のみが最大値を超える
		{"-10000", "/", "2", "左辺: -10000は最小値-9999.9999を下回っています\n"},        // 左辺が最小値未満
		{"-10000", "+", "-10000", "左辺: -10000は最小値-9999.9999を下回っています\n"},   // 両辺が最小値未満
		{"2", "-", "-10000", "右辺: -9999.9999を下回っています\n"},                  // 右辺のみが最小値未満
		{"1.00001", "*", "2", "左辺: 1.00001は小数点第4位以下で収まっていません\n"},          // 左辺が小数点第4位を超える
		{"1.000004", "/", "1.000002", "左辺: 1.000004は小数点第4位以下で収まっていません\n"}, // 両辺が小数点第4位を超える
		{"1", "+", "1.000002", "右辺: 1.000002は小数点第4位以下で収まっていません\n"},        // 右辺のみが小数点第4位を超える
		{"123hello", "-", "1.000002", "左辺: 123helloは数値である必要があります\n"},      // 左辺が数値ではない
		{"123hello", "*", "456world", "左辺: 123helloは数値である必要があります\n"},      // 両辺が数値ではない
		{"123", "/", "456world", "右辺: 456worldは数値である必要があります\n"},           // 両辺が数値ではない
	}
	for _, testCase := range testCases {
		// Act
		got, status, err := execute(testCase.lhs, testCase.operator, testCase.rhs)
		if err != nil {
			t.Fatalf("command output err: %s", err.Error())
		}

		// Assert
		if status != 1 {
			t.Errorf("when %s %s %s, want exit status 1 but actual %d", testCase.lhs, testCase.operator, testCase.rhs, status)
		}

		if testCase.want != got {
			t.Errorf("want: \"%s\" but got: \"%s\"", testCase.want, got)
		}
	}
}

func Test存在しない演算子(t *testing.T) {
	// Arrange
	testCases := []TestCase{
		{"3", "++", "2", "演算子: ++はサポートされていません\n"},
		{"3", "**", "2", "演算子: **はサポートされていません\n"},
	}
	for _, testCase := range testCases {
		// Act
		got, status, err := execute(testCase.lhs, testCase.operator, testCase.rhs)
		if err != nil {
			t.Fatalf("command output err: %s", err.Error())
		}

		// Assert
		if status != 1 {
			t.Errorf("when %s %s %s, want exit status 1 but actual %d", testCase.lhs, testCase.operator, testCase.rhs, status)
		}

		if testCase.want != got {
			t.Errorf("want: \"%s\" but got: \"%s\"", testCase.want, got)
		}
	}
}

func Testゼロ除算(t *testing.T) {
	want := "0で割ることはできません\n"

	// Act
	got, status, err := execute("1", "/", "0")
	if err != nil {
		t.Fatalf("command output err: %s", err.Error())
	}

	// Assert
	if status != 1 {
		t.Errorf("want exit status 1, but %d", status)
	}

	if want != got {
		t.Errorf("want: \"%s\" but got: \"%s\"", want, got)
	}
}

func execute(args ...string) (string, int, error) {
	cmd := exec.Command("./calculator", args...)
	out, err := cmd.Output()
	return string(out), cmd.ProcessState.ExitCode(), err
}

type TestCase struct {
	lhs      string
	operator string
	rhs      string
	want     string
}
