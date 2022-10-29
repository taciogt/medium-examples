from abc import ABC, abstractmethod
from datetime import datetime
from unittest import TestCase
from unittest.case import skip
from uuid import UUID, uuid4

from domain.books.entities import Book
from domain.books.repositories import InMemoryRepository, BooksRepository


class TestBooksRepository(ABC, TestCase):
    repo: BooksRepository

    def setUp(self) -> None:
        self.repo = self.get_repository()

    @abstractmethod
    def get_repository(self) -> BooksRepository:
        ...

    @abstractmethod
    def get_book(self, id_: UUID) -> Book:
        ...

    @abstractmethod
    def save_book(self, book: Book) -> None:
        ...


class TestInMemoryRepository(TestBooksRepository):
    _repo: InMemoryRepository

    def get_repository(self) -> BooksRepository:
        self._repo = InMemoryRepository()
        return self._repo

    def get_book(self, id_: UUID) -> Book:
        return self._repo.books[id_]

    def save_book(self, book: Book) -> None:
        if book.id is None:
            return
        self._repo.books[book.id] = book


class TestRepository(TestInMemoryRepository):
    def test_get(self) -> None:
        book_to_get = Book(title='Lord of the Rings', author='J.R.R. Tolkien',
                           publishing_date=datetime(year=1954, month=7, day=29))
        book_to_get.id = uuid4()

        self.save_book(book=book_to_get)

        got_book = self.repo.get_book(_id=book_to_get.id)
        self.assertEqual(book_to_get, got_book)

    def test_get_non_existing_book(self) -> None:
        book_to_get = Book(title='Lord of the Rings', author='J.R.R. Tolkien',
                           publishing_date=datetime(year=1954, month=7, day=29))
        book_to_get.id = uuid4()

        self.save_book(book=book_to_get)

        self.assertRaisesRegex(Exception, 'book not found', self.repo.get_book, _id=uuid4())

    def test_save(self) -> None:
        book_to_save = Book(title='Lord of the Rings', author='J.R.R. Tolkien',
                            publishing_date=datetime(year=1954, month=7, day=29))

        self.repo.save_book(book=book_to_save)

        self.assertIsInstance(book_to_save.id, UUID)
        if book_to_save.id is None:
            return
        saved_book = self.get_book(id_=book_to_save.id)
        self.assertEqual(book_to_save, saved_book)
