package main

import "testing"

func TestHashes(t *testing.T) {

	repo := GitRepository{"test-data/samplerepo"}
	hs := repo.Hashes()

	if hs == nil || len(hs) == 0 {
		t.Error("No hashes read.")
	}

	found := false

	for _, h := range(hs) {
		if h == "620dae1d8c03c241c6ae27c24e76548dbf5b37a1" {
			found = true
			break
		}
	}

	if found == false {
		t.Error("Unable to find hash 620dae1.")
	}

}

func TestObjects(t *testing.T) {

	repo := GitRepository{"test-data/samplerepo"}
	bs, ts, cs := repo.Objects()

	if ts == nil || len(ts) == 0 {
		t.Error("No trees read.")
	}

	if bs == nil || len(bs) == 0 {
		t.Error("No blobs read.")
	}

	if cs == nil || len(cs) == 0 {
		t.Error("No commits read.")
	}

	for _, c := range(cs) {

		if c.hash == "" {
			t.Errorf("Commit hash was blank.")
		}

		if c.author == "" {
			t.Errorf("Commit author was blank.")
		}

		if c.committer == "" {
			t.Errorf("Commit committer was blank.")
		}

		if c.date == "" {
			t.Errorf("Commit date was blank.")
		}

		if c.tree == "" {
			t.Errorf("Commit tree was blank.")
		}

	}

}
