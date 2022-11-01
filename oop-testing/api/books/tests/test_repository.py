from uuid import UUID

from books.models import BookModel
from books.repositories import DjangoBooksRepository
from domain.books.entities import Book
from domain.books.repositories import BooksRepository
from domain.books.test_repositories import TestBooksRepository, BaseBooksRepositoryTest


class TestDjangoRepository(BaseBooksRepositoryTest):

    # _repo: DjangoBooksRepository

    def get_repository(self) -> BooksRepository:
        return DjangoBooksRepository()
        # return self._repo

    def get_book(self, id_: UUID) -> Book:
        # return self._repo.books[id_]
        raise NotImplemented

    def save_book(self, book: Book) -> None:
        BookModel(id=book.id, author=book.author, title=book.title, publishing_date=book.publishing_date).save()


class TestRepository(TestDjangoRepository, TestBooksRepository):
    ...
