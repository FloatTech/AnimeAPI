package bilibili

import "testing"

func TestVideoInfo(t *testing.T) {
	t.Log(VideoInfo("10007"))
	t.Log(VideoInfo("BV1xx411c7mD"))
	t.Log(VideoInfo("bv1xx411c7mD"))
	t.Log(VideoInfo("BV1mF411j7iU"))
}
