package service

import (
	"context"
	"strconv"

	pb "book/api/book"
	"book/internal/biz"
)

type BookService struct {
	pb.UnimplementedBookServer

	uc *biz.BookUsecase
}

func NewBookService(uc *biz.BookUsecase) *BookService {
	return &BookService{uc: uc}
}

func (s *BookService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookReply, error) {
	g, err := s.uc.CreateBook(ctx, &biz.Book{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb.CreateBookReply{
		Name: "Book Create Success,Name:" + g.Name + ",Count:" + strconv.Itoa(g.Count),
	}, nil
}

func (s *BookService) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookReply, error) {
	return &pb.UpdateBookReply{}, nil
}

func (s *BookService) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookReply, error) {
	return &pb.DeleteBookReply{}, nil
}

func (s *BookService) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookReply, error) {
	return &pb.GetBookReply{}, nil
}

func (s *BookService) ListBook(ctx context.Context, req *pb.ListBookRequest) (*pb.ListBookReply, error) {
	return &pb.ListBookReply{}, nil
}
