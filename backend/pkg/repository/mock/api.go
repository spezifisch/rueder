package mock

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
)

// Repository gives some dummy data for testing
// it implements a controller.FeedRepository
type Repository struct {
	feedCount      int
	feedID         uuid.UUID
	runningCounter int
	seqCounter     int
}

// NewMockRepository returns a FeedRepository that can be used in place of a real database
func NewMockRepository() *Repository {
	return &Repository{
		seqCounter: 100000000000,
	}
}

// Folders returns mock folders
func (r *Repository) Folders(claims *helpers.AuthClaims) ([]controller.Folder, error) {
	r.runningCounter = 0
	feedCounts := []int{1, 2, 8, 3, 1}
	folderCount := 5
	folders := make([]controller.Folder, folderCount)
	for i := 0; i < folderCount; i++ {
		r.feedCount = feedCounts[i]
		folders[i], _ = r.GetFolder(1 + i)
	}

	return folders, nil
}

// Labels returns mock labels
func (r *Repository) Labels() ([]controller.Label, error) {
	count := 4
	labels := make([]controller.Label, count)
	for i := 0; i < count; i++ {
		labels[i], _ = r.GetLabel(1 + i)
	}
	return labels, nil
}

// GetArticle returns a mock article with the given id
func (r *Repository) GetArticle(id uuid.UUID) (controller.Article, error) {
	return controller.Article{
		ID:        getMockUUID(),
		FeedTitle: "A Feed Online News",

		Title:        fmt.Sprintf("Mining %d Bitcoin On The Nintendo Game Boy", r.runningCounter),
		Time:         time.Now(),
		Link:         "about:blank",
		LinkComments: "about:blank",

		Thumbnail:  "/favicon.png",
		Image:      "/favicon.png",
		ImageTitle: "An IMAGE TITLE",

		Content: controller.ArticleContent{
			Text: `<p>Lorem ipsum dolor <a href="about:blank">sit amet</a>, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Dolor morbi non arcu risus quis varius quam quisque id. Sodales ut etiam sit amet nisl purus. Viverra ipsum nunc aliquet bibendum enim facilisis. Et malesuada fames ac turpis egestas integer eget aliquet. Quis risus sed vulputate odio ut enim blandit volutpat. Pellentesque dignissim enim sit amet venenatis urna. Et egestas quis ipsum suspendisse ultrices gravida dictum fusce. Iaculis eu non diam phasellus. Luctus accumsan tortor posuere ac. Praesent semper feugiat nibh sed pulvinar proin gravida hendrerit. Faucibus nisl tincidunt eget nullam non nisi est sit amet. Sed velit dignissim sodales ut eu sem integer.</p>
				<p>Varius vel pharetra vel turpis nunc eget lorem dolor. Sit amet consectetur adipiscing elit. Orci eu lobortis elementum nibh. Sodales neque sodales ut etiam sit amet nisl purus in. Est ante in nibh mauris. In hendrerit gravida rutrum quisque non. Urna molestie at elementum eu facilisis sed. Tortor posuere ac ut consequat semper. Sodales neque sodales ut etiam sit amet nisl purus in. Amet risus nullam eget felis eget nunc lobortis. Purus semper eget duis at tellus. Ullamcorper sit amet risus nullam eget felis eget nunc lobortis. Velit scelerisque in dictum non consectetur a. Gravida quis blandit turpis cursus. Ultrices tincidunt arcu non sodales neque. Pharetra vel turpis nunc eget lorem. Mus mauris vitae ultricies leo integer malesuada nunc vel risus. Arcu ac tortor dignissim convallis aenean et tortor at risus. Sit amet aliquam id diam maecenas ultricies. Tellus elementum sagittis vitae et leo duis ut diam quam.</p>
				<p>Nunc vel risus commodo viverra maecenas accumsan lacus. Dictum non consectetur a erat nam at lectus urna duis. Augue mauris augue neque gravida in fermentum. Arcu dictum varius duis at consectetur. Donec adipiscing tristique risus nec feugiat in fermentum. Eget aliquet nibh praesent tristique. Leo vel fringilla est ullamcorper eget. Purus gravida quis blandit turpis cursus in hac. Nullam eget felis eget nunc lobortis mattis aliquam faucibus. Non curabitur gravida arcu ac. In hac habitasse platea dictumst vestibulum rhoncus est pellentesque. Quam pellentesque nec nam aliquam sem et tortor consequat id. Cras adipiscing enim eu turpis egestas pretium aenean pharetra. Vestibulum sed arcu non odio euismod lacinia.</p>
				<p>At auctor urna nunc id. Consectetur a erat nam at. Massa placerat duis ultricies lacus sed turpis tincidunt id aliquet. Ac placerat vestibulum lectus mauris ultrices eros in cursus. Cum sociis natoque penatibus et magnis dis parturient. Lacus sed viverra tellus in hac. Semper quis lectus nulla at volutpat diam ut venenatis. Turpis egestas pretium aenean pharetra magna. Tortor vitae purus faucibus ornare suspendisse sed nisi. Sodales neque sodales ut etiam sit. Nisi vitae suscipit tellus mauris a diam. Suscipit adipiscing bibendum est ultricies integer quis. Phasellus vestibulum lorem sed risus ultricies. Lobortis mattis aliquam faucibus purus. Quis risus sed vulputate odio ut enim blandit. Vestibulum sed arcu non odio euismod lacinia. Est sit amet facilisis magna etiam. Orci eu lobortis elementum nibh tellus. <a href="about:blank">Faucibus et molestie ac feugiat</a> sed lectus vestibulum mattis.</p>
				<p>Tempus imperdiet nulla malesuada pellentesque elit eget gravida cum sociis. Commodo viverra maecenas accumsan lacus. Quis blandit turpis cursus in hac. Magnis dis parturient montes nascetur ridiculus mus. Aliquam etiam erat velit scelerisque in dictum non consectetur a. Tortor aliquam nulla facilisi cras. Aliquet lectus proin nibh nisl condimentum id venenatis a. Facilisis mauris sit amet massa vitae tortor condimentum lacinia. In nisl nisi scelerisque eu ultrices vitae auctor. Sit amet consectetur adipiscing elit. Non enim praesent elementum facilisis leo. Ut enim blandit volutpat maecenas volutpat. Faucibus interdum posuere lorem ipsum dolor sit amet consectetur adipiscing. Lorem mollis aliquam ut porttitor leo. Volutpat blandit aliquam etiam erat.</p>
				<p>Nam at lectus urna duis convallis convallis tellus id. Sodales ut etiam sit amet nisl purus in. Faucibus ornare suspendisse sed nisi lacus. Lobortis elementum nibh tellus molestie nunc non blandit massa enim. Semper quis lectus nulla at volutpat diam. Enim sed faucibus turpis in. Semper feugiat nibh sed pulvinar. Orci nulla pellentesque dignissim enim sit amet venenatis. Adipiscing diam donec adipiscing tristique risus. Proin libero nunc consequat interdum varius sit. Sed augue lacus viverra vitae congue eu consequat. Ullamcorper morbi tincidunt ornare massa eget egestas. Quam viverra orci sagittis eu. Non sodales neque sodales ut etiam sit amet nisl. Et odio pellentesque diam volutpat commodo sed. Ornare aenean euismod elementum nisi quis eleifend quam adipiscing vitae. Et ultrices neque ornare aenean euismod elementum nisi quis eleifend. Lacus laoreet non curabitur gravida arcu ac tortor dignissim.</p>
				<p>Aliquet eget sit amet tellus cras adipiscing enim. Suspendisse in est ante in nibh mauris. Tincidunt ornare massa eget egestas purus. Ut enim blandit volutpat maecenas volutpat blandit aliquam etiam erat. Parturient montes nascetur ridiculus mus mauris. Sagittis eu volutpat odio facilisis mauris sit. Ultrices mi tempus imperdiet nulla malesuada. Nisl vel pretium lectus quam id. Morbi tempus iaculis urna id volutpat lacus laoreet. Varius sit amet mattis vulputate enim nulla aliquet porttitor lacus. Pellentesque dignissim enim sit amet venenatis urna cursus eget. Neque ornare aenean euismod elementum nisi quis eleifend. Rhoncus mattis rhoncus urna neque. Vivamus arcu felis bibendum ut tristique et egestas quis ipsum.</p>
				<p>Lorem mollis aliquam ut porttitor leo a. Vulputate odio ut enim blandit. Vulputate mi sit amet mauris commodo quis imperdiet massa tincidunt. Blandit aliquam etiam erat velit scelerisque in dictum non. Consequat id porta nibh venenatis cras sed felis eget velit. Adipiscing commodo elit at imperdiet dui accumsan sit amet nulla. Euismod elementum nisi quis eleifend quam adipiscing vitae proin. Adipiscing vitae proin sagittis nisl rhoncus mattis. Faucibus ornare suspendisse sed nisi lacus. Dui id ornare arcu odio ut. Nunc pulvinar sapien et ligula ullamcorper malesuada proin libero nunc. Facilisis leo vel fringilla est ullamcorper. Sit amet porttitor eget dolor morbi. Pharetra massa massa ultricies mi quis hendrerit dolor magna. Auctor elit sed vulputate mi sit amet mauris commodo quis.</p>
				<p>Nullam ac tortor vitae purus faucibus ornare. Risus nec feugiat in fermentum posuere urna nec tincidunt praesent. Id volutpat lacus laoreet non curabitur gravida arcu. Suspendisse sed nisi lacus sed viverra. At tellus at urna condimentum mattis. Et molestie ac feugiat sed lectus vestibulum mattis ullamcorper. Consectetur lorem donec massa sapien faucibus et molestie ac feugiat. Vel turpis nunc eget lorem. Viverra aliquet eget sit amet. Volutpat lacus laoreet non curabitur gravida. Habitasse platea dictumst quisque sagittis purus sit amet volutpat consequat. Pulvinar pellentesque habitant morbi tristique senectus. Sed faucibus turpis in eu. Euismod lacinia at quis risus sed vulputate odio ut. Tellus in hac habitasse platea dictumst vestibulum rhoncus.</p>
				<p>Habitant morbi tristique senectus et netus et. Arcu cursus euismod quis viverra nibh. Massa enim nec dui nunc mattis enim ut tellus. Id venenatis a condimentum vitae sapien pellentesque. Dictum varius duis at consectetur lorem donec massa. Interdum velit laoreet id donec. Vitae et leo duis ut diam quam nulla. Urna neque viverra justo nec. Sem integer vitae justo eget. Cras semper auctor neque vitae tempus quam pellentesque.</p>`,
		},
	}, nil
}

