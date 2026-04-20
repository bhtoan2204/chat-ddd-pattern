package generator

import "testing"

func TestModuleForUsecase(t *testing.T) {
	cases := []struct {
		name       string
		usecase    string
		wantFsRoot string
	}{
		{
			name:       "auth aliases to account",
			usecase:    "AuthUsecase",
			wantFsRoot: "core/modules/account",
		},
		{
			name:       "message aliases to room",
			usecase:    "MessageUsecase",
			wantFsRoot: "core/modules/room",
		},
		{
			name:       "relationship resolves generically",
			usecase:    "RelationshipUsecase",
			wantFsRoot: "core/modules/relationship",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			module, err := moduleForUsecase(tc.usecase)
			if err != nil {
				t.Fatalf("moduleForUsecase returned error: %v", err)
			}
			if module.FsRoot != tc.wantFsRoot {
				t.Fatalf("unexpected fs root: got %s want %s", module.FsRoot, tc.wantFsRoot)
			}
		})
	}
}
