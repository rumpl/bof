package urlutil

import "testing"

var (
	gitUrls = []string{
		"git://github.com/rumpl/bof",
		"git@github.com:rumpl/bof.git",
		"git@bitbucket.org:atlassianlabs/atlassian-docker.git",
		"https://github.com/rumpl/bof.git",
		"http://github.com/rumpl/bof.git",
		"http://github.com/rumpl/bof.git#branch",
		"http://github.com/rumpl/bof.git#:dir",
	}
	incompleteGitUrls = []string{
		"github.com/rumpl/bof",
	}
	invalidGitUrls = []string{
		"http://github.com/rumpl/bof.git:#branch",
	}
)

func TestIsGIT(t *testing.T) {
	for _, url := range gitUrls {
		if !IsGitURL(url) {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}

	for _, url := range incompleteGitUrls {
		if !IsGitURL(url) {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}

	for _, url := range invalidGitUrls {
		if IsGitURL(url) {
			t.Fatalf("%q should not be detected as valid Git prefix", url)
		}
	}
}
