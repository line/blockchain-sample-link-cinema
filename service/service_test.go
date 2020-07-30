package service

import "testing"

func TestRegex(t *testing.T) {
	testdata := [][]string{
		{"a1"},                                    //true
		{"123"},                                   //true
		{"abc"},                                   //true
		{""},                                      //false
		{" "},                                     //false
		{"a 1"},                                   //false
		{"a.1"},                                   //false
		{"a-1"},                                   //true
		{"a_1"},                                   //true
		{"a=1"},                                   //true
		{"a1", "123", "abc", "a-1", "a_1", "a=1"}, //true
		{"a1", "123", "abc", "a 1"},               //false
		{"", " ", "a 1", "a.1"},                   //false
	}

	expectResult := []bool{
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		true,
		true,
		true,
		true,
		false,
		false,
	}

	for i, data := range testdata {
		result := checkUrlParam(data...)
		if result != expectResult[i] {
			t.Error("Unexpected result", data, result, expectResult[i])
		}
	}
}
