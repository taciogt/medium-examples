from uuid import UUID

from django.test import TestCase

from books.models import BookModel
from books.repositories import DjangoBooksRepository
from domain.books.entities import Book
from domain.books.repositories import BooksRepository
from domain.books.test_repositories import TestBooksRepository, BaseBooksRepositoryTest


class TestDjangoRepository(BaseBooksRepositoryTest, TestCase):

    def get_repository(self) -> BooksRepository:
        return DjangoBooksRepository()

    def get_book(self, id_: UUID) -> Book:
        return BookModel.objects.get(id=id_).to_entity()

    def save_book(self, book: Book) -> None:
        BookModel(id=book.id, author=book.author, title=book.title, publishing_date=book.publishing_date).save()


class TestRepository(TestDjangoRepository, TestBooksRepository):
    ...