// GetArticlePreview returns a mock article with the given id
func (r *Repository) GetArticlePreview(id int) (controller.ArticlePreview, error) {
	return controller.ArticlePreview{
		ID:        getMockUUID(),
		Seq:       r.getSeq(),
		Title:     fmt.Sprintf("A preview for article %d on feed %d", r.runningCounter+2, 23),
		Time:      time.Now(),
		FeedTitle: fmt.Sprintf("Feed Online News %d", r.runningCounter),
		FeedIcon:  "/favicon.png",
		Teaser:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}, nil
}

// GetArticles returns mock articles
func (r *Repository) GetArticles(feedID uuid.UUID, limit int, offset int) ([]controller.ArticlePreview, error) {
	r.feedID = feedID
	count := 20
	articles := make([]controller.ArticlePreview, count)
	for i := 0; i < count; i++ {
		articles[i], _ = r.GetArticlePreview(1 + i)
	}
	return articles, nil
}

// GetFeed returns a mock feed
func (r *Repository) GetFeed(id uuid.UUID) (controller.Feed, error) {
	unreadCount := 5
	if r.runningCounter%2 == 0 {
		unreadCount = 0
	} else if r.runningCounter > 5 {
		unreadCount *= 2
	}

	r.runningCounter++

	return controller.Feed{
		ID:           getMockUUID(),
		Title:        fmt.Sprintf("Feed Number %d", r.runningCounter),
		Icon:         "/favicon.png",
		ArticleCount: unreadCount,
	}, nil
}

