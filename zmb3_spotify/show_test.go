package spotify

import "github.com/jkravitz/mytrace"

import (
	"context"
	"net/http"
	"testing"
)

func TestGetShow(t *testing.T) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	c, s := testClientFile(http.StatusOK, "test_data/get_show.txt")
	defer s.Close()

	r, err := c.GetShow(context.Background(), "1234")
	if err != nil {
		t.Fatal(err)
	}
	if r.SimpleShow.Name != "Uncommon Core" {
		t.Error("Invalid data:", r.Name)
	}
	if len(r.Episodes.Episodes) != 25 {
		t.Error("Invalid data", len(r.Episodes.Episodes))
	}
}

func TestGetShowEpisodes(t *testing.T) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	c, s := testClientFile(http.StatusOK, "test_data/get_show_episodes.txt")
	defer s.Close()

	r, err := c.GetShowEpisodes(context.Background(), "1234")
	if err != nil {
		t.Fatal(err)
	}
	if r.Total != 25 {
		t.Error("Invalid data:", r.Total)
	}
	if r.Offset != 0 {
		t.Error("Invalid data:", r.Offset)
	}
	if len(r.Episodes) != 25 {
		t.Error("Invalid data", len(r.Episodes))
	}
}
