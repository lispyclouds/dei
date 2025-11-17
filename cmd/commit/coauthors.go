package commit

import (
	"bytes"
	json "encoding/json/v2"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/lispyclouds/dei/pkg"
)

const coAuthorsKey = "dei.commit.coAuthors"

type CoAuthors = map[string]map[string]string

func loadCoAuthors(cache *pkg.Cache) (CoAuthors, error) {
	data, err := cache.Get(coAuthorsKey)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return CoAuthors{}, nil
	}

	var coAuthors CoAuthors
	if err = json.UnmarshalRead(bytes.NewReader(data), &coAuthors); err != nil {
		return nil, err
	}

	return coAuthors, nil
}

func manageCoAuthor(cache *pkg.Cache, name, email, op string) error {
	coAuthors, err := loadCoAuthors(cache)
	if err != nil {
		return err
	}

	switch op {
	case "add":
		info, ok := coAuthors[email]
		if !ok {
			info = make(map[string]string)
		}

		info["name"] = name
		coAuthors[email] = info
	case "remove":
		delete(coAuthors, email)
	}

	buffer := bytes.NewBuffer([]byte{})
	if err = json.MarshalWrite(buffer, &coAuthors); err != nil {
		return err
	}

	return cache.Put(coAuthorsKey, buffer.Bytes())
}

func listCoAuthors(cache *pkg.Cache) error {
	coAuthors, err := loadCoAuthors(cache)
	if err != nil {
		return err
	}

	if len(coAuthors) == 0 {
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	fmt.Fprintln(w, "Name\tEmail")
	for email, info := range coAuthors {
		fmt.Fprintf(w, "%s\t%s\n", info["name"], email)
	}
	w.Flush()

	return nil
}
