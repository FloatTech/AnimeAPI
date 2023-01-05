package bilibili

import "testing"

func TestGetAllGuard(t *testing.T) {
	guardUser, err := GetAllGuard("628537")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", guardUser)
}
