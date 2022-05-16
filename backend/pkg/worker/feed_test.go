package worker

import (
	"errors"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/spezifisch/rueder3/backend/pkg/worker/scheduler"
	"github.com/stretchr/testify/assert"
)

func readFeed(t *testing.T, filename string) *gofeed.Feed {
	file, err := os.Open(filename)
	assert.NoError(t, err)
	defer file.Close()

	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	assert.NoError(t, err)
	return feed
}

func readRaumzeitFeed(t *testing.T) *gofeed.Feed {
	feed := readFeed(t, "../../test/data/raumzeit-podcast-oga.xml")
	assert.Equal(t, "Raumzeit", feed.Title)
	return feed
}

func readFefeFeed(t *testing.T) *gofeed.Feed {
	feed := readFeed(t, "../../test/data/fefe-html.xml")
	assert.Equal(t, "Fefes Blog", feed.Title)
	return feed
}

func readGolemFeed(t *testing.T) *gofeed.Feed {
	feed := readFeed(t, "../../test/data/golem.xml")
	assert.Equal(t, "Golem.de", feed.Title)
	return feed
}

type mockRepository struct {
	t *testing.T

	// mocking options
	allArticlesNew            bool
	failCheckExistingArticles bool
	failAddArticle            bool

	// repository data
	addedArticles []scheduler.Article
}

func (m *mockRepository) Feeds() ([]scheduler.Feed, error) {
	return nil, errors.New("not implemented")
}
func (m *mockRepository) GetFeed(feedID uuid.UUID) (scheduler.Feed, error) {
	return scheduler.Feed{}, errors.New("not implemented")
}
func (m *mockRepository) RunFeedChangeListener(addedFeeds chan<- uuid.UUID, needRehash chan<- bool) (err error) {
	return errors.New("not implemented")
}
func (m *mockRepository) UpdateFeedInfo(feedID uuid.UUID, updatedFeed *scheduler.Feed) (err error) {
	return errors.New("not implemented")
}
func (m *mockRepository) CheckExistingArticles(feedID uuid.UUID, articleGUIDs []string) (exists []bool, err error) {
	if m.failCheckExistingArticles {
		return nil, errors.New("mock failCheckExistingArticles")
	}

	exists = make([]bool, len(articleGUIDs))
	if m.allArticlesNew {
		return
	}

	// all articles existing (old)
	for i := 0; i < len(exists); i++ {
		exists[i] = true
	}
	return
}
func (m *mockRepository) AddArticle(feedID uuid.UUID, article *scheduler.Article) error {
	assert.NotNil(m.t, article, "foo")
	if m.failAddArticle {
		return errors.New("mock failAddArticle")
	}

	m.addedArticles = append(m.addedArticles, *article)
	return nil
}

func TestFeedWorkerPool_processArticles(t *testing.T) {
	f := &scheduler.Feed{}
	feedRaumzeit := readRaumzeitFeed(t)
	feedFefe := readFefeFeed(t)
	feedGolem := readGolemFeed(t)

	type fields struct {
		config     FeedWorkerConfig
		repository scheduler.Repository
	}
	type args struct {
		f    *scheduler.Feed
		feed *gofeed.Feed
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		check  func(*mockRepository)
	}{
		{
			name: "existing check fails",
			fields: fields{
				repository: &mockRepository{
					t:                         t,
					failCheckExistingArticles: true,
				},
			},
			args: args{
				f:    f,
				feed: feedRaumzeit,
			},
		},
		{
			name: "all articles existing",
			fields: fields{
				repository: &mockRepository{
					t:              t,
					allArticlesNew: false,
				},
			},
			args: args{
				f:    f,
				feed: feedRaumzeit,
			},
		},
		{
			name: "all articles new but inserting fails",
			fields: fields{
				repository: &mockRepository{
					t:              t,
					allArticlesNew: true,
					failAddArticle: true,
				},
			},
			args: args{
				f:    f,
				feed: feedRaumzeit,
			},
		},
		{
			name: "all articles new",
			fields: fields{
				repository: &mockRepository{
					t:              t,
					allArticlesNew: true,
					failAddArticle: false,
					addedArticles:  make([]scheduler.Article, 0),
				},
			},
			args: args{
				f:    f,
				feed: feedRaumzeit,
			},
			check: func(m *mockRepository) {
				assert.Equal(t, 5, len(m.addedArticles))

				// check dates
				prevArticle := &m.addedArticles[0]
				for i := 1; i < len(m.addedArticles); i++ {
					thisArticle := &m.addedArticles[i]
					assert.Greater(t, thisArticle.Time.Unix(), prevArticle.Time.Unix(), "article date ordering not chronological")

					prevArticle = thisArticle
				}
			},
		},
		{
			name: "feed without dates",
			fields: fields{
				repository: &mockRepository{
					t:              t,
					allArticlesNew: true,
					failAddArticle: false,
					addedArticles:  make([]scheduler.Article, 0),
				},
			},
			args: args{
				f:    f,
				feed: feedFefe,
			},
			check: func(m *mockRepository) {
				assert.Equal(t, 10, len(m.addedArticles))

				// check order
				wantedGUIDOrder := []string{
					"https://blog.fefe.de/?ts=9e60a8a7",
					"https://blog.fefe.de/?ts=9e602dd0",
					"https://blog.fefe.de/?ts=9e617d6c",
				}
				for i := 0; i < len(wantedGUIDOrder); i++ {
					thisArticle := &m.addedArticles[i]
					assert.Equal(t, wantedGUIDOrder[i], thisArticle.SiteGUID, "articles not in wanted order")
				}
			},
		},
		{
			name: "feed with CEST timestamps",
			fields: fields{
				repository: &mockRepository{
					t:              t,
					allArticlesNew: true,
					failAddArticle: false,
					addedArticles:  make([]scheduler.Article, 0),
				},
			},
			args: args{
				f:    f,
				feed: feedGolem,
			},
			check: func(m *mockRepository) {
				assert.Equal(t, 3, len(m.addedArticles))

				// check for correct timezone handling
				wantedTimeHours := []int{
					9,
					9,
					10,
				}
				for i := 0; i < len(wantedTimeHours); i++ {
					thisArticle := &m.addedArticles[i]
					assert.Equal(t, wantedTimeHours[i], thisArticle.Time.UTC().Hour(), "article timestamps not parsed correctly")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FeedWorkerPool{
				config:     tt.fields.config,
				repository: tt.fields.repository,
			}
			p.processArticles(tt.args.f, tt.args.feed)
			if tt.check != nil {
				tt.check(p.repository.(*mockRepository))
			}
		})
	}
}
