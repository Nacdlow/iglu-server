package models

import (
	"testing"
)

func TestFakeDataGeneration(t *testing.T) {
	t.Log("Devices...")
	for i := 0; i < 10; i++ {
		t.Log(GetFakeDevice())
	}
	t.Log("Room...")
	for i := 0; i < 10; i++ {
		t.Log(GetFakeRoom())
	}
	t.Log("Room stats...")
	for i := 0; i < 10; i++ {
		t.Log(GetFakeRoomStat())
	}
	t.Log("Stats...")
	for i := 0; i < 10; i++ {
		t.Log(GetFakeStat())
	}
	t.Log("Users...")
	for i := 0; i < 10; i++ {
		t.Log(GetFakeUser())
	}
}
