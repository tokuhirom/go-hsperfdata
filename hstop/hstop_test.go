package main

import (
	"testing"
)

func Test_shorten_main_class_name(t *testing.T) {
	if shorten_main_class_name("hoge", 15) != "hoge" {
		t.Fail()
	}

	shorten := shorten_main_class_name("abcdefghijklmnopqrstu", 15)
	if shorten != "ghijklmnopqrstu" {
		t.Errorf("shorten failed: %s", shorten)
	}
	if len(shorten) != 15 {
		t.Errorf("shorten failed: %s", shorten)
	}
}
