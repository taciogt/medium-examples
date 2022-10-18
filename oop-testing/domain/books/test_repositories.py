from datetime import datetime
from unittest import TestCase

from domain.books.entities import Book
from domain.books.repositories import InMemoryRepository, BooksRepository


class TestInMemoryRepository(TestCase):
    repo: BooksRepository = InMemoryRepository()

    def test_save(self):
        book_to_save = Book(title="Lord of the Rings", author="J.R.R. Tolkien",
                            publishing_date=datetime(year=1954, month=7, day=29))

        returned_book = self.repo.save_book(book=book_to_save)

        saved_book = self.repo.get_book(_id=returned_book.id)
        self.assertEqual(returned_book, saved_book)
