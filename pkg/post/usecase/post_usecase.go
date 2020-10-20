package usecase

import (
	"context"
	"time"

	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type postUsecase struct {
	postRepo       domain.PostRepository
	authorRepo     domain.AuthorRepository
	contextTimeout time.Duration
}

// NewPostUsecase will create new an postUsecase object representation of domain.PostUsecase interface
func NewPostUsecase(pr domain.PostRepository, ar domain.AuthorRepository, timeout time.Duration) domain.PostUsecase {
	return &postUsecase{
		postRepo:       pr,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

func (p *postUsecase) fillAuthorDetails(c context.Context, data []domain.Post) ([]domain.Post, error) {
	g, ctx := errgroup.WithContext(c)

	// Get author's id
	mapAuthors := map[int64]domain.Author{}

	for _, post := range data {
		mapAuthors[post.Author.ID] = domain.Author{}
	}

	// channel author
	chanAuthor := make(chan domain.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := p.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}

			chanAuthor <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != (domain.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data to post's data
	for i, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[i].Author = a
		}
	}

	return data, nil
}

func (p *postUsecase) Store(c context.Context, e *domain.Post) error {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	// Check if post already posted
	existedPost, _ := p.GetByTitle(ctx, e.Title)
	if existedPost != (domain.Post{}) {
		return domain.ErrConflict
	}

	err := p.postRepo.Store(ctx, e)
	return err
}

func (p *postUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Post, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	res, nextCursor, err = p.postRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = p.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}

	return
}

func (p *postUsecase) GetByID(c context.Context, id int64) (res domain.Post, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	res, err = p.postRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := p.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Post{}, err
	}

	res.Author = resAuthor
	return
}

func (p *postUsecase) GetByTitle(c context.Context, title string) (res domain.Post, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	res, err = p.postRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := p.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return
	}

	res.Author = resAuthor
	return
}

func (p *postUsecase) Update(c context.Context, e *domain.Post) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	e.UpdatedAt = time.Now()
	return p.postRepo.Update(ctx, e)
}

func (p *postUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	// check existedpost
	existedPost, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	if existedPost == (domain.Post{}) {
		return domain.ErrNotFound
	}

	return p.postRepo.Delete(ctx, id)
}