// Feeds returns some feeds
func (r *Repository) Feeds() (ret []controller.Feed, err error) {
	ret = make([]controller.Feed, 10)
	for i := 0; i < len(ret); i++ {
		ret[i], _ = r.GetFeed(getMockUUID())
	}
	return
}

// GetFolder returns a single mock folder
func (r *Repository) GetFolder(id int) (controller.Folder, error) {
	feeds := make([]controller.Feed, r.feedCount)
	for i := 0; i < r.feedCount; i++ {
		feeds[i], _ = r.GetFeed(getMockUUID())
	}

	return controller.Folder{
		ID:    getMockUUID(),
		Title: fmt.Sprintf("Folder %d", id),
		Feeds: feeds,
	}, nil
}

// GetLabel returns a single mock label
func (*Repository) GetLabel(id int) (controller.Label, error) {
	title := fmt.Sprintf("Label %d", id)
	color := []string{"red", "green", "blue"}[id%3]
	return controller.Label{
		ID:    getMockUUID(),
		Title: title,
		Color: color,
	}, nil
}

// AddFeed does nothing
func (*Repository) AddFeed(url string) (feedID uuid.UUID, err error) {
	err = errors.New("not implemented")
	return
}

// GetFeedByURL does nothing
func (*Repository) GetFeedByURL(url string) (ret controller.Feed, err error) {
	err = errors.New("not implemented")
	return
}

// ChangeFolders does nothing
func (*Repository) ChangeFolders(claims *helpers.AuthClaims, folders []controller.Folder) (err error) {
	err = errors.New("not implemented")
	return
}

func getMockUUID() uuid.UUID {
	id, _ := uuid.NewGen().NewV4()
	return id
}

func (r *Repository) getSeq() int {
	r.seqCounter--
	return r.seqCounter
}
