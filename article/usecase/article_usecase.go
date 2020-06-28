package usecase

import (
	"context"
	"fmt"
	"time"

	"demo-echo/model"
	models "demo-echo/model"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type articleUsecase struct {
	articleRepo     models.ArticleRepository
	authorRepo      models.AuthorRepository
	contextTimeoput time.Duration
}

// NewArticleUsecase func
func NewArticleUsecase(artRepo models.ArticleRepository, auRepo model.AuthorRepository, timeout time.Duration) model.ArticleUsecase {
	return &articleUsecase{
		articleRepo:     artRepo,
		authorRepo:      auRepo,
		contextTimeoput: timeout,
	}
}

func (a *articleUsecase) fillAuthorDetails(c context.Context, data []models.Article) ([]model.Article, error) {
	g, ctx := errgroup.WithContext(c)
	mapAuthors := map[int64]models.Author{}
	for _, article := range data {
		mapAuthors[article.Author.ID] = models.Author{}
	}

	chanAuthor := make(chan models.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
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
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != (models.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for index, item := range data {
		if au, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = au
		}
	}
	return data, nil
}

func (a *articleUsecase) Fetch(c context.Context, cursor string, num int64) (res []model.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}
	ctx, cancel := context.WithTimeout(c, a.contextTimeoput)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *articleUsecase) GetByID(c context.Context, id int64) (res model.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeoput)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return models.Article{}, err
	}
	res.Author = resAuthor
	return
}

func (a *articleUsecase) Update(c context.Context, ar *model.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeoput)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleUsecase) GetByTitle(c context.Context, title string) (res model.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeoput)
	defer cancel()
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return model.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *articleUsecase) Store(c context.Context, ar *models.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeoput)
	defer cancel()
	existArticle, _ := a.GetByTitle(ctx, ar.Title)
	if existArticle != (models.Article{}) {
		return model.ErrConflict
	}
	ar.SetCreatedAt()
	ar.SetUpdatedAt()
	fmt.Println("ar:", ar)
	err = a.articleRepo.Store(ctx, ar)
	if err != nil {
		return model.ErrInternalServerError
	}
	return
}

func (a *articleUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeoput)
	defer cancel()
	existArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existArticle == (models.Article{}) {
		return models.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
