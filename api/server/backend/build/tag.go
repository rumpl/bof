package build

import (
	"fmt"
	"io"

	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
	"github.com/rumpl/bof/image"
)

// tagImages creates image tags for the imageID.
func tagImages(ic ImageComponent, stdout io.Writer, imageID image.ID, repoAndTags []reference.Named) error {
	for _, rt := range repoAndTags {
		if err := ic.TagImageWithReference(imageID, rt); err != nil {
			return err
		}
		_, _ = fmt.Fprintln(stdout, "Successfully tagged", reference.FamiliarString(rt))
	}
	return nil
}

// sanitizeRepoAndTags parses the raw "t" parameter received from the client
// to a slice of repoAndTag.
// It also validates each repoName and tag.
func sanitizeRepoAndTags(names []string) ([]reference.Named, error) {
	var (
		repoAndTags []reference.Named
		// This map is used for deduplicating the "-t" parameter.
		uniqNames = make(map[string]struct{})
	)
	for _, repo := range names {
		if repo == "" {
			continue
		}

		ref, err := reference.ParseNormalizedNamed(repo)
		if err != nil {
			return nil, err
		}

		if _, isCanonical := ref.(reference.Canonical); isCanonical {
			return nil, errors.New("build tag cannot contain a digest")
		}

		ref = reference.TagNameOnly(ref)

		nameWithTag := ref.String()

		if _, exists := uniqNames[nameWithTag]; !exists {
			uniqNames[nameWithTag] = struct{}{}
			repoAndTags = append(repoAndTags, ref)
		}
	}
	return repoAndTags, nil
}
