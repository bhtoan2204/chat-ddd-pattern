package generator

import "testing"

func TestNormalizeRoutePath(t *testing.T) {
	path := "/relationship/friend-requests/{target_user_id}"
	got := normalizeRoutePath(path)
	want := "/relationship/friend-requests/:target_user_id"
	if got != want {
		t.Fatalf("unexpected normalized path: got %s want %s", got, want)
	}
}
