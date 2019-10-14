package timemap

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTimeMapSet(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		values []string
		want   bool
	}{
		{
			name:   "successfully setting a value in a key",
			key:    "name",
			values: []string{"Pavlos"},
			want:   true,
		},
		{
			name:   "successfully setting multiple values in a key",
			key:    "name",
			values: []string{"Pavlos", "Yannis", "Theo"},
			want:   true,
		},
		{
			name:   "successfully setting a nil value in a key",
			key:    "name",
			values: nil,
			want:   true,
		},
		{
			name:   "failing setting a value to an empty key",
			key:    "",
			values: []string{"Pavlos"},
			want:   false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := New()
			if tc.values == nil {
				ok := l.Set(tc.key, nil, time.Now())
				if ok != tc.want {
					t.Fatalf("Set() key %q failed to be set, got = %v, want = %v", tc.key, ok, tc.want)
				}
			}
			if len(tc.values) >= 1 {
				for _, v := range tc.values {
					time.Sleep(10 * time.Millisecond)
					ok := l.Set(tc.key, v, time.Now())
					if ok != tc.want {
						t.Fatalf("Set() key %q failed to be set, got = %v, want = %v", tc.key, ok, tc.want)
					}
				}
			}
			if _, ok := l.timemap[tc.key]; ok != tc.want {
				t.Fatalf("Set() key %q should exist, got = %v, want = %v", tc.key, ok, true)
			}
		})
	}
}

func TestTimeMapGet(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		want     interface{}
		wantSucc bool
	}{
		{
			name:     "successfully getting a value",
			values:   []string{"Pavlos"},
			want:     "Pavlos",
			wantSucc: true,
		},
		{
			name:     "successfully getting most recent value from multiple values",
			values:   []string{"Pavlos", "Yannis", "Theo"},
			want:     "Theo",
			wantSucc: true,
		},
		{
			name:     "failing getting a key",
			values:   []string{},
			want:     nil,
			wantSucc: false,
		},
		{
			name:     "successfully getting a nil value from a key",
			values:   nil,
			want:     nil,
			wantSucc: true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := New()
			if tc.values == nil {
				l.Set("name", nil, time.Now())
			}
			if len(tc.values) >= 1 {
				for _, v := range tc.values {
					time.Sleep(10 * time.Millisecond)
					l.Set("name", v, time.Now())
				}
			}
			got, ok := l.Get("name", time.Now())
			if ok != tc.wantSucc {
				t.Fatalf("Get() should have a value set got = %v, want = %v", ok, tc.wantSucc)
			}
			if got != tc.want {
				t.Fatalf("Get() should have matched the value in key 'name' got = %s want = %s", got, tc.want)
			}
		})
	}
}

func TestTimeMapRemove(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		values []string
		want   bool
	}{
		{
			name:   "successfully removing key with single value",
			key:    "name",
			values: []string{"Pavlos"},
			want:   true,
		},
		{
			name:   "successfully removing a key with multiple values",
			key:    "name",
			values: []string{"Pavlos", "Yannis", "Theo"},
			want:   true,
		},
		{
			name:   "failing removing a key",
			key:    "noname",
			values: []string{},
			want:   false,
		},
		{
			name:   "failing removing an empty key",
			key:    "",
			values: []string{},
			want:   false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := New()
			if len(tc.values) >= 1 {
				for _, v := range tc.values {
					time.Sleep(10 * time.Millisecond)
					l.Set(tc.key, v, time.Now())
				}
			}
			ok := l.Remove(tc.key)
			if ok != tc.want {
				t.Fatalf("Remove() should have not failed to remove the key %q, got = %v, want = %v", tc.key, ok, tc.want)
			}
			if ok := l.Contains(tc.key); ok {
				t.Fatalf("Remove() should have removed the key %q already, got = %v, want = %v", tc.key, ok, tc.want)
			}
		})
	}
}

func TestTimeMapContains(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		values []string
		want   bool
	}{
		{
			name:   "successfully has a key",
			key:    "name",
			values: []string{"Pavlos"},
			want:   true,
		},
		{
			name:   "failing checking for a key",
			key:    "noname",
			values: []string{},
			want:   false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := New()
			if len(tc.values) >= 1 {
				for _, v := range tc.values {
					time.Sleep(10 * time.Millisecond)
					l.Set(tc.key, v, time.Now())
				}
			}
			ok := l.Contains(tc.key)
			if ok != tc.want {
				t.Fatalf("Contains() should have keys set got = %v, want = %v", ok, tc.want)
			}
		})
	}
}
func TestTimeMapKeys(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		want     []string
		wantKeys bool
	}{
		{
			name:     "successfully getting keys from the timemap",
			want:     []string{"name", "country"},
			wantKeys: true,
		},
		{
			name:     "failing getting keys",
			want:     []string{},
			wantKeys: false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := New()
			if tc.wantKeys {
				l.Set("name", "Pavlos", time.Now())
				time.Sleep(10 * time.Millisecond)
				l.Set("country", "Greece", time.Now())
			}
			got := l.Keys()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Keys() returned unexpected results (-got +want):\n%s", diff)
			}
		})
	}
}
