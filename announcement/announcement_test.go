package announcement

import (
	"testing"

	"github.com/UnipayFI/go-gate/v4/internal/testutil"
)

func TestAnnouncement(t *testing.T) {
	t.Run("ListArticles", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		got, err := c.NewListArticlesService().SetPage(1).SetSize(5).Do(cx)
		if err != nil {
			t.Fatalf("list articles: %v", err)
		}
		if got.Code != 0 {
			t.Fatalf("list articles: code=%d message=%s", got.Code, got.Message)
		}
		if len(got.Data.List) == 0 {
			t.Fatal("no articles returned")
		}
		t.Logf("total=%d first=%+v", got.Data.Total, got.Data.List[0])
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/ann/list_article",
			map[string]any{"page": "1", "size": "5"}, false)
		testutil.AssertCovers(t, "ann/list_article", raw, got)
	})

	t.Run("ListArticlesSearch", func(t *testing.T) {
		c := testPublicClient()
		cx := testutil.Ctx(t)

		got, err := c.NewListArticlesService().SetTitleQuery("BTC").SetSize(3).Do(cx)
		if err != nil {
			t.Fatalf("search articles: %v", err)
		}
		t.Logf("total=%d list=%d", got.Data.Total, len(got.Data.List))
		raw := testutil.FetchRawPost(t, c, cx, "/api/v4/ann/list_article",
			map[string]any{"title_query": "BTC", "size": "3"}, false)
		testutil.AssertCovers(t, "ann/list_article(title_query)", raw, got)
	})
}
