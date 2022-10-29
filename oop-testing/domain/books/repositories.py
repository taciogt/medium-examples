from typing import List, Dict
from uuid import UUID, uuid4
from abc import ABC, abstractmethod

from domain.books.entities import Book


class WriteBooksRepository(ABC):
    @abstractmethod
    def save_book(self, book: Book) -> None:
        ...


class ReadBooksRepository(ABC):
    @abstractmethod
    def get_book(self, _id: UUID) -> Book:
        ...


class BooksRepository(WriteBooksRepository, ReadBooksRepository, ABC):
    ...


class InMemoryRepository(BooksRepository):
    books: Dict[UUID, Book]

    def __init__(self) -> None:
        self.books = dict()

    def save_book(self, book: Book) -> None:
        book.id = uuid4()
        self.books[book.id] = book

    def get_book(self, _id: UUID) -> Book:
        if _id in self.books:
            return self.books[_id]
        raise Exception("book not found")
