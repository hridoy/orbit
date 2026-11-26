package database

import "testing"

func TestCatalogEntriesAreComplete(t *testing.T) {
	for _, c := range Catalog {
		if c.Engine == "" || c.Version == "" || c.Bin == "" {
			t.Fatalf("incomplete entry: %+v", c)
		}
		switch c.Method {
		case MethodApt:
			if c.Pkg == "" {
				t.Fatalf("%s %s: apt entry needs a package", c.Engine, c.Version)
			}
		case MethodTarball, MethodSource:
			if c.URL == "" {
				t.Fatalf("%s %s: archive entry needs a URL", c.Engine, c.Version)
			}
		default:
			t.Fatalf("%s %s: unknown method %q", c.Engine, c.Version, c.Method)
		}
	}
}

func TestCatalogRejectsUnknownVersion(t *testing.T) {
	if _, ok := catalogFor("postgres", "9.6"); ok {
		t.Fatal("unknown version found in catalog")
	}
	if _, ok := catalogFor("postgres", "17"); !ok {
		t.Fatal("postgres 17 missing from catalog")
	}
}
