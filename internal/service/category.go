package service

import (
	"context"
	"io"
	"sync"

	"github.com/mariofelesdossantosjunior/sample-grpc/internal/database"
	"github.com/mariofelesdossantosjunior/sample-grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB  database.Category
	subscribers map[chan *pb.Category]struct{}
	mu          sync.Mutex
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB:  categoryDB,
		subscribers: make(map[chan *pb.Category]struct{}),
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	go c.broadcast(categoryResponse)

	return categoryResponse, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var categoryResponses []*pb.Category
	for _, category := range categories {
		categoryResponse := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
		categoryResponses = append(categoryResponses, categoryResponse)
	}

	return &pb.CategoryList{Categories: categoryResponses}, nil
}

func (c *CategoryService) ListCategoriesStream(in *pb.Blank, stream pb.CategoryService_ListCategoriesStreamServer) error {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return err
	}

	// Envia categorias existentes, uma por uma
	for _, category := range categories {
		err := stream.Send(&pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
		if err != nil {
			return err
		}
	}

	// Novo canal para este cliente
	clientCh := make(chan *pb.Category, 10)
	c.addSubscriber(clientCh)
	defer c.removeSubscriber(clientCh)

	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return nil
		case newCategory := <-clientCh:
			if err := stream.Send(newCategory); err != nil {
				return err
			}
		}
	}
}

func (c *CategoryService) addSubscriber(ch chan *pb.Category) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subscribers[ch] = struct{}{}
}

func (c *CategoryService) removeSubscriber(ch chan *pb.Category) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.subscribers, ch)
	close(ch)
}

func (c *CategoryService) broadcast(category *pb.Category) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for ch := range c.subscribers {
		select {
		case ch <- category:
		default:
			// Canal cheio ou desconectado, remover
			delete(c.subscribers, ch)
			close(ch)
		}
	}
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryId) (*pb.Category, error) {
	category, err := c.CategoryDB.Find(in.Id)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		category, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(categories)
		}

		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})

		if err != nil {
			return err
		}
	}
}
