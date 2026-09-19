package announcement

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-gate/v4/request"
	"github.com/shopspring/decimal"
)

// ListArticlesService -- POST /api/v4/ann/list_article (public)
//
// Lists announcement articles, paginated and filtered by title, tags, category,
// language and time. The filters travel in the JSON body; page, size, timer,
// cate_level and sub_website_id are sent as strings.
type ListArticlesService struct {
	c    *AnnouncementClient
	body map[string]any
}

func (c *AnnouncementClient) NewListArticlesService() *ListArticlesService {
	return &ListArticlesService{c: c, body: map[string]any{}}
}

// SetTitleQuery filters by article title.
func (s *ListArticlesService) SetTitleQuery(titleQuery string) *ListArticlesService {
	s.body["title_query"] = titleQuery
	return s
}

// SetPage selects the page number.
func (s *ListArticlesService) SetPage(page int) *ListArticlesService {
	s.body["page"] = strconv.Itoa(page)
	return s
}

// SetSize caps the number of articles per page.
func (s *ListArticlesService) SetSize(size int) *ListArticlesService {
	s.body["size"] = strconv.Itoa(size)
	return s
}

// SetTags filters by article tags.
func (s *ListArticlesService) SetTags(tags string) *ListArticlesService {
	s.body["tags"] = tags
	return s
}

// SetTimer limits the result to announcements from the last N days (e.g. on the
// 10th, 10 covers the 1st through the 10th).
func (s *ListArticlesService) SetTimer(days int) *ListArticlesService {
	s.body["timer"] = strconv.Itoa(days)
	return s
}

// SetCategoryName filters by announcement category name.
func (s *ListArticlesService) SetCategoryName(categoryName string) *ListArticlesService {
	s.body["cate_name"] = categoryName
	return s
}

// SetCategoryLevel selects the category level (1 or 2).
func (s *ListArticlesService) SetCategoryLevel(categoryLevel int) *ListArticlesService {
	s.body["cate_level"] = strconv.Itoa(categoryLevel)
	return s
}

// SetSubWebsiteID selects the subsite ("0" = main site, the default; "177" =
// Turkey site).
func (s *ListArticlesService) SetSubWebsiteID(subWebsiteID string) *ListArticlesService {
	s.body["sub_website_id"] = subWebsiteID
	return s
}

// SetPinned controls whether pinned articles are included (1 = include, the
// default; 0 = exclude).
func (s *ListArticlesService) SetPinned(pinned int) *ListArticlesService {
	s.body["pinned"] = pinned
	return s
}

// SetUpdateAfter limits the result to articles updated after this time.
func (s *ListArticlesService) SetUpdateAfter(updateAfter time.Time) *ListArticlesService {
	s.body["update_after"] = updateAfter.Unix()
	return s
}

// SetLang selects the language code (e.g. "cn", "en").
func (s *ListArticlesService) SetLang(lang string) *ListArticlesService {
	s.body["lang"] = lang
	return s
}

// SetFilterEmptyContent controls whether articles with empty content in the
// current language are excluded (1 = exclude, the default; 0 = keep).
func (s *ListArticlesService) SetFilterEmptyContent(filterEmptyContent int) *ListArticlesService {
	s.body["filter_empty_content"] = filterEmptyContent
	return s
}

func (s *ListArticlesService) Do(ctx context.Context) (*AnnouncementArticlesResponse, error) {
	req := request.Post(ctx, s.c, "/api/v4/ann/list_article", s.body)
	return request.Do[AnnouncementArticlesResponse](req)
}

// AnnouncementArticlesResponse is the {code,message,data,version} business
// envelope of the article-list query; Code 0 means success.
type AnnouncementArticlesResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Version string `json:"version"`
	Data    struct {
		Total int64                 `json:"total"`
		List  []AnnouncementArticle `json:"list"`
	} `json:"data"`
}

// AnnouncementArticle is a single announcement article. created / updated are
// display strings (a date like "2026-08-03 UTC" or a relative age like
// "1 days"); created_t / updated_t / release_timestamp carry the Unix-second
// times. url is the site-relative article path. cate and content_lang_list are
// documented but not currently returned by the live API; _score is the search
// relevance score.
type AnnouncementArticle struct {
	ID                  int64                         `json:"id"`
	Title               string                        `json:"title"`
	Brief               string                        `json:"brief"`
	Created             string                        `json:"created"`
	Updated             string                        `json:"updated"`
	ReleaseTime         string                        `json:"release_time"`
	ReleaseTimestamp    time.Time                     `json:"release_timestamp,string,format:unix"`
	Views               int64                         `json:"views"`
	Author              string                        `json:"author"`
	AuthorID            int64                         `json:"author_id"`
	Tags                string                        `json:"tags"`
	Category            string                        `json:"cate"`
	IsTop               int                           `json:"is_top"`
	CategoryID          int64                         `json:"cate_id"`
	Source              string                        `json:"source"`
	ContentLanguageList []AnnouncementArticleLanguage `json:"content_lang_list"`
	CreatedTime         time.Time                     `json:"created_t,format:unix"`
	UpdatedTime         time.Time                     `json:"updated_t,format:unix"`
	Votes1              int64                         `json:"votes1"`
	URL                 string                        `json:"url"`
	HighlightTitle      string                        `json:"highlight_title"`
	HighlightBrief      string                        `json:"highlight_brief"`
	Score               decimal.Decimal               `json:"_score"`
}

// AnnouncementArticleLanguage is one language an article is available in.
type AnnouncementArticleLanguage struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
